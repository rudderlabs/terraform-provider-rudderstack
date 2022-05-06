---
page_title: "rudderstack_destination_redshift Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_redshift (Resource)

This resource represents an Amazon Redshift warehouse destination. For more information check 
https://www.rudderstack.com/docs/reverse-etl/amazon-redshift

## Example Usage

```terraform
resource "rudderstack_destination_redshift" "example" {
  name = "my-redshift"

  config {
    host     = "localhost"
    port     = "5432"
    database = "example"
    user     = "postgres"
    password = "postgres"

    namespace          = "example"
    enable_sse         = true
    use_rudder_storage = false


    # s3 {
    #   bucket_name   = ""
    #   access_key_id = ""
    #   access_key    = ""
    # }

    sync {
      frequency = "30"

      # start_at                  = "10:00"
      # exclude_window_start_time = "11:00"
      # exclude_window_end_time   = "12:00"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config` (Block List, Min: 1, Max: 1) Destination specific configuration. Check the nested block documenation for more information. (see [below for nested schema](#nestedblock--config))
- `name` (String) Human readable name of the destination. The value has to be unique across all destinations.

### Optional

- `enabled` (Boolean) An enabled destination allows data to be sent to it.

### Read-Only

- `created_at` (String) Time when the resource was created, in ISO 8601 format.
- `id` (String) The ID of this resource.
- `updated_at` (String) Time when the resource was last updated, in ISO 8601 format.

<a id="nestedblock--config"></a>
### Nested Schema for `config`

Required:

- `database` (String)
- `host` (String)
- `namespace` (String)
- `password` (String, Sensitive)
- `port` (String)
- `sync` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--config--sync))
- `use_rudder_storage` (Boolean)
- `user` (String)

Optional:

- `enable_sse` (Boolean)
- `s3` (Block List, Max: 1) (see [below for nested schema](#nestedblock--config--s3))

<a id="nestedblock--config--sync"></a>
### Nested Schema for `config.sync`

Required:

- `frequency` (String)

Optional:

- `exclude_window_end_time` (String)
- `exclude_window_start_time` (String)
- `start_at` (String)


<a id="nestedblock--config--s3"></a>
### Nested Schema for `config.s3`

Required:

- `bucket_name` (String)

Optional:

- `access_key` (String, Sensitive)
- `access_key_id` (String)