# Source or Destination Config
Configuration of Rudderstack sources and destinations are JSON objects. Any such JSON object can be mapped into [terraform schema](#nestedatt-config) of Rudderstack's Terraform provider as shown in the examples.

## Example usage
Config attribute for a supported resource is set as follows.
```
  .
  .
  config = {
    object = {
      "trackingID": { str = "UA-908213012-193" },
      "doubleClick": { bool = true },
      "enhancedLinkAttribution": { bool = true },
      "includeSearch": { bool = true },
      "enableServerSideIdentify": { bool = true },
      "serverSideIdentifyEventCategory": { str = "myEventCategory1" },
      "serverSideIdentifyEventAction": { str = "myEventAction1" },
      "disableMd5": { bool = true },
      "anonymizeIp": { bool = true },
      "enhancedEcommerce": { bool = true },
      "nonInteraction": { bool = true },
      "sendUserId": { bool = true },
      "dimensions": {
        objects_list = [
          {
             object = {
               "from": { str = "myDimensionSource" },
               "to": { str = "3" },
             }
          }
        ]
      },
      "metrics": {
        objects_list = [
          {
             object = {
               "from": { str = "myMetricSource" },
               "to": { str = "2" },
             }
          }
        ]
      },
      "contentGroupings": {
        objects_list = [
          {
             object = {
               "from": { str = "myContentGroupingSource" },
               "to": { str = "myContentGroupingDest" },
             }
          }
        ]
      },
    },
  }
  .
  .
```

Check out a complete destination config example [here](destination.md#example). 

## Config Schema

<a id="nestedatt--config"></a>
### Nested Schema for `config`
Config of any source or destination is represented as a JSON object. Any such JSON object can be converted into equivalent
rudderstack provider config as in examples below.

#### Examples
|JSON Config Object  | Equivalent Terraform Config                             |
|--------------------|---------------------------------------------------------|
|{}                  | { object = {} }                                         |
|{"a":1,"b":"strval"}| { object = {"a":{int = 1},"b":{str="strval"}} }         |
|{"a":[]}            | { object = {"a":{object_list=[]}} }                     |
|{"a":{}}            | { object = {"a":{object={}}} }                     |

#### Attributes

- **object** (Attributes Map, Required): object is a required attribute, mapping string keys to values. Values in the attribute map must comply with [config value schema](#nestedatt--config--value).

<a id="nestedatt--config--str"></a>
### Nested Schema for `config value`
Each JSON value can be converted into equivalent rudderstack provider config as follows.

#### Examples
|JSON Config Value   | Remark                   | Equivalent Representation in Terraform                  |
|--------------------|--------------------------|---------------------------------------------------------|
|123                 | Integer                  | { int = 123 }                                           |
|"arbit string"      | String                   | { str = "arbit string" }                                |
|true                | Boolean                  | { bool = true }                                         |
|{"a":1,"b":"strval"}| JSON Object              | { object = {"a":{int = 1},"b":{str="strval"}} }         |
|[{}, {"a":1}]       | JSON List of Objects     | { object_list = [<BR/>.  { object = {} },<BR/>.  { object = {"a":{int = 1},"b":{str="strval"}} }<BR/>] }|

#### Attributes:

Depending on the kind of JSON value, *EXACTLY ONE* of the following attributes must be set.

- **bool** (Boolean) Set this attribute if the JSON value is a boolean.
- **num** (Number) Set this attribute if the JSON value is an integer or float.
- **str** (String) Set this attribute if the JSON value is a string.
- **object** (Attributes Map) Set this attribute if the JSON value is an object. Define it as a string map, each value in attribute map following [this schema](#nestedatt--config--object)
- **objects_list** (Attributes List) Set this attribute if the JSON value is a list of objects. Define it as a list, each value in the attribute list following [this schema](#nestedatt--config--objects)

