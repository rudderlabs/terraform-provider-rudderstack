# JSON Mapper flow — identifiers + per-field mappings, no `object`.
# `basic` schedule runs every N minutes. `constants` adds static key/value
# pairs to every synced row (mutable in JSON Mapper; ForceNew in Object
# Mapping). Destination-specific flows (e.g. Customer.io Audience) have their
# own typed resources — see retl_connection_customerio_audience.tf.
resource "rudderstack_retl_connection" "table_to_webhook" {
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_webhook.example.id
  enabled        = true
  sync_behaviour = "full"

  schedule {
    type          = "basic"
    every_minutes = 60
  }

  event {
    type = "identify"
  }

  identifiers {
    from = "email"
    to   = "user_id"
  }

  mappings {
    from = "created_at"
    to   = "traits.obj.created_at"
  }

  constants {
    key   = "context.source"
    value = "warehouse-users"
  }

  constants {
    key   = "context.env"
    value = "production"
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

# Object Mapping flow — `object` selects the destination-side object (e.g.
# "customers" in Customer.io). `cron` schedule fires on a cron expression.
resource "rudderstack_retl_connection" "table_to_customerio" {
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_customerio.example.id
  enabled        = true
  sync_behaviour = "upsert"
  object         = "customers"

  schedule {
    type            = "cron"
    cron_expression = "0 */6 * * *"
  }

  identifiers {
    from = "email"
    to   = "email"
  }

  mappings {
    from = "user_id"
    to   = "id"
  }

  mappings {
    from = "created_at"
    to   = "created_at"
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

# Incremental upsert — `cursor_column` makes each sync read only rows whose
# cursor value is newer than the previous run. Only valid when
# `sync_behaviour = "upsert"`. ForceNew: changes recreate the connection.
resource "rudderstack_retl_connection" "table_to_webhook_incremental" {
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_webhook.example.id
  enabled        = true
  sync_behaviour = "upsert"
  cursor_column  = "updated_at"

  schedule {
    type          = "basic"
    every_minutes = 15
  }

  event {
    type = "identify"
  }

  identifiers {
    from = "user_id"
    to   = "user_id"
  }

  mappings {
    from = "email"
    to   = "traits.email"
  }

  mappings {
    from = "updated_at"
    to   = "traits.updated_at"
  }

  sync_settings {
    sync_logs_config {
      enabled               = true
      log_retention_in_days = 30
      snapshots_to_retain   = 5
    }
    failed_keys_config {
      enable_failed_keys_retry = true
    }
  }
}

# JSON Mapper with `track` events — each synced row becomes a track call with
# a fixed event `name`. Use `name_column` instead if the event name should be
# read from a column in the source.
resource "rudderstack_retl_connection" "track_to_webhook" {
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_webhook.example.id
  enabled        = true
  sync_behaviour = "full"

  schedule {
    type = "manual"
  }

  event {
    type = "track"
    name = "user_synced"
  }

  identifiers {
    from = "user_id"
    to   = "userId"
  }

  mappings {
    from = "email"
    to   = "properties.email"
  }

  mappings {
    from = "plan"
    to   = "properties.plan"
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

# JSON Mapper with `track` events, dynamic event name — `name_column` reads
# the event name from a column in the source row, so each row can emit a
# different track event. Mutually exclusive with `name`.
resource "rudderstack_retl_connection" "track_to_webhook_named_column" {
  source_id      = rudderstack_retl_source_table.users.id
  destination_id = rudderstack_destination_webhook.example.id
  enabled        = true
  sync_behaviour = "full"

  schedule {
    type = "manual"
  }

  event {
    type        = "track"
    name_column = "event_name"
  }

  identifiers {
    from = "user_id"
    to   = "userId"
  }

  mappings {
    from = "email"
    to   = "properties.email"
  }

  mappings {
    from = "plan"
    to   = "properties.plan"
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

