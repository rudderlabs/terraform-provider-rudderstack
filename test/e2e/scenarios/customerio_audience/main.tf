# Scenario: BigQuery → Customer.io Audience.
#
# terraform{}, provider, and all variable declarations come from the symlinked
# _shared.tf. This file holds only the scenario-specific resources + the four
# standardized outputs (account_id, source_id, destination_id, connection_id).

module "bq" {
  source         = "../../modules/bigquery_source"
  name_prefix    = "tf-e2e-customerio-audience"
  bq_project     = var.bq_project
  bq_location    = var.bq_location
  bq_dataset     = var.bq_dataset
  bq_table       = var.bq_table
  bq_credentials = var.bq_credentials
}

resource "rudderstack_destination_customerio_audience" "cio_aud" {
  name = "tf-e2e-customerio-audience"
  config {
    site_id     = var.customerio_site_id
    api_key     = var.customerio_api_key
    app_api_key = var.customerio_audience_app_api_key
    region      = var.customerio_audience_region
  }
}

# Typed resource carrying audience_id as a top-level field (packed into the API's
# destinationConfig). manual schedule so no syncs fire.
resource "rudderstack_retl_connection_customerio_audience" "to_customerio_audience" {
  source_id      = module.bq.source_id
  destination_id = rudderstack_destination_customerio_audience.cio_aud.id
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

output "account_id" {
  description = "ID of the created BigQuery rETL account."
  value       = module.bq.account_id
}

output "source_id" {
  description = "ID of the created rETL source table."
  value       = module.bq.source_id
}

output "destination_id" {
  description = "ID of the created Customer.io Audience destination."
  value       = rudderstack_destination_customerio_audience.cio_aud.id
}

output "connection_id" {
  description = "ID of the created BigQuery→Customer.io Audience rETL connection."
  value       = rudderstack_retl_connection_customerio_audience.to_customerio_audience.id
}
