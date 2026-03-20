// integration-verify is a CLI tool for validating integrations onboarded via
// the /onboard-integration Claude Code skill.
//
// After the skill generates the .go, test, example, and docs files for a new
// integration, the optional E2E step deploys the resource to a real RudderStack
// instance. This tool then compares the Terraform .tf config against the live
// API response to confirm the onboarded integration works end-to-end.
//
// Usage:
//
//	go run ./cmd/integration-verify/ -file <path.tf> -id <resource-id> [-resource <name>]
//
// Environment variables:
//
//	RUDDERSTACK_ACCESS_TOKEN  (required) API token for the RudderStack control plane.
//	RUDDERSTACK_API_URL       (optional) Override the default API base URL.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rudderlabs/rudder-api-go/client"

	// Blank-import integrations so every registered source/destination is
	// available for schema lookup when parsing .tf files.
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	if code != 0 {
		os.Exit(code)
	}
}

// run executes the integration-verify CLI logic and returns an exit code.
// It writes output to stdout and errors to stderr.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("integration-verify", flag.ContinueOnError)
	fs.SetOutput(stderr)

	filePath := fs.String("file", "", "path to the .tf file containing the onboarded integration resource (required)")
	resourceID := fs.String("id", "", "resource ID from terraform apply (required)")
	targetResource := fs.String("resource", "", "resource name to verify (optional, defaults to first rudderstack resource)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *filePath == "" || *resourceID == "" {
		fmt.Fprintf(stderr, "Usage: integration-verify -file <path.tf> -id <resource-id> [-resource <name>]\n\n")
		fmt.Fprintf(stderr, "Validates that an integration onboarded via /onboard-integration matches the live API.\n\n")
		fs.PrintDefaults()
		return 1
	}

	// Parse the .tf file to extract the onboarded integration's resource info.
	info, err := ParseTFFile(*filePath, *targetResource)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err.Error())
		return 1
	}

	cl, err := setupClient()
	if err != nil {
		fmt.Fprintf(stderr, "Error: could not create API client: %s\n", err.Error())
		return 1
	}

	// Verify the onboarded integration: convert .tf state → expected API JSON,
	// fetch the actual config from the API, and compare.
	result, err := Verify(context.Background(), cl, info, *resourceID)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err.Error())
		return 1
	}

	fmt.Fprint(stdout, FormatResult(info, *resourceID, result))

	if !result.Match {
		return 1
	}
	return 0
}

func setupClient() (*client.Client, error) {
	accessToken := os.Getenv("RUDDERSTACK_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("RUDDERSTACK_ACCESS_TOKEN environment variable is required")
	}

	baseURL := os.Getenv("RUDDERSTACK_API_URL")
	if baseURL == "" {
		baseURL = client.BASE_URL_V2
	}

	return client.New(accessToken, client.WithBaseURL(baseURL))
}
