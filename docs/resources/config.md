# Source or Destination Config (Common Attribute)
REST API for Rudderstack defines CRUD operations on Rudderstack sources and destinations. In those CRUD operatoins, configuration of sources and destinations is modelled as a JSON object. To define config of such a source or a destination in this provider, the JSON object needs to be converted into terraform schema [Config](#nestedatt-config).

## Example usage
Checkout destination config example [here](destination.md#example). 

## Schema

<a id="nestedatt--config"></a>
### Nested Schema for `config`
Any JSON object can be converted into equivalent rudderstack provider config. Examples as follows.

#### Examples
Any JSON object can be converted into equivalent terraform config.
|--------------------|---------------------------------------------------------|
|JSON Config Object  | Equivalent Terraform Config                             |
|--------------------|---------------------------------------------------------|
|{}                  | { object = {} }                                         |
|{"a":1,"b":"strval"}| { object = {"a":{int = 1},"b":{str="strval"}} }         |
|{"a":[]}            | { object = {"a":{object_list=[]}} }                     |
|{"a":{}}            | { object = {"a":{object={}}} }                     |
|--------------------|---------------------------------------------------------|

Required:

- **object** (Attributes Map, Required): object is a required attribute. Values in the attribute map must comply with [nested config value schema](#nestedatt--config--value).

<a id="nestedatt--config--str"></a>
### Nested Schema for `config value`
Each JSON value can be converted equivalent rudderstack provider config as follows.

#### Examples
|--------------------|--------------------------|---------------------------------------------------------|
|JSON Config Value   | Remark                   | Equivalent Representation in Terraform                  |
|--------------------|--------------------------|---------------------------------------------------------|
|123                 | Integer                  | { int = 123 }                                           |
|"arbit string"      | String                   | { str = "arbit string" }                                |
|true                | Boolean                  | { bool = true }                                         |
|{"a":1,"b":"strval"}| JSON Object              | { object = {"a":{int = 1},"b":{str="strval"}} }         |
|[{}, {"a":1}]       | JSON List of Objects     | { object_list = [                                       |
|                    |                          |        { object = {} },                                    |
|                    |                          |        { object = {"a":{int = 1},"b":{str="strval"}} }  |
|                    |                          | ] },                                                    |
|---------------------------------------------------------------------------------------------------------|

#### Attributes:

Depending on the kind of JSON object, *EXACTLY ONE* of the following attributes can be set.

- **bool** (Boolean) Set this attribute if the JSON value is a boolean.
- **num** (Number) Set this attribute if the JSON value is an integer or float.
- **str** (String) Set this attribute if the JSON value is a string.
- **object** (Attributes Map) Define this attribute as a map, if the JSON value is an object. Each value in attribute map follows [this nested schema](#nestedatt--config--object)
- **objects_list** (Attributes List) Define this attribute as a list, if the JSON value is a list of objects. Each value in attribute list follows [this nested schema](#nestedatt--config--objects)

