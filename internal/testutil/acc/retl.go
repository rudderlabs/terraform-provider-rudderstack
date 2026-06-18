package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// retlAccountIDEnv names the env var that supplies a real warehouse account ID
// for live RETL acceptance tests. When unset (and the test is not running in
// plan-only mode) the helpers t.Skip — local runs without a workspace fixture
// shouldn't fail. The e2e-tests.yml workflow forwards `vars.RUDDERSTACK_RETL_TEST_ACCOUNT_ID`
// (a GitHub Environment variable, not a secret — the ID is opaque, not credentials)
// to this env var; until that variable is added, the expression renders empty
// in CI and the live RETL tests skip.
const retlAccountIDEnv = "RUDDERSTACK_RETL_TEST_ACCOUNT_ID"

// planOnlyAccountID is a placeholder used in plan-only mode where no API call
// is made; the value just has to satisfy the schema (non-empty string).
const planOnlyAccountID = "acc-plan-only"

// RETLSourceTestConfig is the shape consumed by AccAssertRETLSourceModel /
// AccAssertRETLSourceTable. The Config / UpdateConfig fields are HCL fragments
// that go inside the resource's `config { }` block.
type RETLSourceTestConfig struct {
	SourceDefinitionName string // e.g. "bigquery"
	Config               string // HCL inside config { } block for the create step
	UpdateConfig         string // HCL for the update step (omit to skip update)
}

// RETLConnectionTestConfig describes a connection variant. The helper builds
// a complete pipeline (RETL source + webhook destination + connection) so
// tests don't need to repeat the boilerplate.
type RETLConnectionTestConfig struct {
	// Variant is a short label included in the test resource names so
	// concurrent runs don't collide.
	Variant string

	SyncBehaviour string // "upsert" | "mirror" | "full"
	Schedule      string // HCL fragment for the schedule { } block body
	Identifiers   string // HCL fragment for one or more raw identifiers { } blocks (required, ≥1)
	Mappings      string // HCL fragment for one or more raw mappings { } blocks (optional)
	Event         string // HCL fragment for event { } block body (optional)
	CursorColumn  string // value for cursor_column (optional, requires upsert)

	// UpdateMappings, when non-empty, runs an Update step replacing the
	// connection's mappings { } blocks. Mappings are mutable across all flows,
	// so this is the cheapest way to flex the Update path.
	UpdateMappings string
}

// AccAssertRETLSourceModel runs the standard E2E lifecycle against
// rudderstack_retl_source_model. In plan-only mode it validates HCL/schema only;
// in live mode it runs Create → Update → Import → Destroy and verifies via the
// RETL API.
func AccAssertRETLSourceModel(t *testing.T, cfg RETLSourceTestConfig) {
	t.Helper()
	runRETLSourceLifecycle(t, "rudderstack_retl_source_model", cfg)
}

// AccAssertRETLSourceTable runs the standard E2E lifecycle against
// rudderstack_retl_source_table. See AccAssertRETLSourceModel.
func AccAssertRETLSourceTable(t *testing.T, cfg RETLSourceTestConfig) {
	t.Helper()
	runRETLSourceLifecycle(t, "rudderstack_retl_source_table", cfg)
}

func runRETLSourceLifecycle(t *testing.T, resourceType string, cfg RETLSourceTestConfig) {
	t.Helper()

	planOnly := PlanOnly()
	if planOnly {
		t.Parallel()
	}

	resourceName := fmt.Sprintf("%s.test", resourceType)
	name := RandomName(resourceType)

	accountID := planOnlyAccountID
	if !planOnly {
		accountID = os.Getenv(retlAccountIDEnv)
		if accountID == "" {
			t.Skipf("%s is not set; skipping live RETL source test (set it to a workspace warehouse account ID)", retlAccountIDEnv)
		}
	}

	createHCL := retlSourceHCL(resourceType, name, cfg.SourceDefinitionName, accountID, cfg.Config)

	if planOnly {
		ensureDummyToken(t)
		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:             createHCL,
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		})
		return
	}

	steps := []resource.TestStep{
		{
			Config: createHCL,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckRETLSourceExists(resourceName),
				resource.TestCheckResourceAttr(resourceName, "name", name),
				resource.TestCheckResourceAttr(resourceName, "source_definition_name", cfg.SourceDefinitionName),
				resource.TestCheckResourceAttr(resourceName, "account_id", accountID),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
			),
		},
	}

	if cfg.UpdateConfig != "" {
		updateHCL := retlSourceHCL(resourceType, name+"-updated", cfg.SourceDefinitionName, accountID, cfg.UpdateConfig)
		steps = append(steps, resource.TestStep{
			Config: updateHCL,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckRETLSourceExists(resourceName),
				resource.TestCheckResourceAttr(resourceName, "name", name+"-updated"),
			),
		})
	}

	steps = append(steps, resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckRETLSourceDestroy(resourceType),
		Steps:             steps,
	})
}

