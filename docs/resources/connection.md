# rudderstack_connection (Resource)
Manages a connection between Rudderstack source and Rudderstack destination.

## Example Usage
```
resource "rudderstack_connection" "cnxn1" {
  source_id = "${rudderstack_source.src1.id}"
  destination_id = "${rudderstack_destination.dst1.id}"
}
```
## Argument Reference 

### Required

- **destination_id** (String, Required) Represents ID of the source object.
- **source_id** (String, Required) Represents ID of the destination object.

### Read-Only

- **id** (String) The ID of the connection object.


