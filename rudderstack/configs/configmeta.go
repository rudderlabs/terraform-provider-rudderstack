package configs

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type ConfigProperty struct {
	ConfigurationKey string
	APIKey           string
}

func prop(configurationKey, apiKey string) ConfigProperty {
	return ConfigProperty{
		ConfigurationKey: configurationKey,
		APIKey:           apiKey,
	}
}

type ConfigMeta struct {
	APIType        string
	OptionalConfig bool
	ConfigSchema   map[string]*schema.Schema
	Properties     []ConfigProperty
	TestConfigs    []TestConfig
}

type TestConfig struct {
	TerraformCreate string
	APICreate       string
	TerraformUpdate string
	APIUpdate       string
}

var emptyTestConfig = TestConfig{TerraformCreate: "", APICreate: "{}", TerraformUpdate: "", APIUpdate: "{}"}

func (cm *ConfigMeta) ParseResourceData(d *schema.ResourceData) (json.RawMessage, bool) {
	result := "{}"

	for _, p := range cm.Properties {
		v, ok := d.GetOk(fmt.Sprintf("config.0.%s", p.ConfigurationKey))
		if ok && v != nil {
			result, _ = sjson.Set(result, p.APIKey, v)
		}
	}

	return json.RawMessage(result), true
}

func (cm *ConfigMeta) StoreResourceData(config json.RawMessage, d *schema.ResourceData) error {
	json := string(config)
	properties := make(map[string]interface{})
	for _, p := range cm.Properties {
		r := gjson.Get(json, p.APIKey)
		if r.Exists() {
			properties[p.ConfigurationKey] = r.Value()
		}
	}

	if len(properties) > 0 {
		d.Set("config", []interface{}{properties})
	} else {
		d.Set("config", nil)
	}

	return nil
}
