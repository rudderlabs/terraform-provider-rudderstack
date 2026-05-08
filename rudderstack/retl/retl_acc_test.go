package retl_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
)

// TestAccRETLSourceModel_Snowflake exercises rudderstack_retl_source_model
// against a snowflake account. SQL is intentionally trivial; the API accepts
// it at create time and any sync-time SQL validation is out of scope here.
func TestAccRETLSourceModel_Snowflake(t *testing.T) {
	acc.AccAssertRETLSourceModel(t, acc.RETLSourceTestConfig{
		SourceDefinitionName: "snowflake",
		Config: `
			primary_key = "id"
			sql         = "select 1 as id"
		`,
		UpdateConfig: `
			primary_key = "id"
			sql         = "select 1 as id, 'two' as name"
			description = "v2"
		`,
	})
}

// TestAccRETLSourceTable_Snowflake exercises rudderstack_retl_source_table.
func TestAccRETLSourceTable_Snowflake(t *testing.T) {
	acc.AccAssertRETLSourceTable(t, acc.RETLSourceTestConfig{
		SourceDefinitionName: "snowflake",
		Config: `
			primary_key = "id"
			schema      = "public"
			table       = "users"
		`,
		UpdateConfig: `
			primary_key = "id"
			schema      = "public"
			table       = "events"
		`,
	})
}

// TestAccRETLConnection_JSONMapperBasicSchedule covers the JSON-mapper flow
// (no `object`, no `destination_config`) with a basic-interval schedule and an
// identify event. Mappings update verifies the mutable Update path.
func TestAccRETLConnection_JSONMapperBasicSchedule(t *testing.T) {
	acc.AccAssertRETLConnection(t, acc.RETLConnectionTestConfig{
		Variant:       "jm-basic",
		SyncBehaviour: "upsert",
		Schedule: `
			type          = "basic"
			every_minutes = 60
		`,
		Identifiers: `
			from = "email"
			to   = "user_id"
		`,
		Mappings: `
			mappings {
				from = "name"
				to   = "first_name"
			}
		`,
		Event: `type = "identify"`,
		UpdateMappings: `
			mappings {
				from = "name"
				to   = "first_name"
			}
			mappings {
				from = "phone"
				to   = "phone_number"
			}
		`,
	})
}

// TestAccRETLConnection_CronSchedule covers the cron schedule branch.
func TestAccRETLConnection_CronSchedule(t *testing.T) {
	acc.AccAssertRETLConnection(t, acc.RETLConnectionTestConfig{
		Variant:       "cron",
		SyncBehaviour: "upsert",
		Schedule: `
			type            = "cron"
			cron_expression = "30 13 * * *"
		`,
		Identifiers: `
			from = "email"
			to   = "user_id"
		`,
		Event: `type = "identify"`,
	})
}

// TestAccRETLConnection_ManualSchedule covers the manual schedule branch.
// `mirror` sync_behaviour is used here to flex a non-upsert path.
func TestAccRETLConnection_ManualSchedule(t *testing.T) {
	acc.AccAssertRETLConnection(t, acc.RETLConnectionTestConfig{
		Variant:       "manual",
		SyncBehaviour: "mirror",
		Schedule:      `type = "manual"`,
		Identifiers: `
			from = "email"
			to   = "user_id"
		`,
		Event: `type = "identify"`,
	})
}
