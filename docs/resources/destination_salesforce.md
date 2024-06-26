---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rudderstack_destination_salesforce Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-
  
---

# rudderstack_destination_salesforce (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config` (Block List, Min: 1, Max: 1) Destination specific configuration. Check the nested block documenation for more information. (see [below for nested schema](#nestedblock--config))
- `name` (String) Human readable name of the destination. The value has to be unique across all the destinations.

### Optional

- `enabled` (Boolean) An enabled destination allows data to be sent to it.

### Read-Only

- `created_at` (String) Time when the resource was created, in ISO 8601 format.
- `id` (String) The ID of this resource.
- `updated_at` (String) Time when the resource was last updated, in ISO 8601 format.

<a id="nestedblock--config"></a>
### Nested Schema for `config`

Required:

- `initial_access_token` (String, Sensitive) Enter your Salesforce security token.
- `password` (String, Sensitive) Enter the password for the above user.
- `user_name` (String) Enter the Salesforce username.

Optional:

- `map_properties` (Boolean) Use this setting to map RudderStack event properties to specific Salesforce fields.
- `onetrust_cookie_categories` (List of String) Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.
- `sandbox` (Boolean) Use this setting to enable Salesforce sandbox mode.
- `use_contact_id` (Boolean) When enabled, RudderStack uses contactId for the converted leads.


