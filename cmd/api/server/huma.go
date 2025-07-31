package server

import (
	"encoding/json"
	"reflect"
	"strings"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/common"

	"github.com/danielgtaylor/huma/v2"
	"gopkg.in/yaml.v3"
)

func newHumaConfig(title string, version string) huma.Config {
	huma.NewError = apierrors.NewApiError()
	huma.DefaultArrayNullable = false

	humaConfig := huma.DefaultConfig(title, version)
	humaConfig.Components.Schemas = NewCustomRegistry()

	return humaConfig
}

// CustomRegistry overrides the default Huma registry behaviour by inlining unique requests/responses  and only creating reusable schemas for shared models.
// This is needed because the default Huma registry creates a schema for every request and response object which fills the docs with a lot of noise.
type CustomRegistry struct {
	huma.Registry
}

func NewCustomRegistry() *CustomRegistry {
	return &CustomRegistry{
		huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer),
	}
}

func (r *CustomRegistry) Schema(t reflect.Type, allowRef bool, hint string) *huma.Schema {
	v := reflect.New(t).Interface()

	if _, ok := v.(common.InlinedSchema); ok ||
		// Types defined in the `handlers` package are tied to a specific handler and are inlined into its endpoint's definition.
		strings.Contains(t.PkgPath(), "/handlers/") {
		return huma.SchemaFromType(r, t)
	}

	return r.Registry.Schema(t, allowRef, hint)
}

func (r *CustomRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Registry) //nolint:wrapcheck
}

func (r *CustomRegistry) MarshalYAML() (any, error) {
	return yaml.Marshal(r.Registry) //nolint:wrapcheck
}
