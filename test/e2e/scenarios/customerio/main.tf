# Scenario: BigQuery → Customer.io (VDM v2 / destination-specific flow).
#
# terraform{}, provider, and all variable declarations come from the symlinked
# _shared.tf. This file holds only the scenario-specific resources + the four
# standardized outputs (account_id, source_id, destination_id, connection_id)
# that run.sh verifies for every scenario.

locals {
  base_name = var.run_id == "" ? "tf-e2e-customerio" : "tf-e2e-customerio-${var.run_id}"
}

module "bq" {
  source         = "../../modules/bigquery_source"
  name_prefix    = local.base_name
  bq_project     = var.bq_project
  bq_location    = var.bq_location
  bq_dataset     = var.bq_dataset
  bq_table       = var.bq_table
  bq_credentials = var.bq_credentials
}

resource "rudderstack_destination_customerio" "cio" {
  name = local.base_name
  config {
    site_id         = var.customerio_site_id
    api_key         = var.customerio_api_key
    datacenter      = var.customerio_datacenter
    api_version     = "v2"
    user_id_mapping = "id"
  }
}

# Customer.io is a VDM v2 / "destination-specific flow": the object must travel
# inside the API's destinationConfig as {"object":"..."}, which only the typed
# resource packs. Both supported objects are exercised against the same
# destination — manual schedule so no syncs fire.
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

# object = "event". Unlike `person`, the event sync mode is chosen by the
# backend, so `sync_behaviour` is omitted from config and stamped into state on
# apply. This connection exists to put that claim in front of the real API: the
# create must be accepted without `syncBehaviour`, and the drift assertion
# run.sh performs after apply (terraform plan -detailed-exitcode) must stay at
# exit 0, which is what proves the Optional+Computed round-trip does not produce
# a perpetual diff. Unit tests can only assert this against a mock that returns
# whatever the fixtures say.
resource "rudderstack_retl_connection_customerio" "to_customerio_event" {
  source_id      = module.bq.source_id
  destination_id = rudderstack_destination_customerio.cio.id
  object         = "event"

  schedule {
    type = "manual"
  }

  identifiers {
    from = "email"
    to   = "email"
  }

  # Event objects additionally require a `name` identifier carrying the event
  # name — the API rejects the connection without it. The schedule is manual so
  # no sync ever runs; the mapping only has to name a column that really exists
  # on the source table, and `user_id` is its primary key (see
  # modules/bigquery_source).
  identifiers {
    from = "user_id"
    to   = "name"
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
  description = "ID of the created BigQuery→Customer.io rETL connection (object = person)."
  value       = rudderstack_retl_connection_customerio.to_customerio.id
}

output "event_connection_id" {
  description = "ID of the created BigQuery→Customer.io rETL connection (object = event)."
  value       = rudderstack_retl_connection_customerio.to_customerio_event.id
}
