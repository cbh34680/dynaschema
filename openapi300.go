package dynaschema

import (
	"fmt"

	"github.com/cbh34680/dynajson"
	"github.com/xeipuuv/gojsonschema"
)

// SchemaOpenAPI300 ... func
type SchemaOpenAPI300 struct {
	SchemaAbstract
}

// NewSchemaOpenAPI300 ... func
func NewSchemaOpenAPI300(argJSON *dynajson.JSONElement) JSONSchema {
	ret := SchemaOpenAPI300{}
	ret.SchemaAbstract.objJSON = argJSON
	return &ret
}

// String ... func
func (me *SchemaOpenAPI300) String() string {
	return me.objJSON.String()
}

// ValidateParameters ... func
func (me *SchemaOpenAPI300) ValidateParameters(argPath, argMethod, argIn string, argData map[string]interface{}) (*gojsonschema.Result, error) {

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
			case "schema":
				val.EachMap(func(sKey string, sVal *dynajson.JSONElement) (bool, error) {

					property[sKey] = sVal.Raw()
					return true, nil
				})
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
func (me *SchemaOpenAPI300) ValidateJSONRequestBody(argPath, argMethod, argData string) (*gojsonschema.Result, error) {

	root := me.objJSON

	requestBody := root.Select("paths", argPath, argMethod, "requestBody")
	if requestBody.IsNil() {
		//return nil, fmt.Errorf("requestBody: Select return nil")
		return &gojsonschema.Result{}, nil
	}

	if !requestBody.IsMap() {
		return nil, fmt.Errorf("not requestBody.IsMap()")
	}

	//
	// TODO: requestBody.required
	//

	schema := requestBody.Select("content", "application/json", "schema")
	if schema.IsNil() {
		//return nil, fmt.Errorf("schema: Select return nil")
		return &gojsonschema.Result{}, nil
	}

	dataLoader := gojsonschema.NewStringLoader(argData)
	schemaLoader := gojsonschema.NewStringLoader(schema.String())

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return nil, fmt.Errorf("gojsonschema.Validate: %w", err)
	}

	return result, nil
}
