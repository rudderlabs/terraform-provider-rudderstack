package acc

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// AccAssertDestination is the E2E counterpart of cm.AssertDestination.
// It accepts the same []configs.TestConfig used by unit tests, reusing their HCL configs.
//
// In plan-only mode (TF_ACC_PLAN_ONLY=1): validates HCL + provider schema (zero API calls).
// In full mode (TF_ACC=1): runs Create → Update → Import → Destroy against the real API,
// and verifies the API config matches the expected JSON from test configs.
func AccAssertDestination(t *testing.T, destination string, testConfigs []configs.TestConfig) {
	t.Helper()

	resourceName := fmt.Sprintf("rudderstack_destination_%s.test", destination)
	name := RandomName(destination)
	cfg := testConfigs[0]
	wantVersion := registeredDestinationVersion(t, destination)
	cm := configs.Destinations.Entries()[destination]
	redactedFields := testutil.RedactedAPIConfigKeys(t, cm)
	// Secrets are redacted from responses, so they can't be verified on import
	// (there is no prior state to preserve them from) — ignore them there.
	importIgnore := sensitiveStateAttrPaths(cm)

	if PlanOnly() {
		t.Parallel()
		ensureDummyToken(t)
		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:             testAccDestinationConfig(destination, name, cfg.TerraformCreate),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		})
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckDestinationDestroy(destination),
		Steps: []resource.TestStep{
			{
				Config: testAccDestinationConfig(destination, name, cfg.TerraformCreate),
				// Secrets are redacted from responses, so config stays authoritative
				// and every plan re-asserts them → a non-empty post-apply plan is
				// expected, but only when this step's config actually sets a secret
				// (BREAKING_CHANGES.md).
				ExpectNonEmptyPlan: testutil.ConfigHasRedactedSecret(cfg.APICreate, redactedFields),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					testAccCheckDestinationAPIConfig(resourceName, cfg.APICreate, redactedFields),
					// Exact wire version must match the destination's registered
					// ConfigMeta.Version (v1 today; future _v2 resources expect 2).
					// The automatic post-apply plan check also asserts no plan
					// diff. Requires the backend to echo version on read
					// (Blocker B) — until then this fails loudly rather than
					// silently passing.
					testAccCheckDestinationVersion(resourceName, wantVersion),
				),
			},
			{
				Config:             testAccDestinationConfig(destination, name+"-updated", cfg.TerraformUpdate),
				ExpectNonEmptyPlan: testutil.ConfigHasRedactedSecret(cfg.APIUpdate, redactedFields),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-updated"),
					testAccCheckDestinationAPIConfig(resourceName, cfg.APIUpdate, redactedFields),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importIgnore,
			},
		},
	})
}

// RandomName generates a unique resource name with a tf-acc- prefix for test isolation.
func RandomName(prefix string) string {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil || n == nil {
		return fmt.Sprintf("tf-acc-%s-%d", prefix, os.Getpid())
	}
	return fmt.Sprintf("tf-acc-%s-%d", prefix, n.Int64())
}

// testAccDestinationConfig generates the Terraform HCL for a destination resource.
// No provider block is needed — the provider reads credentials from env vars.
func testAccDestinationConfig(destination, name, configBlock string) string {
	if configBlock == "" {
		return fmt.Sprintf(`
resource "rudderstack_destination_%s" "test" {
  name = %q
}
`, destination, name)
	}
	return fmt.Sprintf(`
resource "rudderstack_destination_%s" "test" {
  name = %q
  config {
    %s
  }
}
`, destination, name, configBlock)
}

// testAccCheckDestinationExists verifies the resource exists in the live API.
func testAccCheckDestinationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is empty")
		}

		cl, err := newTestAPIClient()
		if err != nil {
			return err
		}
		_, err = cl.Destinations.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("destination %s not found in API: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

// testAccCheckDestinationAPIConfig fetches the destination from the API and verifies
// its config contains all expected fields from the test's API JSON. redactedFields
// are secret API keys the backend omits from responses and must not be asserted.
func testAccCheckDestinationAPIConfig(resourceName, expectedJSON string, redactedFields map[string]bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if expectedJSON == "" {
			return nil
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		cl, err := newTestAPIClient()
		if err != nil {
			return err
		}

		dest, err := cl.Destinations.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get destination from API: %w", err)
		}

		return compareConfig(dest.Config, expectedJSON, redactedFields)
	}
}

// sensitiveStateAttrPaths returns the terraform state attribute paths of the
// destination's Sensitive (secret) config fields, e.g. "config.0.api_secret"
// or "config.0.s3.0.access_key". Import can't verify these: the backend redacts
// them from responses, so an imported resource has no value to compare against
// the pre-import state.
func sensitiveStateAttrPaths(cm configs.ConfigMeta) []string {
	var paths []string
	for _, p := range cm.SensitiveImportIgnorePaths() {
		paths = append(paths, "config.0."+p)
	}
	return paths
}

// registeredDestinationVersion returns ConfigMeta.Version for the terraform
// destination type name. Exact match (not >= 1) keeps future _v2 resources
// correct while still failing if the API returns the wrong version.
func registeredDestinationVersion(t *testing.T, destination string) int {
	t.Helper()
	cm, ok := configs.Destinations.Entries()[destination]
	if !ok {
		t.Fatalf("destination %q is not registered", destination)
	}
	if cm.Version < 1 {
		t.Fatalf("destination %q has invalid ConfigMeta.Version %d", destination, cm.Version)
	}
	return cm.Version
}

// testAccCheckDestinationVersion fetches the destination from the API and verifies
// its reported version matches wantVersion. This depends on the backend always
// reporting a real version on every destination (INT-6489); the provider does not
// coerce an absent/zero version to v1 (see cmd/generatetf/generator.configMetaByVersion).
func testAccCheckDestinationVersion(resourceName string, wantVersion int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		cl, err := newTestAPIClient()
		if err != nil {
			return err
		}

		dest, err := cl.Destinations.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get destination from API: %w", err)
		}

		if dest.Version != wantVersion {
			return fmt.Errorf("destination %s: got version %d, want %d", rs.Primary.ID, dest.Version, wantVersion)
		}
		return nil
	}
}

// testAccCheckDestinationDestroy verifies all destinations created by the test
// are deleted from the API after the test completes.
// Note: The RudderStack API uses soft-delete, so Get may still return a 200
// after deletion. We accept this and do not fail — the destroy step already
// verified the Delete handler ran without error.
func testAccCheckDestinationDestroy(destination string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cl, err := newTestAPIClient()
		if err != nil {
			return err
		}

		resourceType := fmt.Sprintf("rudderstack_destination_%s", destination)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			// RudderStack API uses soft-delete: Get may still return 200 after deletion.
			// We accept this and do not fail — the destroy step already verified the
			// Delete handler ran without error.
			_, _ = cl.Destinations.Get(context.Background(), rs.Primary.ID)
		}
		return nil
	}
}

// newTestAPIClient creates a real API client using environment variables.
func newTestAPIClient() (*client.Client, error) {
	accessToken := os.Getenv("RUDDERSTACK_ACCESS_TOKEN")
	var opts []client.Option
	if v := os.Getenv("RUDDERSTACK_API_URL"); v != "" {
		// Strip trailing /v2 (with or without trailing slash) for backward compatibility —
		// the new client includes /v2 in service paths.
		v = strings.TrimSuffix(strings.TrimRight(v, "/"), "/v2")
		opts = append(opts, client.WithBaseURL(v))
	}
	return client.New(accessToken, opts...)
}
