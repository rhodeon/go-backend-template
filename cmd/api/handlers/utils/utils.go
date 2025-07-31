package handlerutils

import (
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/rhodeon/go-backend-template/cmd/api/models/common"
)

// GenerateSchemaName builds the schema names of request and response data.
// The naming format is: /path-1/path-2/data -> Path1Path2Data
// E.g. /users/create-with-list -> UsersCreateWithListRequestBody & UsersCreateWithListResponseData
func GenerateSchemaName(schema common.OasSchema) string {
	t := reflect.TypeOf(schema)
	name := t.Name()
	pkg := t.PkgPath()
	pkgs := strings.Split(pkg, "/")

	prefix := ""
	// The iteration is done backwards to stop at the "handlers" package as that's the point of reference for all handlers.
	for i := len(pkgs) - 1; i >= 0; i-- {
		if pkgs[i] == "handlers" {
			break
		}
		prefix = titleCase(pkgs[i]) + prefix
	}

	return prefix + name
}

func titleCase(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(r)) + s[size:]
}
