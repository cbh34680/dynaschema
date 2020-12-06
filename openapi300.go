package dynaschema

import (
	"fmt"

	"github.com/cbh34680/dynajson"
	"github.com/xeipuuv/gojsonschema"
)

// SchemaOpenAPI300 ... func
type SchemaOpenAPI300 struct {
	objJSON *dynajson.JSONElement
}

// NewSchemaOpenAPI300 ... func
func NewSchemaOpenAPI300(argJSON *dynajson.JSONElement) JSONSchema {
	return &SchemaOpenAPI300{
		objJSON: argJSON,
	}
}

// RawJSON ... func
func (me *SchemaOpenAPI300) RawJSON() *dynajson.JSONElement {
	return me.objJSON
}

// String ... func
func (me *SchemaOpenAPI300) String() string {
	return me.objJSON.String()
}

// ValidateJSONRequestBody ... func
func (me *SchemaOpenAPI300) ValidateJSONRequestBody(argPath, argMethod, argJSON string) (*gojsonschema.Result, error) {

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

	dataLoader := gojsonschema.NewStringLoader(argJSON)
	schemaLoader := gojsonschema.NewStringLoader(schema.String())

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return nil, fmt.Errorf("gojsonschema.Validate: %w", err)
	}

	return result, nil
}
