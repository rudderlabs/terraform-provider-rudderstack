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

  # Independent of the VDM v2 chain above: the Audience chain needs its own App
  # API key plus a real audience ID (on top of the shared site_id/api_key).
  enable_customerio_audience = var.customerio_site_id != "" && var.customerio_api_key != "" && var.customerio_audience_app_api_key != "" && var.customerio_audience_id > 0
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

# Customer.io rETL connection. Customer.io is a VDM v2 / "destination-specific
# flow": the object must travel inside the API's destinationConfig as
# {"object":"person"}, which only the typed resource packs. The generic
# rudderstack_retl_connection sends object as a top-level field, which the API
# rejects ("Fields not allowed for Destination-specific flow: object"), so this
# must use the typed resource (added in #275). manual schedule so no syncs fire.
resource "rudderstack_retl_connection_customerio" "to_customerio" {
  count          = local.enable_customerio ? 1 : 0
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_customerio.cio[0].id
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

# 4. (optional) Customer.io Audience destination + rETL connection from the same
#    BigQuery source. Mirrors the Customer.io (VDM v2) chain above but targets the
#    Audience destination via the typed rudderstack_retl_connection_customerio_audience
#    resource, which carries audience_id as a top-level field (packed into the API's
#    destinationConfig). count-gated on its own creds so it's independent of the
#    VDM v2 chain. manual schedule so no syncs fire.
resource "rudderstack_destination_customerio_audience" "cio_aud" {
  count = local.enable_customerio_audience ? 1 : 0
  name  = "tf-e2e-customerio-audience"
  config {
    site_id     = var.customerio_site_id
    api_key     = var.customerio_api_key
    app_api_key = var.customerio_audience_app_api_key
    region      = var.customerio_audience_region
  }
}

resource "rudderstack_retl_connection_customerio_audience" "to_customerio_audience" {
  count          = local.enable_customerio_audience ? 1 : 0
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_customerio_audience.cio_aud[0].id
  sync_behaviour = "mirror"
  audience_id    = var.customerio_audience_id

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
  value       = try(rudderstack_retl_connection_customerio.to_customerio[0].id, "")
}

output "customerio_destination_id" {
  description = "ID of the Customer.io destination (empty when creds not supplied)."
  value       = try(rudderstack_destination_customerio.cio[0].id, "")
}

output "customerio_connection_id" {
  description = "ID of the BigQuery→Customer.io rETL connection (empty when creds not supplied)."
  value       = try(rudderstack_retl_connection_customerio.to_customerio[0].id, "")
}

output "customerio_audience_destination_id" {
  description = "ID of the Customer.io Audience destination (empty when audience creds not supplied)."
  value       = try(rudderstack_destination_customerio_audience.cio_aud[0].id, "")
}

output "customerio_audience_connection_id" {
  description = "ID of the BigQuery→Customer.io Audience rETL connection (empty when audience creds not supplied)."
  value       = try(rudderstack_retl_connection_customerio_audience.to_customerio_audience[0].id, "")
}
