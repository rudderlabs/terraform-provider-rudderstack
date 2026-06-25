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

# 3. (optional) Customer.io destination + rETL connection from the BigQuery source.
#    BigQuery source. ponytail: count-gated on creds so the webhook-only smoke
#    still runs when Customer.io creds aren't supplied.
locals {
  enable_customerio = var.customerio_api_key != "" && var.customerio_site_id != ""
}

resource "rudderstack_destination_customerio" "cio" {
  count = local.enable_customerio ? 1 : 0
  name  = "tf-e2e-customerio"
  config {
    site_id    = var.customerio_site_id
    api_key    = var.customerio_api_key
    datacenter = var.customerio_datacenter
  }
}

# Customer.io rETL connection via the GENERIC rudderstack_retl_connection in
# object-mapping mode: object="person" selects the Customer.io object — same
# {"object":"person"} payload the typed resource sends. manual schedule so no
# syncs fire during the smoke.
resource "rudderstack_retl_connection" "to_customerio" {
  count          = local.enable_customerio ? 1 : 0
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_customerio.cio[0].id
  enabled        = true
  sync_behaviour = "mirror"
  object         = "person"

  schedule {
    type = "manual"
  }

  identifiers {
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
  description = "ID of the created Customer.io destination (empty when creds not supplied)."
  value       = try(rudderstack_destination_customerio.cio[0].id, "")
}

output "connection_id" {
  description = "ID of the created rETL connection (empty when creds not supplied)."
  value       = try(rudderstack_retl_connection.to_customerio[0].id, "")
}

output "customerio_destination_id" {
  description = "ID of the Customer.io destination (empty when creds not supplied)."
  value       = try(rudderstack_destination_customerio.cio[0].id, "")
}

output "customerio_connection_id" {
  description = "ID of the BigQuery→Customer.io rETL connection (empty when creds not supplied)."
  value       = try(rudderstack_retl_connection.to_customerio[0].id, "")
}
