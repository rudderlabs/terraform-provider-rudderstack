---
page_title: "rudderstack_retl_connection Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-

---

# rudderstack_retl_connection (Resource)

A generic RETL (Reverse ETL) connection between a RETL source and a destination, covering JSON Mapper and Object Mapping flows. Destination-specific flows (e.g. Customer.io Audience) have their own typed resources.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination_id` (String) ID of the destination.
- `identifiers` (Block List, Min: 1) Source-to-destination identifier mappings (mutable). (see [below for nested schema](#nestedblock--identifiers))
- `schedule` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--schedule))
- `source_id` (String) ID of the RETL source.
- `sync_behaviour` (String) How records are synced to the destination: `upsert`, `mirror`, or `full`.

### Optional

- `constants` (Block List) User-defined constants. Mutable for JSON Mapper; ForceNew for Object Mapping. (see [below for nested schema](#nestedblock--constants))
- `cursor_column` (String) Column name for incremental upsert syncs (only valid when sync_behaviour is `upsert`).
- `enabled` (Boolean) Whether the connection is enabled.
- `event` (Block List, Max: 1) CDP event configuration. Optional in the Terraform schema; flow-specific requirements (required for JSON Mapper, absent for Object Mapping) are enforced by the API. (see [below for nested schema](#nestedblock--event))
- `mappings` (Block List) Source-to-destination field mappings (mutable). (see [below for nested schema](#nestedblock--mappings))
- `object` (String) Destination entity for Object Mapping flows (e.g. `Contact`, `Lead`).
- `sync_settings` (Block List, Max: 1) (see [below for nested schema](#nestedblock--sync_settings))

### Read-Only

- `created_at` (String)
- `id` (String) The ID of this resource.
- `updated_at` (String)

<a id="nestedblock--identifiers"></a>
### Nested Schema for `identifiers`

Required:

- `from` (String)
- `to` (String)


<a id="nestedblock--schedule"></a>
### Nested Schema for `schedule`

Required:

- `type` (String) Schedule type: `basic`, `manual`, or `cron`.

Optional:

- `cron_expression` (String) Cron expression. Required when `type` is `cron`.
- `every_minutes` (Number) Sync interval in minutes. Required when `type` is `basic`.


<a id="nestedblock--constants"></a>
### Nested Schema for `constants`

Required:

- `key` (String)
- `value` (String)


<a id="nestedblock--event"></a>
### Nested Schema for `event`

Required:

- `type` (String)

Optional:

- `name` (String)
- `name_column` (String)


<a id="nestedblock--mappings"></a>
### Nested Schema for `mappings`

Required:

- `from` (String)
- `to` (String)


<a id="nestedblock--sync_settings"></a>
### Nested Schema for `sync_settings`

Optional:

- `failed_keys_config` (Block List, Max: 1) (see [below for nested schema](#nestedblock--sync_settings--failed_keys_config))
- `sync_logs_config` (Block List, Max: 1) (see [below for nested schema](#nestedblock--sync_settings--sync_logs_config))

<a id="nestedblock--sync_settings--failed_keys_config"></a>
### Nested Schema for `sync_settings.failed_keys_config`

Optional:

- `enable_failed_keys_retry` (Boolean)


<a id="nestedblock--sync_settings--sync_logs_config"></a>
### Nested Schema for `sync_settings.sync_logs_config`

Optional:

- `enabled` (Boolean)
- `log_retention_in_days` (Number)
- `snapshots_to_retain` (Number)
