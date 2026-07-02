# Shared base: a BigQuery rETL account + an rETL source table reading from it.
# Every scenario stands up its own copy (isolated create/destroy) by calling this
# module.
#
# The required_providers block maps the local name `rudderstack` to the
# rudderstack.com/rudderlabs/rudderstack source so it matches the root's
# dev_override (without it, Terraform infers hashicorp/rudderstack for the module
# and `init` hits the registry). NO `provider` block here — provider config is
# inherited from the root, keeping init offline for the locally-built provider.
terraform {
  required_providers {
    rudderstack = { source = "rudderstack.com/rudderlabs/rudderstack" }
  }
}

# 1. BigQuery rETL account — stores warehouse credentials.
resource "rudderstack_account_source_bigquery" "acct" {
  name = "${var.name_prefix}-bigquery"
  config {
    project     = var.bq_project
    location    = var.bq_location
    credentials = var.bq_credentials
  }
}

# 2. rETL source that reads from a single BigQuery table.
resource "rudderstack_retl_source_table" "users" {
  name                   = "${var.name_prefix}-users-table"
  source_definition_name = "bigquery"
  account_id             = rudderstack_account_source_bigquery.acct.id
  enabled                = true
  config {
    primary_key = "user_id"
    schema      = var.bq_dataset
    table       = var.bq_table
  }
}
