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
)

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
		return nil, nil, fmt.Errorf("no access token in specified. Please provide one through the RUDDERSTACK_ACCESS_TOKEN environmental variable")
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
		sources = append(sources, sourcesPage.Sources...)
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
		connections = append(connections, connectionsPage.Connections...)
		connectionsPage, err = cl.Connections.Next(context.Background(), connectionsPage.Paging)
		if err != nil {
			return nil, err
		}
	}

	return connections, nil
}

func getAPIRetlSources() ([]retl.RETLSource, error) {
	retlSources := []retl.RETLSource{}
	resp, err := retlSvc.ListRetlSources(context.Background(), retl.WithSourceType(string(retl.ModelSourceType)))
	if err != nil {
		return nil, err
	}
	retlSources = append(retlSources, resp.Data...)
	resp, err = retlSvc.ListRetlSources(context.Background(), retl.WithSourceType(string(retl.TableSourceType)))
	if err != nil {
		return nil, err
	}
	retlSources = append(retlSources, resp.Data...)
	return retlSources, nil
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
		// Stop when the server says we have everything, or as a defensive fallback
		// when an empty page comes back (in case Total is missing/zero).
		if len(resp.Data) == 0 || len(out) >= resp.Paging.Total {
			break
		}
	}
	return out, nil
}
