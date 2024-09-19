package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/cmd/generatetf/generator"
)

var cl *client.Client

func main() {
	importFlag := flag.Bool("import", false, "generate terraform import commands")
	flag.Parse()

	var err error
	cl, err = setupClient()
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

	if *importFlag {
		bytes, err := generator.GenerateImportScript(sources, destinations, connections)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not generate terraform import script: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else {
		bytes, err := generator.GenerateTerraform(sources, destinations, connections)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not generate terraform HCL: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	}
}

func setupClient() (*client.Client, error) {
	accessToken := os.Getenv("RUDDERSTACK_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("no access token in specified. Please provide one through the RUDDERSTACK_ACCESS_TOKEN environmental variable")
	}

	baseURL := os.Getenv("RUDDERSTACK_API_URL")
	if baseURL == "" {
		baseURL = client.BASE_URL_V2
	}

	return client.New(accessToken, client.WithBaseURL(baseURL))
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
