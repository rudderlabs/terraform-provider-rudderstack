package configs

import "fmt"

var (
	Sources      *Registry = &Registry{name: "sources"}
	Destinations *Registry = &Registry{name: "destinations"}
)

type Registry struct {
	name    string
	entries map[string]ConfigMeta
}

func (r *Registry) Register(name string, cm ConfigMeta) {
	if r.entries == nil {
		r.entries = map[string]ConfigMeta{}
	}

	if _, ok := r.entries[name]; ok {
		panic(fmt.Errorf("name '%s' is already registered with %s", name, r.name))
	}

	r.entries[name] = cm
}

func (r *Registry) Entries() map[string]ConfigMeta {
	return r.entries
}
