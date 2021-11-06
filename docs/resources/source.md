# rudderstack_source (Resource)
Manages a Rudderstack CDP source.

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

## Argument Reference 

### Required

- **name** (String, Required) Specifies name of the resource.
- **type** (String, Required) Selects category of the CDP source to be created. Examples include iOS, Android, Auth0, etc.
  Consult Rudderstack documentation for list of supported source types.  
- **config** (Attributes) (Check [this](config.md) for schema)

### Read-Only

- **id** (String) The ID of this resource.


