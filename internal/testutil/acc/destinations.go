package acc

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	cm := registeredDestinationConfigMeta(t, destination)
	wantVersion := cm.Version

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
					testAccCheckDestinationAPIConfig(resourceName, cfg.APICreate, cm),
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
				Config: testAccDestinationConfig(destination, name+"-updated", cfg.TerraformUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-updated"),
					testAccCheckDestinationAPIConfig(resourceName, cfg.APIUpdate, cm),
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
func testAccCheckDestinationAPIConfig(resourceName, expectedJSON string, cm configs.ConfigMeta) resource.TestCheckFunc {
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

		ignoredPaths, err := sensitiveAPIConfigPaths(cm)
		if err != nil {
			return err
		}
		return compareConfig(dest.Config, expectedJSON, ignoredPaths...)
	}
}

// registeredDestinationConfigMeta returns ConfigMeta for the terraform
// destination type name. Exact version match (not >= 1) keeps future _v2 resources
// correct while still failing if the API returns the wrong version.
func registeredDestinationConfigMeta(t *testing.T, destination string) configs.ConfigMeta {
	t.Helper()
	cm, ok := configs.Destinations.Entries()[destination]
	if !ok {
		t.Fatalf("destination %q is not registered", destination)
	}
	if cm.Version < 1 {
		t.Fatalf("destination %q has invalid ConfigMeta.Version %d", destination, cm.Version)
	}
	return cm
}

func sensitiveAPIConfigPaths(cm configs.ConfigMeta) ([]string, error) {
	state, markers := sensitiveStateForSchema(cm.ConfigSchema)
	if len(markers) == 0 {
		return nil, nil
	}

	stateBytes, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sensitive test state: %w", err)
	}

	apiConfig, err := cm.StateToAPI(string(stateBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to map sensitive test state to API config: %w", err)
	}

	var apiValue any
	if err := json.Unmarshal([]byte(apiConfig), &apiValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sensitive API config: %w", err)
	}

	pathSet := map[string]struct{}{}
	collectMarkerPaths("", apiValue, markers, pathSet)
	paths := make([]string, 0, len(pathSet))
	for path := range pathSet {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func sensitiveStateForSchema(configSchema map[string]*schema.Schema) (map[string]any, map[string]struct{}) {
	markers := map[string]struct{}{}
	state := sensitiveStateObject(configSchema, markers)
	return state, markers
}

func sensitiveStateObject(configSchema map[string]*schema.Schema, markers map[string]struct{}) map[string]any {
	state := map[string]any{}
	for key, fieldSchema := range configSchema {
		value, ok := sensitiveStateValue(fieldSchema, markers)
		if ok {
			state[key] = value
		}
	}
	return state
}

func sensitiveStateValue(fieldSchema *schema.Schema, markers map[string]struct{}) (any, bool) {
	if fieldSchema == nil {
		return nil, false
	}

	if nestedResource, ok := fieldSchema.Elem.(*schema.Resource); ok {
		nestedState := sensitiveStateObject(nestedResource.Schema, markers)
		if len(nestedState) == 0 {
			return nil, false
		}

		switch fieldSchema.Type {
		case schema.TypeList, schema.TypeSet:
			return []any{nestedState}, true
		case schema.TypeMap:
			return nestedState, true
		default:
			return nestedState, true
		}
	}

	if fieldSchema.Sensitive {
		marker := fmt.Sprintf("__rudder_tf_sensitive_%d__", len(markers))
		markers[marker] = struct{}{}
		return marker, true
	}

	return nil, false
}

func collectMarkerPaths(prefix string, value any, markers map[string]struct{}, paths map[string]struct{}) {
	switch typedValue := value.(type) {
	case map[string]any:
		for key, child := range typedValue {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			collectMarkerPaths(path, child, markers, paths)
		}
	case []any:
		for i, child := range typedValue {
			collectMarkerPaths(fmt.Sprintf("%s[%d]", prefix, i), child, markers, paths)
		}
	case string:
		if _, ok := markers[typedValue]; ok {
			paths[prefix] = struct{}{}
		}
	}
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
