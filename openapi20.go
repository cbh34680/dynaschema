package dynaschema

import (
	"fmt"

	"github.com/cbh34680/dynajson"
	"github.com/xeipuuv/gojsonschema"
)

// SchemaOpenAPI20 ... func
type SchemaOpenAPI20 struct {
	SchemaAbstract
}

// NewSchemaOpenAPI20 ... func
func NewSchemaOpenAPI20(argJSON *dynajson.JSONElement) JSONSchema {
	ret := SchemaOpenAPI20{}
	ret.SchemaAbstract.objJSON = argJSON
	return &ret
}

// ---------------------------------------------------------------------------

// ValidateParameters ... func
func (me *SchemaOpenAPI20) ValidateParameters(argPath, argMethod, argIn string, argData map[string]interface{}) (*gojsonschema.Result, error) {

	schema := dynajson.NewAsMap()
	schema.Put("type", "object")

	required, err := schema.PutEmptyArray("required")
	if err != nil {
		return nil, fmt.Errorf("schema.PutEmptyArray: %w", err)
	}

	properties, err := schema.PutEmptyMap("properties")
	if err != nil {
		return nil, fmt.Errorf("schema.PutEmptyMap: %w", err)
	}

	err = me.eachParams(argPath, argMethod, argIn, func(pos int, spec *dynajson.JSONElement) (bool, error) {

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
		return nil, fmt.Errorf("me.eachParams: %w", err)
	}

	data := dynajson.New(argData)

	strSchema := schema.String()
	strData := data.String()

	//fmt.Println(strSchema)
	//fmt.Println(strData)

	schemaLoader := gojsonschema.NewStringLoader(strSchema)
	dataLoader := gojsonschema.NewStringLoader(strData)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return nil, fmt.Errorf("gojsonschema.Validate: %w", err)
	}

	return result, nil
}

// ValidateJSONRequestBody ... func
func (me *SchemaOpenAPI20) ValidateJSONRequestBody(argPath, argMethod, argData string) (*gojsonschema.Result, error) {

	dataLoader := gojsonschema.NewStringLoader(argData)

	var lastResult *gojsonschema.Result = &gojsonschema.Result{}

	me.eachParams(argPath, argMethod, "body", func(pos int, spec *dynajson.JSONElement) (bool, error) {

		//
		// TODO: required
		//

		schema := spec.Select("schema")
		if spec.Select("schema").IsNil() {
			return true, nil
		}

		//fmt.Println(schema)
		//fmt.Println(argData)

		schemaLoader := gojsonschema.NewStringLoader(schema.String())

		result, err := gojsonschema.Validate(schemaLoader, dataLoader)
		if err != nil {
			return false, fmt.Errorf("gojsonschema.Validate: %w", err)
		}

		lastResult = result

		if !result.Valid() {
			return false, nil
		}

		return true, nil
	})

	return lastResult, nil
}
