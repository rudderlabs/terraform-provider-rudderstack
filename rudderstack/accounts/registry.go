package accounts

import (
	"fmt"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

type AccountCategory string

const (
	CategorySource      AccountCategory = "source"
	CategoryDestination AccountCategory = "destination"
)

type AccountConfigMeta struct {
	configs.ConfigMeta
	Category AccountCategory
	Secret   []configs.ConfigProperty
}

func (cm *AccountConfigMeta) SecretStateToAPI(state string) (string, error) {
	api := "{}"

	for _, p := range cm.Secret {
		r, err := p.FromStateFunc(api, state)
		if err != nil {
			return api, err
		}
		api = r
	}

	return api, nil
}

func (cm *AccountConfigMeta) SecretAPIToState(api string) (string, error) {
	state := "{}"
	for _, p := range cm.Secret {
		s, err := p.ToStateFunc(state, api)
		if err != nil {
			return state, err
		}
		state = s
	}

	return state, nil
}

type registry struct {
	entries map[AccountCategory]map[string]AccountConfigMeta
}

func (r *registry) Register(name string, cm AccountConfigMeta) error {
	if r.entries == nil {
		r.entries = map[AccountCategory]map[string]AccountConfigMeta{}
	}

	categoryEntries := r.entries[cm.Category]
	if categoryEntries == nil {
		r.entries[cm.Category] = map[string]AccountConfigMeta{}
		categoryEntries = r.entries[cm.Category]
	}

	if _, ok := categoryEntries[name]; ok {
		return fmt.Errorf("name '%s' is already registered with %s accounts repository", cm.Category, name)
	}

	categoryEntries[name] = cm

	return nil
}

func (r *registry) Entries() map[AccountCategory]map[string]AccountConfigMeta {
	return r.entries
}

var Accounts = &registry{}
