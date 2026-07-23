package configs

import "fmt"

var (
	Sources      *Registry = &Registry{name: "sources"}
	Destinations *Registry = &Registry{name: "destinations", versioned: true}
	Accounts     *Registry = &Registry{name: "accounts"}
)

type apiTypeVersion struct {
	apiType string
	version int
}

type Registry struct {
	name    string
	entries map[string]ConfigMeta

	// versioned enables the (APIType, Version) uniqueness guard and the
	// Version-must-be-set guard below. Only Destinations is versioned today;
	// Sources and Accounts don't carry a meaningful Version yet.
	versioned bool
	byVersion map[apiTypeVersion]string // (apiType, version) -> name, only populated when versioned
}

func (r *Registry) Register(name string, cm ConfigMeta) {
	if r.entries == nil {
		r.entries = map[string]ConfigMeta{}
	}

	if _, ok := r.entries[name]; ok {
		panic(fmt.Errorf("name '%s' is already registered with %s", name, r.name))
	}

	if r.versioned {
		if cm.Version == 0 {
			panic(fmt.Errorf("'%s' must set ConfigMeta.Version when registering with %s", name, r.name))
		}

		if r.byVersion == nil {
			r.byVersion = map[apiTypeVersion]string{}
		}

		key := apiTypeVersion{apiType: cm.APIType, version: cm.Version}
		if existing, ok := r.byVersion[key]; ok {
			panic(fmt.Errorf("apiType '%s' version %d is already registered as '%s' in %s", cm.APIType, cm.Version, existing, r.name))
		}
		r.byVersion[key] = name
	}

	r.entries[name] = cm
}

func (r *Registry) Entries() map[string]ConfigMeta {
	return r.entries
}
