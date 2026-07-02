# Scenario: BigQuery → Customer.io (VDM v2 / destination-specific flow).
#
# terraform{}, provider, and all variable declarations come from the symlinked
# _shared.tf. This file holds only the scenario-specific resources + the four
# standardized outputs (account_id, source_id, destination_id, connection_id)
# that run.sh verifies for every scenario.

module "bq" {
  source         = "../../modules/bigquery_source"
  name_prefix    = "tf-e2e-customerio"
  bq_project     = var.bq_project
  bq_location    = var.bq_location
  bq_dataset     = var.bq_dataset
  bq_table       = var.bq_table
  bq_credentials = var.bq_credentials
}

resource "rudderstack_destination_customerio" "cio" {
  name = "tf-e2e-customerio"
  config {
    site_id    = var.customerio_site_id
    api_key    = var.customerio_api_key
    datacenter = var.customerio_datacenter
  }
}

# Customer.io is a VDM v2 / "destination-specific flow": the object must travel
# inside the API's destinationConfig as {"object":"person"}, which only the typed
# resource packs. manual schedule so no syncs fire.
resource "rudderstack_retl_connection_customerio" "to_customerio" {
  source_id      = module.bq.source_id
  destination_id = rudderstack_destination_customerio.cio.id
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

output "account_id" {
  description = "ID of the created BigQuery rETL account."
  value       = module.bq.account_id
}

output "source_id" {
  description = "ID of the created rETL source table."
  value       = module.bq.source_id
}

output "destination_id" {
  description = "ID of the created Customer.io destination."
  value       = rudderstack_destination_customerio.cio.id
}

output "connection_id" {
  description = "ID of the created BigQuery→Customer.io rETL connection."
  value       = rudderstack_retl_connection_customerio.to_customerio.id
}
