# Customer.io — RETL connection scoped to a Customer.io person object.
# `object` is a typed top-level field (ForceNew — changing it recreates the
# connection). For person syncs, `sync_behaviour` is required and must be either
# "upsert" or "mirror".
resource "rudderstack_retl_connection_customerio" "model_to_customerio_person" {
  source_id      = rudderstack_retl_source_model.users_revenue.id
  destination_id = rudderstack_destination_customerio.example.id
  sync_behaviour = "upsert"
  object         = "person"

  # Optional: incremental watermark column. Only valid when sync_behaviour is
  # "upsert".
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

# For event objects the sync mode is determined by the backend, so
# `sync_behaviour` must be omitted here; Terraform stores whatever value the
# backend returns in state after apply.
resource "rudderstack_retl_connection_customerio" "model_to_customerio_event" {
  source_id      = rudderstack_retl_source_model.users_revenue.id
  destination_id = rudderstack_destination_customerio.example.id
  object         = "event"

  schedule {
    type          = "basic"
    every_minutes = 30
  }

  identifiers {
    from = "email"
    to   = "email"
  }
}
