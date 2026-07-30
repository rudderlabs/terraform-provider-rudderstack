package rudderstack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations/destinations"
)

func resourceDestination(cm configs.ConfigMeta) *schema.Resource {
	return &schema.Resource{
		Schema:        resourceDestinationSchema(cm),
		CreateContext: resourceDestinationCreate(cm),
		ReadContext:   resourceDestinationRead(cm),
		UpdateContext: resourceDestinationUpdate(cm),
		DeleteContext: resourceDestinationDelete(cm),
		CustomizeDiff: resourceDestinationCustomizeDiff(cm),
		Importer: &schema.ResourceImporter{
			StateContext: resourceDestinationImportState(cm),
		},
	}
}

func resourceDestinationCustomizeDiff(cm configs.ConfigMeta) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
		if cm.SkipConfig {
			return nil
		}

		consentManagement := d.Get("config.0.consent_management.0")
		if consentManagement == nil {
			return nil
		}

		return destinations.ValidateConsentManagementUniqueProviders(consentManagement)
	}
}

func resourceDestinationSchema(cm configs.ConfigMeta) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Human readable name of the destination. The value has to be unique across all the destinations.",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "An enabled destination allows data to be sent to it.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time when the resource was created, in ISO 8601 format.",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time when the resource was last updated, in ISO 8601 format.",
		},
	}

	if !cm.SkipConfig {
		s["config"] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: cm.SkipConfig,
			Required: !cm.SkipConfig,
			Description: "Destination specific configuration. Check the nested block documenation " +
				"for more information.",
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: cm.ConfigSchema,
			},
		}
	}

	return s
}

func resourceDestinationCreate(cm configs.ConfigMeta) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		destination := &client.Destination{}
		err := populateDestinationFromState(cm, destination, d)
		if err != nil {
			return diag.FromErr(err)
		}

		destination, err = c.Destinations.Create(ctx, destination)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create destination: %w", err))
		}

		d.SetId(destination.ID)

		return resourceDestinationRead(cm)(ctx, d, m)
	}
}

func resourceDestinationRead(cm configs.ConfigMeta) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		id := d.Id()

		destination, err := c.Destinations.Get(ctx, id)
		if err != nil {
			var apiErr *client.APIError
			if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
			return diag.FromErr(err)
		}

		if err := storeDestinationToState(cm, destination, d); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{}
	}
}

func resourceDestinationUpdate(cm configs.ConfigMeta) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		destination := &client.Destination{}
		err := populateDestinationFromState(cm, destination, d)
		if err != nil {
			return diag.FromErr(err)
		}

		destination, err = c.Destinations.Update(ctx, destination)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not update destination: %w", err))
		}

		d.SetId(destination.ID)

		return resourceDestinationRead(cm)(ctx, d, m)
	}
}

func resourceDestinationDelete(cm configs.ConfigMeta) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		if err := c.Destinations.Delete(ctx, d.Id()); err != nil {
			var apiErr *client.APIError
			if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
			return diag.FromErr(err)
		}

		d.SetId("")
		return diag.Diagnostics{}
	}
}

func resourceDestinationImportState(cm configs.ConfigMeta) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		diagnostics := resourceDestinationRead(cm)(ctx, d, m)
		if diagnostics.HasError() {
			for _, diagnostic := range diagnostics {
				if diagnostic.Severity == diag.Error {
					return nil, fmt.Errorf("could not import connection: %s", diagnostic.Summary)
				}
			}
		}
		return []*schema.ResourceData{d}, nil
	}
}

func populateDestinationFromState(cm configs.ConfigMeta, destination *client.Destination, d *schema.ResourceData) error {
	destination.ID = d.Id()
	destination.Type = cm.APIType
	destination.Version = cm.Version
	destination.Name = d.Get("name").(string)
	destination.IsEnabled = d.Get("enabled").(bool)

	if c := d.Get("config.0"); c != nil {
		state, err := json.Marshal(c)
		if err != nil {
			return err
		}
		apiConfig, err := cm.StateToAPI(string(state))
		if err != nil {
			return err
		}
		destination.Config = json.RawMessage(apiConfig)
	}

	return nil
}

func storeDestinationToState(cm configs.ConfigMeta, destination *client.Destination, d *schema.ResourceData) error {
	d.SetId(destination.ID)
	if err := d.Set("name", destination.Name); err != nil {
		return err
	}
	if err := d.Set("enabled", destination.IsEnabled); err != nil {
		return err
	}
	if destination.CreatedAt != nil {
		createdAt := destination.CreatedAt.Format(time.RFC3339)
		if err := d.Set("created_at", createdAt); err != nil {
			return err
		}
	}
	if destination.UpdatedAt != nil {
		updatedAt := destination.UpdatedAt.Format(time.RFC3339)
		if err := d.Set("updated_at", updatedAt); err != nil {
			return err
		}
	}

	state, err := cm.APIToState(string(destination.Config))
	if err != nil {
		return err
	}

	properties := make(map[string]interface{})
	if err := json.Unmarshal([]byte(state), &properties); err != nil {
		return err
	}
	properties = mergeDestinationConfigWithPriorState(properties, cm, d)

	if len(properties) > 0 {
		if err := d.Set("config", []interface{}{properties}); err != nil {
			return err
		}
	} else {
		if err := d.Set("config", []interface{}{}); err != nil {
			return err
		}
	}

	return nil
}

