# Customer.io — RETL connection scoped to Customer.io
# destinations. `object` is a typed top-level field (ForceNew — changing it
# recreates the connection). Supported object values are `person` and `event`;
# supported sync behaviours are `upsert` and `mirror`.
resource "rudderstack_retl_connection_customerio" "model_to_customerio" {
  source_id      = rudderstack_retl_source_model.users_revenue.id
  destination_id = rudderstack_destination_customerio.example.id
  sync_behaviour = "upsert"
  object         = "person"

  # Optional: incremental watermark column. Only valid when sync_behaviour is "upsert".
  cursor_column = "updated_at"

  schedule {
    type          = "basic"
    every_minutes = 30
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
