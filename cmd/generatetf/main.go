package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/cmd/generatetf/generator"
)

var (
	cl      *client.Client
	retlSvc retl.RETLStore

	// retlSourceIDsSeenViaLegacyEndpoint collects IDs of warehouse-backed
	// RETL sources that the legacy /v2/sources endpoint also returns. We
	// drop them from the event-streaming source list (no rudderstack_source_*
	// resource exists for them — they're emitted via the RETL path instead)
	// and use the set to silence connection warnings for connections that
	// reference them in /v2/connections.
	retlSourceIDsSeenViaLegacyEndpoint = map[string]bool{}
)

// retlWarehouseSourceTypes is the set of `type` values returned by the legacy
// /v2/sources endpoint for sources that are actually RETL warehouse sources.
// The bootstrap tool emits them via the typed /v2/retl-sources path instead
// — silencing them here avoids "type X not supported" noise for every RETL
// source in the workspace.
var retlWarehouseSourceTypes = map[string]bool{
	"redshift":   true,
	"snowflake":  true,
	"bigquery":   true,
	"postgres":   true,
	"mysql":      true,
	"databricks": true,
	"trino":      true,
	"s3":         true,
}

// retlListPageSize is the page size used when paginating RETL connections.
const retlListPageSize = 50

func main() {
	importFlag := flag.Bool("import", false, "generate terraform import commands")
	flag.Parse()

	var err error
	cl, retlSvc, err = setupClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create API client: %s\n", err.Error())
		os.Exit(1)
	}

	sources, err := getAPISources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get list of sources: %s\n", err.Error())
		os.Exit(1)
	}

	destinations, err := getAPIDestinations()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get list of destinations: %s\n", err.Error())
		os.Exit(1)
	}

	connections, err := getAPIConnections()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get list of connections: %s\n", err.Error())
		os.Exit(1)
	}

	retlSources, err := getAPIRetlSources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get list of RETL sources: %s\n", err.Error())
		os.Exit(1)
	}

	retlConnections, err := getAPIRetlConnections()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get list of RETL connections: %s\n", err.Error())
		os.Exit(1)
	}

	if *importFlag {
		bytes, err := generator.GenerateImportScript(sources, destinations, connections, retlSources, retlConnections)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not generate terraform import script: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else {
		bytes, err := generator.GenerateTerraform(sources, destinations, connections, retlSources, retlConnections)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not generate terraform HCL: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	}
}

func setupClient() (*client.Client, retl.RETLStore, error) {
	accessToken := os.Getenv("RUDDERSTACK_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, nil, fmt.Errorf("no access token is specified. Please provide one through the RUDDERSTACK_ACCESS_TOKEN environmental variable")
	}

	baseURL := os.Getenv("RUDDERSTACK_API_URL")
	if baseURL == "" {
		baseURL = client.BASE_URL
	}
	// Strip trailing /v2 (with or without trailing slash) for backward compatibility —
	// the new client includes /v2 in service paths.
	baseURL = strings.TrimSuffix(strings.TrimRight(baseURL, "/"), "/v2")

	c, err := client.New(accessToken, client.WithBaseURL(baseURL))
	if err != nil {
		return nil, nil, err
	}
	return c, retl.NewRudderRETLStore(c), nil
}

func getAPISources() ([]client.Source, error) {
	var sources []client.Source
	sourcesPage, err := cl.Sources.List(context.Background())
	if err != nil {
		return nil, err
	}

	for sourcesPage != nil {
		for _, src := range sourcesPage.Sources {
			if retlWarehouseSourceTypes[strings.ToLower(src.Type)] {
				retlSourceIDsSeenViaLegacyEndpoint[src.ID] = true
				continue
			}
			sources = append(sources, src)
		}
		sourcesPage, err = cl.Sources.Next(context.Background(), sourcesPage.Paging)
		if err != nil {
			return nil, err
		}
	}

	return sources, nil
}

func getAPIDestinations() ([]client.Destination, error) {
	var destinations []client.Destination
	destinationsPage, err := cl.Destinations.List(context.Background())
	if err != nil {
		return nil, err
	}

	for destinationsPage != nil {
		destinations = append(destinations, destinationsPage.Destinations...)
		destinationsPage, err = cl.Destinations.Next(context.Background(), destinationsPage.Paging)
		if err != nil {
			return nil, err
		}
	}

	return destinations, nil
}

func getAPIConnections() ([]client.Connection, error) {
	var connections []client.Connection
	connectionsPage, err := cl.Connections.List(context.Background())
	if err != nil {
		return nil, err
	}

	for connectionsPage != nil {
		for _, cnxn := range connectionsPage.Connections {
			// Drop connections whose source was filtered out as a RETL
			// warehouse source — those are emitted via rudderstack_retl_connection
			// (or their RETL connection equivalent) on the RETL path.
			if retlSourceIDsSeenViaLegacyEndpoint[cnxn.SourceID] {
				continue
			}
			connections = append(connections, cnxn)
		}
		connectionsPage, err = cl.Connections.Next(context.Background(), connectionsPage.Paging)
		if err != nil {
			return nil, err
		}
	}

	return connections, nil
}

func getAPIRetlSources() ([]retl.RETLSource, error) {
	// Two filtered calls — one per supported sourceType — even though the SDK
	// doc claims "Pass an empty sourceType to return sources of all types".
	// In practice the server treats sourceType as required: omitting it
	// silently returns only `model` sources and drops tables entirely. Until
	// the API contract is fixed, we must paginate per type and concatenate.
	//
	// TODO: ListRetlSources is not paginated in the SDK today (RETLSources
	// has no Paging field). Once the upstream SDK adds pagination, switch
	// each branch to a paged loop — otherwise large workspaces will silently
	// truncate at whatever limit the server applies.
	var sources []retl.RETLSource
	for _, st := range []retl.SourceType{retl.ModelSourceType, retl.TableSourceType} {
		resp, err := retlSvc.ListRetlSources(context.Background(), retl.WithSourceType(string(st)))
		if err != nil {
			return nil, fmt.Errorf("listing RETL sources of type %q: %w", st, err)
		}
		sources = append(sources, resp.Data...)
	}
	return sources, nil
}

func getAPIRetlConnections() ([]retl.RETLConnection, error) {
	var out []retl.RETLConnection
	for page := 1; ; page++ {
		resp, err := retlSvc.ListConnections(context.Background(), &retl.ListRETLConnectionsRequest{
			Page:     page,
			PageSize: retlListPageSize,
		})
		if err != nil {
			return nil, err
		}
		out = append(out, resp.Data...)
		// Stop when an empty page comes back, or when the server-reported Total
		// is populated and we've fetched it all. We only consult Total when
		// it's > 0 — a zero Total can mean either "really empty" (handled by
		// the empty-data check) or "server didn't populate it", and in the
		// latter case treating 0 as the stop condition would prematurely break.
		if len(resp.Data) == 0 || (resp.Paging.Total > 0 && len(out) >= resp.Paging.Total) {
			break
		}
	}
	return out, nil
}
