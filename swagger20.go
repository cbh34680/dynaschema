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

// FindParameters ... func
func (me *SchemaSwagger20) FindParameters(argPath, argMethod, argIn string) (string, error) {

	required := dynajson.NewAsArray()
	properties := dynajson.NewAsMap()

	err := me.eachParams(argPath, argMethod, argIn, func(pos int, spec *dynajson.JSONElement) (bool, error) {

		spName := spec.Select("name").AsString()
		if spName == "" {
			return false, fmt.Errorf("name is empty")
		}

		if spec.Select("required").AsBool() {
			required.Append(spName)
		}

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

		properties.Put(spName, property)

		return true, nil
	})

	if err != nil {
		return "", fmt.Errorf("me.eachParams: %w", err)
	}

	if properties.Count() == 0 {
		return "", nil
	}

	schema := dynajson.NewAsMap()
	schema.Put("type", "object")

	err = schema.Put("required", required)
	if err != nil {
		return "", fmt.Errorf("schema.Put(required): %w", err)
	}

	err = schema.Put("properties", properties)
	if err != nil {
		return "", fmt.Errorf("schema.Put(properties): %w", err)
	}

	return schema.String(), nil
}

// FindBody ... func
func (me *SchemaSwagger20) FindBody(argPath, argMethod, _ string) (string, error) {

	var schema *dynajson.JSONElement

	err := me.eachParams(argPath, argMethod, "body", func(pos int, spec *dynajson.JSONElement) (bool, error) {

		//
		// TODO: required
		//

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
		return "", nil
	}

	return schema.String(), nil
}
