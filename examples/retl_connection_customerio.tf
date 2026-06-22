# Customer.io (VDM v2) flow — destination-specific RETL connection scoped to
# Customer.io destinations. `object` is a typed top-level field (ForceNew
# because the destinationConfig shape is not mutable in place on this flow —
# changing it recreates the connection). identifiers map to VDM v2
# identifierMappings and mappings to fieldMappings; config-be assembles the
# VDM v2 connectionConfig server-side from the Customer.io destination
# definition. Only `upsert` and `mirror` sync behaviours are supported.
resource "rudderstack_retl_connection_customerio" "model_to_customerio" {
  source_id      = rudderstack_retl_source_model.users_revenue.id
  destination_id = rudderstack_destination_customerio.example.id
  sync_behaviour = "upsert"
  object         = "customers"

  # Optional: incremental watermark column. Only valid when sync_behaviour
  # is "upsert". Sent as a generic top-level source field (not destinationConfig).
  cursor_column = "updated_at"

  schedule {
    type          = "basic"
    every_minutes = 30
  }

  identifiers {
    from = "email"
    to   = "email"
  }

  mappings {
    from = "name"
    to   = "plan"
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
