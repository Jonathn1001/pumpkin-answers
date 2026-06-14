// Package configschema exposes, per dimension, a JSON Schema generated from the Go
// struct (invopop/jsonschema — struct is the source of truth) plus UI widget
// metadata from the dimension. Computed once and cached.
package configschema

import (
	"sync"

	_ "claimsplatform/internal/dimensions" // register all six dimensions
	"claimsplatform/internal/registry"
	"github.com/invopop/jsonschema"
)

type DimensionSchema struct {
	Key        string                     `json:"key"`
	JSONSchema *jsonschema.Schema         `json:"jsonSchema"`
	UI         []registry.FieldDescriptor `json:"ui"`
}

type Response struct {
	Dimensions []DimensionSchema `json:"dimensions"`
}

var (
	once   sync.Once
	cached Response
)

// Get returns the cached schema set, building it on first call.
func Get() Response {
	once.Do(build)
	return cached
}

func build() {
	r := &jsonschema.Reflector{DoNotReference: true}
	dims := make([]DimensionSchema, 0, len(registry.All()))
	for _, d := range registry.All() {
		dims = append(dims, DimensionSchema{
			Key:        d.Key(),
			JSONSchema: r.Reflect(d.DefaultConfig()),
			UI:         d.UISchema(),
		})
	}
	cached = Response{Dimensions: dims}
}
