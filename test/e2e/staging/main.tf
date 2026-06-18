terraform {
  required_providers {
    rudderstack = { source = "rudderstack.com/rudderlabs/rudderstack" }
  }
}

provider "rudderstack" {
  api_url      = var.api_url
  access_token = var.access_token
}

# 1. BigQuery rETL account — stores warehouse credentials.
resource "rudderstack_account_source_bigquery" "acct" {
  name = "tf-e2e-bigquery"
  config {
    project     = var.bq_project
    location    = var.bq_location
    credentials = var.bq_credentials
  }
}

# 2. rETL source that reads from a single BigQuery table.
resource "rudderstack_retl_source_table" "users" {
  name                   = "tf-e2e-users-table"
  source_definition_name = "bigquery"
  account_id             = rudderstack_account_source_bigquery.acct.id
  enabled                = true
  config {
    primary_key = "user_id"
    schema      = var.bq_dataset
    table       = var.bq_table
  }
}

# 3. Webhook destination — throwaway endpoint; delivery is not the point here.
resource "rudderstack_destination_webhook" "demo" {
  name = "tf-e2e-webhook"
  config {
    webhook_url    = "https://example.com/test"
    webhook_method = "POST"
  }
}

# 4. rETL connection wiring the source to the destination.
#    Uses a manual schedule so no automatic syncs are triggered during the test.
resource "rudderstack_retl_connection" "to_webhook" {
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_webhook.demo.id
  enabled        = true
  sync_behaviour = "full"

  schedule {
    type = "manual"
  }

  event {
    type = "identify"
  }

  identifiers {
    from = "user_id"
    to   = "user_id"
  }

  mappings {
    from = "email"
    to   = "email"
  }
}

# Outputs — consumed by run.sh to verify IDs were created.
output "account_id" {
  description = "ID of the created BigQuery rETL account."
  value       = rudderstack_account_source_bigquery.acct.id
}

output "retl_source_id" {
  description = "ID of the created rETL source table."
  value       = rudderstack_retl_source_table.users.id
}

output "destination_id" {
  description = "ID of the created webhook destination."
  value       = rudderstack_destination_webhook.demo.id
}

output "connection_id" {
  description = "ID of the created rETL connection."
  value       = rudderstack_retl_connection.to_webhook.id
}
