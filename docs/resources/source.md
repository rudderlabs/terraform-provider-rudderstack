# Resource rudderstack_source
Manages a RudderStack CDP source.

## Attribute Reference 

### Required

- **name** (String, Required) Specifies name of the resource.
- **type** (String, Required) Selects category of the CDP source to be created. Examples include iOS, Android, Auth0, etc.
  Consult RudderStack documentation for list of supported source types.  
- **config** (Attributes, Required) Check [this](../guides/config.md) for schema and <mark>how to [create](../guides/config.md#creating-config) config</mark>.

### Read-Only

- **id** (String) The ID of this resource.

## Example Usage
```
resource "rudderstack_source" "src1" {
  name = "tfsource"
  type = "Auth0"
  config = {
    object = {
    }
  }
}
```
