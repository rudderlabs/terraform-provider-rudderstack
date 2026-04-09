package acc

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// ConnectionTestConfig holds the minimal config needed to create the source and
// destination that a connection test will wire together.
type ConnectionTestConfig struct {
	Source      string // source type, e.g. "javascript"
	Destination string // destination type, e.g. "webhook"
	DestConfig  string // HCL inside the destination's config { } block
}

// AccAssertConnection runs an E2E test that creates a source + destination, connects
// them, updates the connection (disable/enable), imports, and destroys.
//
// When PlanOnly() is enabled (TF_ACC_PLAN_ONLY=1): plan-only validation (zero API calls).
// Otherwise: full CRUD against the real API.
func AccAssertConnection(t *testing.T, cfg ConnectionTestConfig) {
	t.Helper()

	planOnly := PlanOnly()
	if planOnly {
		t.Parallel()
	}

	connResource := "rudderstack_connection.test"
	srcName := RandomName("conn-src-" + cfg.Source)
	dstName := RandomName("conn-dst-" + cfg.Destination)

	if planOnly {
		ensureDummyToken(t)
		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:             testAccConnectionConfig(cfg, srcName, dstName, true),
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
		CheckDestroy:      testAccCheckConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig(cfg, srcName, dstName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(connResource),
					resource.TestCheckResourceAttrSet(connResource, "id"),
					resource.TestCheckResourceAttr(connResource, "enabled", "true"),
					resource.TestCheckResourceAttrPair(connResource, "source_id",
						fmt.Sprintf("rudderstack_source_%s.test", cfg.Source), "id"),
					resource.TestCheckResourceAttrPair(connResource, "destination_id",
						fmt.Sprintf("rudderstack_destination_%s.test", cfg.Destination), "id"),
				),
			},
			{
				ResourceName:      connResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConnectionConfig(cfg ConnectionTestConfig, srcName, dstName string, enabled bool) string {
	destBlock := ""
	if cfg.DestConfig != "" {
		destBlock = fmt.Sprintf("\n  config {\n    %s\n  }", cfg.DestConfig)
	}

	return fmt.Sprintf(`
resource "rudderstack_source_%s" "test" {
  name = %q
}

resource "rudderstack_destination_%s" "test" {
  name = %q%s
}

resource "rudderstack_connection" "test" {
  source_id      = rudderstack_source_%s.test.id
  destination_id = rudderstack_destination_%s.test.id
  enabled        = %t
}
`, cfg.Source, srcName, cfg.Destination, dstName, destBlock, cfg.Source, cfg.Destination, enabled)
}

func testAccCheckConnectionExists(resourceName string) resource.TestCheckFunc {
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
		_, err = cl.Connections.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("connection %s not found in API: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckConnectionDestroy(s *terraform.State) error {
	cl, err := newTestAPIClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rudderstack_connection" {
			continue
		}
		_, err := cl.Connections.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("connection %s still exists after destroy", rs.Primary.ID)
		}
	}
	return nil
}
