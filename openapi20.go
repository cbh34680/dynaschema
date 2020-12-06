package dynaschema

import (
	"fmt"

	"github.com/cbh34680/dynajson"
	"github.com/xeipuuv/gojsonschema"
)

// SchemaOpenAPI20 ... func
type SchemaOpenAPI20 struct {
	objJSON *dynajson.JSONElement
}

// NewSchemaOpenAPI20 ... func
func NewSchemaOpenAPI20(argJSON *dynajson.JSONElement) JSONSchema {
	return &SchemaOpenAPI20{
		objJSON: argJSON,
	}
}

// RawJSON ... func
func (me *SchemaOpenAPI20) RawJSON() *dynajson.JSONElement {
	return me.objJSON
}

// String ... func
func (me *SchemaOpenAPI20) String() string {
	return me.objJSON.String()
}

// ValidateJSONRequestBody ... func
func (me *SchemaOpenAPI20) ValidateJSONRequestBody(argPath, argMethod, argJSON string) (*gojsonschema.Result, error) {

	root := me.objJSON

	parameters := root.Select("paths", argPath, argMethod, "parameters")
	if parameters.IsNil() {
		return nil, fmt.Errorf("parameters: Select return nil")
	}

	if !parameters.IsArray() {
		return nil, fmt.Errorf("not parameters.IsArray()")
	}

	dataLoader := gojsonschema.NewStringLoader(argJSON)

	var result *gojsonschema.Result = &gojsonschema.Result{}

	parameters.EachArray(func(pos int, spec *dynajson.JSONElement) (bool, error) {

		if spec.Select("in").AsString() != "body" {
			return true, nil
		}

		//
		// TODO: required
		//

		schema := spec.Select("schema")
		if schema.IsNil() {
			return true, nil
		}
		//fmt.Println(schema)

		schemaLoader := gojsonschema.NewStringLoader(schema.String())

		tmpResult, err := gojsonschema.Validate(schemaLoader, dataLoader)
		if err != nil {
			return false, fmt.Errorf("gojsonschema.Validate: %w", err)
		}

		result = tmpResult

		return true, nil
	})

	return result, nil
}
