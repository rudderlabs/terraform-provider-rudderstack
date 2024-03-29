---
page_title: "rudderstack_connection Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_connection (Resource)

This resource represents a connection between a Rudderstack Source and Destination.

## Example Usage

```terraform
resource "rudderstack_connection" "example" {
  source_id      = rudderstack_source_javascript.example.id
  destination_id = rudderstack_destination_redshift.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination_id` (String) The ID of the connected destination.
- `source_id` (String) The ID of the connected source.

### Optional

- `enabled` (Boolean) An enabled connection allows data to be transferred from the connected source to the connected destination.

### Read-Only

- `created_at` (String) Time when the resource was created, in ISO 8601 format.
- `id` (String) The ID of this resource.
- `updated_at` (String) Time when the resource was last updated, in ISO 8601 format.
