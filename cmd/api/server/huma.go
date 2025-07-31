package server

import (
	"reflect"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/common"

	"github.com/danielgtaylor/huma/v2"
)

func newHumaConfig(title string, version string) huma.Config {
	huma.NewError = apierrors.NewApiError()
	huma.DefaultArrayNullable = false

	humaConfig := huma.DefaultConfig(title, version)

	schemaPrefix := "#/components/schemas/"
	registry := huma.NewMapRegistry(schemaPrefix, CustomSchemaNamer)
	humaConfig.OpenAPI.Components.Schemas = registry

	// CreateHooks of the config is set to an empty slice to omit the $schema field in the body of responses.
	// This helps to obscure struct paths in the codebase from the API.
	humaConfig.CreateHooks = []func(huma.Config) huma.Config{}

	return humaConfig
}

// CustomSchemaNamer uses the custom schema names of OasSchema types and falls back to the default huma.DefaultSchemaNamer otherwise.
func CustomSchemaNamer(t reflect.Type, hint string) string {
	if t.Implements(reflect.TypeOf((*common.OasSchema)(nil)).Elem()) ||
		reflect.PointerTo(t).Implements(reflect.TypeOf((*common.OasSchema)(nil)).Elem()) {
		v := reflect.New(t).Elem()
		namer := v.Interface().(common.OasSchema)
		return namer.Name()
	}

	return huma.DefaultSchemaNamer(t, hint)
}
