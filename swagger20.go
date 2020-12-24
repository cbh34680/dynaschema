package dynaschema

import (
	"fmt"

	"github.com/cbh34680/dynajson"
)

// SchemaSwagger20 ... func
type SchemaSwagger20 struct {
	SchemaAbstract
}

// NewSchemaSwagger20 ... func
func NewSchemaSwagger20(argJSON *dynajson.JSONElement) JSONSchema {
	ret := SchemaSwagger20{}
	ret.SchemaAbstract.objJSON = argJSON
	return &ret
}

// ---------------------------------------------------------------------------

// GetServers ... func
func (me *SchemaSwagger20) GetServers() ([]string, error) {

	root := me.objJSON

	host := root.Select("host")
	if host.IsNil() {
		return nil, fmt.Errorf("host undefined")
	}

	basePath := root.Select("basePath").AsString()

	schemes := root.Select("schemes")
	schemesLen := schemes.Count()

	if schemesLen == 0 {
		return nil, fmt.Errorf("empty schemes")
	}

	ret := make([]string, schemesLen)

	err := schemes.EachArray(func(pos int, scheme *dynajson.JSONElement) (bool, error) {

		ret[pos] = fmt.Sprintf("%s://%s%s", scheme.AsString(), host.AsString(), basePath)

		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("schemes.EachArray: %w", err)
	}

	return ret, nil
}

// FindParameter ... func
func (me *SchemaSwagger20) FindParameter(argPath, argMethod, argIn string) (string, error) {

	return me.SchemaAbstract.findParameterHelper(argPath, argMethod, argIn, func(spec *dynajson.JSONElement) map[string]interface{} {

		property := map[string]interface{}{}

		spec.EachMap(func(key string, val *dynajson.JSONElement) (bool, error) {

			switch key {
			case "name", "in", "required":
				break
			default:
				property[key] = val.Raw()
			}

			return true, nil
		})

		return property
	})
}

// FindBody ... func
func (me *SchemaSwagger20) FindBody(argPath, argMethod, _ string) (string, error) {

	var schema *dynajson.JSONElement

	err := me.eachParams(argPath, argMethod, "body", func(pos int, spec *dynajson.JSONElement) (bool, error) {

		chkSchema := spec.Select("schema")

		if chkSchema.IsNil() {
			return true, nil
		}

		schema = chkSchema

		return false, nil
	})

	if err != nil {
		return "", fmt.Errorf("me.eachParams: %w", err)
	}

	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}

	return schema.String(), nil
}