func mergeDestinationConfigWithPriorState(properties map[string]interface{}, cm configs.ConfigMeta, d *schema.ResourceData) map[string]interface{} {
	merged := make(map[string]interface{}, len(properties))
	for key, value := range properties {
		merged[key] = value
	}

	prior := priorDestinationConfig(d)
	if prior == nil {
		return merged
	}

	mergeDestinationConfigMap(merged, prior, cm.ConfigSchema)
	return merged
}

func priorDestinationConfig(d *schema.ResourceData) map[string]interface{} {
	existing := d.Get("config")
	list, ok := existing.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}

	prior, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return prior
}

func mergeDestinationConfigMap(merged, prior map[string]interface{}, schemaMap map[string]*schema.Schema) {
	for key, fieldSchema := range schemaMap {
		priorValue, ok := prior[key]
		if !ok {
			continue
		}

		apiValue, hasAPIValue := merged[key]
		if value, preserve := preservedDestinationConfigValue(apiValue, priorValue, fieldSchema, hasAPIValue); preserve {
			merged[key] = value
		}
	}
}

func preservedDestinationConfigValue(apiValue, priorValue interface{}, fieldSchema *schema.Schema, hasAPIValue bool) (interface{}, bool) {
	if fieldSchema.Sensitive {
		return priorValue, true
	}

	nestedSchema, hasNestedSchema := nestedDestinationConfigSchema(fieldSchema)
	if hasNestedSchema {
		if hasAPIValue {
			return mergeNestedDestinationConfigValue(apiValue, priorValue, nestedSchema)
		}
		return preservedNestedDestinationConfigValue(priorValue, nestedSchema)
	}

	if !hasAPIValue && isEmptyDestinationConfigCollection(priorValue) {
		return priorValue, true
	}

	return nil, false
}

func nestedDestinationConfigSchema(fieldSchema *schema.Schema) (map[string]*schema.Schema, bool) {
	resource, ok := fieldSchema.Elem.(*schema.Resource)
	if !ok {
		return nil, false
	}

	return resource.Schema, true
}

func mergeNestedDestinationConfigValue(apiValue, priorValue interface{}, schemaMap map[string]*schema.Schema) (interface{}, bool) {
	switch prior := priorValue.(type) {
	case []interface{}:
		api, ok := apiValue.([]interface{})
		if !ok {
			return nil, false
		}

		merged := make([]interface{}, len(api))
		copy(merged, api)
		preserved := false
		for index := range api {
			if index >= len(prior) {
				continue
			}

			apiItem, apiOK := api[index].(map[string]interface{})
			priorItem, priorOK := prior[index].(map[string]interface{})
			if !apiOK || !priorOK {
				continue
			}

			mergedItem := make(map[string]interface{}, len(apiItem))
			for key, value := range apiItem {
				mergedItem[key] = value
			}
			mergeDestinationConfigMap(mergedItem, priorItem, schemaMap)
			if !reflect.DeepEqual(mergedItem, apiItem) {
				preserved = true
			}
			merged[index] = mergedItem
		}

		return merged, preserved
	case map[string]interface{}:
		api, ok := apiValue.(map[string]interface{})
		if !ok {
			return nil, false
		}

		merged := make(map[string]interface{}, len(api))
		for key, value := range api {
			merged[key] = value
		}
		mergeDestinationConfigMap(merged, prior, schemaMap)
		return merged, !reflect.DeepEqual(merged, api)
	default:
		return nil, false
	}
}

func preservedNestedDestinationConfigValue(priorValue interface{}, schemaMap map[string]*schema.Schema) (interface{}, bool) {
	switch prior := priorValue.(type) {
	case []interface{}:
		merged := make([]interface{}, 0, len(prior))
		for _, item := range prior {
			priorItem, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			mergedItem := map[string]interface{}{}
			mergeDestinationConfigMap(mergedItem, priorItem, schemaMap)
			if len(mergedItem) > 0 {
				merged = append(merged, mergedItem)
			}
		}

		if len(merged) > 0 {
			return merged, true
		}
	case map[string]interface{}:
		merged := map[string]interface{}{}
		mergeDestinationConfigMap(merged, prior, schemaMap)
		if len(merged) > 0 {
			return merged, true
		}
	}

	return nil, false
}

func isEmptyDestinationConfigCollection(value interface{}) bool {
	switch v := value.(type) {
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	case *schema.Set:
		return v.Len() == 0
	default:
		return false
	}
}