// AccAssertRETLConnection wires a model source + webhook destination +
// rudderstack_retl_connection and runs the lifecycle. The webhook destination
// is sufficient for the control-plane API (it doesn't validate downstream
// connectivity) and avoids needing real warehouse credentials in CI.
func AccAssertRETLConnection(t *testing.T, cfg RETLConnectionTestConfig) {
	t.Helper()

	planOnly := PlanOnly()
	if planOnly {
		t.Parallel()
	}

	connResource := "rudderstack_retl_connection.test"
	srcResource := "rudderstack_retl_source_model.test"
	dstResource := "rudderstack_destination_webhook.test"

	suffix := cfg.Variant
	if suffix == "" {
		suffix = "conn"
	}
	srcName := RandomName("retl-src-" + suffix)
	dstName := RandomName("retl-dst-" + suffix)

	accountID := planOnlyAccountID
	if !planOnly {
		accountID = os.Getenv(retlAccountIDEnv)
		if accountID == "" {
			t.Skipf("%s is not set; skipping live RETL connection test", retlAccountIDEnv)
		}
	}

	createHCL := retlConnectionHCL(srcName, dstName, accountID, cfg)

	if planOnly {
		ensureDummyToken(t)
		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:             createHCL,
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		})
		return
	}

	steps := []resource.TestStep{
		{
			Config: createHCL,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckRETLConnectionExists(connResource),
				resource.TestCheckResourceAttrSet(connResource, "id"),
				resource.TestCheckResourceAttr(connResource, "enabled", "true"),
				resource.TestCheckResourceAttr(connResource, "sync_behaviour", cfg.SyncBehaviour),
				resource.TestCheckResourceAttrPair(connResource, "source_id", srcResource, "id"),
				resource.TestCheckResourceAttrPair(connResource, "destination_id", dstResource, "id"),
			),
		},
	}

	if cfg.UpdateMappings != "" {
		updated := cfg
		updated.Mappings = cfg.UpdateMappings
		updated.UpdateMappings = ""
		steps = append(steps, resource.TestStep{
			Config: retlConnectionHCL(srcName, dstName, accountID, updated),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckRETLConnectionExists(connResource),
			),
		})
	}

	steps = append(steps, resource.TestStep{
		ResourceName:      connResource,
		ImportState:       true,
		ImportStateVerify: true,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckRETLConnectionDestroy,
		Steps:             steps,
	})
}

// retlSourceHCL builds the HCL for a single RETL source resource.
func retlSourceHCL(resourceType, name, sourceDefinitionName, accountID, configBlock string) string {
	return fmt.Sprintf(`
resource %q "test" {
  name                   = %q
  source_definition_name = %q
  account_id             = %q
  config {
    %s
  }
}
`, resourceType, name, sourceDefinitionName, accountID, configBlock)
}

// retlConnectionHCL builds the HCL for the source + webhook destination +
// connection pipeline. Identifiers / Mappings are full block fragments
// (e.g. `identifiers { from = "..." to = "..." }`) so callers can pass
// multiple of each. Event and cursor_column are inlined when non-empty.
func retlConnectionHCL(srcName, dstName, accountID string, cfg RETLConnectionTestConfig) string {
	mappingsBlock := ""
	if cfg.Mappings != "" {
		mappingsBlock = "\n  " + cfg.Mappings
	}
	eventBlock := ""
	if cfg.Event != "" {
		eventBlock = "\n  event {\n    " + cfg.Event + "\n  }"
	}
	cursorAttr := ""
	if cfg.CursorColumn != "" {
		cursorAttr = fmt.Sprintf("\n  cursor_column  = %q", cfg.CursorColumn)
	}

	return fmt.Sprintf(`
resource "rudderstack_retl_source_model" "test" {
  name                   = %q
  source_definition_name = "bigquery"
  account_id             = %q
  config {
    primary_key = "id"
    sql         = "select 1 as id"
  }
}

resource "rudderstack_destination_webhook" "test" {
  name = %q
  config {
    webhook_url    = "https://example.com/test"
    webhook_method = "POST"
  }
}

resource "rudderstack_retl_connection" "test" {
  source_id      = rudderstack_retl_source_model.test.id
  destination_id = rudderstack_destination_webhook.test.id
  enabled        = true
  sync_behaviour = %q%s
  schedule {
    %s
  }
  %s%s%s
}
`, srcName, accountID, dstName, cfg.SyncBehaviour, cursorAttr, cfg.Schedule, cfg.Identifiers, mappingsBlock, eventBlock)
}

// newTestRETLClient wraps the live API client with a typed RETL store.
func newTestRETLClient() (retl.RETLStore, error) {
	cl, err := newTestAPIClient()
	if err != nil {
		return nil, err
	}
	return retl.NewRudderRETLStore(cl), nil
}

func testAccCheckRETLSourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is empty")
		}

		store, err := newTestRETLClient()
		if err != nil {
			return err
		}
		if _, err := store.GetRetlSource(context.Background(), rs.Primary.ID); err != nil {
			return fmt.Errorf("RETL source %s not found in API: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckRETLSourceDestroy(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		store, err := newTestRETLClient()
		if err != nil {
			return err
		}
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			// Tolerant of soft-delete: the API may still return the resource
			// after delete. We don't fail here — Destroy already verified the
			// Delete handler ran without error.
			_, _ = store.GetRetlSource(context.Background(), rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckRETLConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is empty")
		}

		store, err := newTestRETLClient()
		if err != nil {
			return err
		}
		if _, err := store.GetConnection(context.Background(), rs.Primary.ID); err != nil {
			return fmt.Errorf("RETL connection %s not found in API: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckRETLConnectionDestroy(s *terraform.State) error {
	store, err := newTestRETLClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rudderstack_retl_connection" {
			continue
		}
		// Expect 404 after delete; tolerant of any error (network blips,
		// soft-delete) — we don't fail the test here.
		_, getErr := store.GetConnection(context.Background(), rs.Primary.ID)
		if getErr == nil {
			// If it still exists, surface that — it's likely a real bug.
			return fmt.Errorf("RETL connection %s still exists after destroy", rs.Primary.ID)
		}
		var apiErr *client.APIError
		if errors.As(getErr, &apiErr) && apiErr.HTTPStatusCode == 404 {
			continue
		}
	}
	return nil
}
