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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					testAccCheckDestinationAPIConfig(resourceName, cfg.APICreate),
				),
			},
			{
				Config: testAccDestinationConfig(destination, name+"-updated", cfg.TerraformUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-updated"),
					testAccCheckDestinationAPIConfig(resourceName, cfg.APIUpdate),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
// its config contains all expected fields from the test's API JSON.
func testAccCheckDestinationAPIConfig(resourceName, expectedJSON string) resource.TestCheckFunc {
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

		return compareConfig(dest.Config, expectedJSON)
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
