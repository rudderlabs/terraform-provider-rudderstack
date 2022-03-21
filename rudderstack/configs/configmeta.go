package configs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConfigMeta struct {
	APIType        string
	OptionalConfig bool
	ConfigSchema   map[string]*schema.Schema
	Properties     []ConfigProperty
	TestConfigs    []TestConfig
}

func (cm *ConfigMeta) StateToAPI(state string) (string, error) {
	api := "{}"

	for _, p := range cm.Properties {
		r, err := p.FromStateFunc(api, state)
		if err != nil {
			return api, err
		}
		api = r
	}

	return api, nil
}

func (cm *ConfigMeta) APIToState(api string) (string, error) {
	state := "{}"
	for _, p := range cm.Properties {
		s, err := p.ToStateFunc(state, api)
		if err != nil {
			return state, err
		}
		state = s
	}

	return state, nil
}
