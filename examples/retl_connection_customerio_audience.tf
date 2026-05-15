# Customer.io Audience flow — destination-specific RETL connection scoped
# to Customer.io Audience destinations. `audience_id` is a typed top-level
# field (ForceNew because the Customer.io Audience API does not accept
# destinationConfig changes on update — bumping audience_id recreates the
# connection). `manual` schedule only runs when triggered explicitly.
resource "rudderstack_retl_connection_customerio_audience" "model_to_customerio_audience" {
  source_id      = rudderstack_retl_source_model.users_revenue.id
  destination_id = rudderstack_destination_customerio_audience.example.id
  sync_behaviour = "mirror"
  audience_id    = 16

  schedule {
    type = "manual"
  }

  identifiers {
    from = "email"
    to   = "email"
  }

  sync_settings {
    sync_logs_config {
      enabled               = true
      log_retention_in_days = 30
      snapshots_to_retain   = 5
    }
    failed_keys_config {
      enable_failed_keys_retry = false
    }
  }
}
