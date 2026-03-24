// integration-verify is a CLI tool for validating integrations onboarded via
// the /onboard-integration Claude Code skill.
//
// After the skill generates the .go, test, example, and docs files for a new
// integration, the optional E2E step deploys the resource to a real RudderStack
// instance. This tool then reads the terraform state and compares it against the
// live API response to confirm the onboarded integration works end-to-end.
//
// Usage:
//
//	go run ./cmd/integration-verify/ -dir <terraform-dir> [-resource <name>]
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
	"log"
	"os"
	"os/exec"

	"github.com/rudderlabs/rudder-api-go/client"

	// Blank-import integrations so every registered source/destination is
	// available for schema lookup when parsing terraform state.
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: integration-verify -dir <terraform-dir> [-resource <name>]\n\n")
		fmt.Fprintf(os.Stderr, "Validates that an integration onboarded via /onboard-integration matches the live API.\n\n")
		flag.PrintDefaults()
	}

	dir := flag.String("dir", "", "terraform working directory containing applied state (required)")
	targetResource := flag.String("resource", "", "resource name to verify (optional, defaults to all rudderstack resources)")
	flag.Parse()

	if *dir == "" {
		flag.Usage()
		os.Exit(1)
	}

	cmd := exec.Command("terraform", "show", "-json")
	cmd.Dir = *dir
	stateJSON, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error: running terraform show: %s", err)
	}

	if err := verifyFromState(stateJSON, *targetResource); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

// verifyFromState performs verification using pre-fetched terraform state JSON.
// This is separated from main() to enable testing without requiring the terraform CLI.
func verifyFromState(stateJSON []byte, targetResource string) error {
	resources, err := ParseTerraformState(stateJSON, targetResource)
	if err != nil {
		return err
	}

	cl, err := setupClient()
	if err != nil {
		return fmt.Errorf("could not create API client: %w", err)
	}

	for _, info := range resources {
		result, err := Verify(context.Background(), cl, info)
		if err != nil {
			return err
		}

		if !result.Match {
			return fmt.Errorf("verification failed for %s (ID: %s)\n\nDiff (-expected +actual):\n%s", info.ResourceType, info.ResourceID, result.Diff)
		}
	}

	return nil
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
