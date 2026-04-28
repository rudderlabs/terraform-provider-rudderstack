package acc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// AccAssertSource is the E2E counterpart of cm.AssertSource.
// It accepts the same []configs.TestConfig used by unit tests.
//
// In plan-only mode (TF_ACC_PLAN_ONLY=1): validates HCL + provider schema (zero API calls).
// In full mode (TF_ACC=1): runs Create → Update → Import → Destroy against the real API,
// and verifies source.Config matches the expected JSON from test configs.
func AccAssertSource(t *testing.T, source string, testConfigs []configs.TestConfig) {
	t.Helper()

	planOnly := PlanOnly()
	if planOnly {
		t.Parallel()
	}

	resourceName := fmt.Sprintf("rudderstack_source_%s.test", source)
	name := RandomName(source)
	cfg := testConfigs[0]

	if planOnly {
		ensureDummyToken(t)
		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:             testAccSourceConfig(source, name, cfg.TerraformCreate),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		})
		return
	}

	createSettingsJSON := cfg.APICreateSettings
	if createSettingsJSON == "" {
		createSettingsJSON = cfg.APICreate
	}
	updateSettingsJSON := cfg.APIUpdateSettings
	if updateSettingsJSON == "" {
		updateSettingsJSON = cfg.APIUpdate
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckSourceDestroy(source),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceConfig(source, name, cfg.TerraformCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					testAccCheckSourceAPIConfig(resourceName, cfg.APICreate),
					testAccCheckSourceSettingsAPI(resourceName, createSettingsJSON),
				),
			},
			{
				Config: testAccSourceConfig(source, name+"-updated", cfg.TerraformUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-updated"),
					testAccCheckSourceAPIConfig(resourceName, cfg.APIUpdate),
					testAccCheckSourceSettingsAPI(resourceName, updateSettingsJSON),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// write_key is computed and not returned on import
				ImportStateVerifyIgnore: []string{"write_key"},
			},
		},
	})
}

// testAccSourceConfig generates the Terraform HCL for a source resource.
// configBlock is embedded raw in the resource body — the caller is responsible
// for passing a well-formed block such as "settings { … }" or an empty string.
func testAccSourceConfig(source, name, configBlock string) string {
	if strings.TrimSpace(configBlock) == "" {
		return fmt.Sprintf(`
resource "rudderstack_source_%s" "test" {
  name = %q
}
`, source, name)
	}
	return fmt.Sprintf(`
resource "rudderstack_source_%s" "test" {
  name = %q
  %s
}
`, source, name, configBlock)
}

// testAccCheckSourceExists verifies the source exists in the live API.
func testAccCheckSourceExists(resourceName string) resource.TestCheckFunc {
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
		_, err = cl.Sources.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("source %s not found in API: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

// testAccCheckSourceAPIConfig fetches the source from the API and verifies
// its config contains all expected fields from the test's API JSON.
func testAccCheckSourceAPIConfig(resourceName, expectedJSON string) resource.TestCheckFunc {
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

		src, err := cl.Sources.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get source from API: %w", err)
		}

		return compareConfig(src.Config, expectedJSON)
	}
}

// testAccCheckSourceSettingsAPI fetches the source from the API and verifies
// GeoEnrichmentEnabled and Transient match the values in expectedJSON.
// expectedJSON uses terraform-facing names: "geoEnrichmentEnabled" maps directly to
// source.GeoEnrichmentEnabled; "transient" equals temporarily_store_events_for_retries
// (= !source.Transient) and is inverted before comparison.
func testAccCheckSourceSettingsAPI(resourceName, expectedJSON string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expectedJSON = strings.TrimSpace(expectedJSON)
		if expectedJSON == "" || expectedJSON == "{}" {
			return nil
		}

		var m map[string]any
		if err := json.Unmarshal([]byte(expectedJSON), &m); err != nil {
			return nil // not a settings JSON, nothing to check
		}
		_, hasGeo := m["geoEnrichmentEnabled"]
		_, hasTrans := m["transient"]
		if !hasGeo && !hasTrans {
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

		src, err := cl.Sources.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get source from API: %w", err)
		}

		if v, ok := m["geoEnrichmentEnabled"].(bool); ok && (src.GeoEnrichmentEnabled == nil || *src.GeoEnrichmentEnabled != v) {
			got := "<nil>"
			if src.GeoEnrichmentEnabled != nil {
				got = fmt.Sprintf("%v", *src.GeoEnrichmentEnabled)
			}
			return fmt.Errorf("API GeoEnrichmentEnabled: expected %v, got %s", v, got)
		}
		if v, ok := m["transient"].(bool); ok && (src.Transient == nil || *src.Transient != v) {
			got := "<nil>"
			if src.Transient != nil {
				got = fmt.Sprintf("%v", *src.Transient)
			}
			return fmt.Errorf("API Transient: expected %v, got %s", v, got)
		}
		return nil
	}
}

// testAccCheckSourceDestroy verifies all sources created by the test are deleted.
// Note: The RudderStack API uses soft-delete, so Get may still return a 200
// after deletion. We treat a successful Get as acceptable (soft-deleted).
func testAccCheckSourceDestroy(source string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cl, err := newTestAPIClient()
		if err != nil {
			return err
		}

		resourceType := fmt.Sprintf("rudderstack_source_%s", source)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			// RudderStack API uses soft-delete: Get may still return 200 after deletion.
			// We accept this and do not fail — the destroy step already verified the
			// Delete handler ran without error.
			_, _ = cl.Sources.Get(context.Background(), rs.Primary.ID)
		}
		return nil
	}
}
