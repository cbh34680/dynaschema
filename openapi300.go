package dynaschema

import (
	"fmt"

	"github.com/cbh34680/dynajson"
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

// FindParameter ... func
func (me *SchemaOpenAPI300) FindParameter(argPath, argMethod, argIn string) (string, error) {

	return me.SchemaAbstract.findParameterHelper(argPath, argMethod, argIn, func(spec *dynajson.JSONElement) map[string]interface{} {

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

		return property
	})
	/*
		required := dynajson.NewAsArray()
		properties := dynajson.NewAsMap()

		err := me.eachParams(argPath, argMethod, argIn, func(pos int, spec *dynajson.JSONElement) (bool, error) {

			spName := spec.Select("name").AsString()
			if spName == "" {
				return false, fmt.Errorf("name is empty")
			}

			spRequired := spec.Select("required").AsBool()

			if spRequired {
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

			if spRequired {
				if me.GetFlag(Flag.SetMinlenIfRequired) {
					if propType, ok := property["type"]; ok {
						if propType == "string" {
							if _, ok := property["minLength"]; !ok {
								property["minLength"] = 1
							}
						}
					}
				}
			}

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
	*/
}

// FindBody ... func
func (me *SchemaOpenAPI300) FindBody(argPath, argMethod, argContent string) (string, error) {

	root := me.objJSON

	requestBody := root.Select("paths", argPath, argMethod, "requestBody")
	if requestBody.IsNil() {
		return "", nil
	}

	if !requestBody.IsMap() {
		return "", fmt.Errorf("not requestBody.IsMap()")
	}

	//
	// TODO: requestBody.required
	//

	schema := requestBody.Select("content", argContent, "schema")
	if schema.IsNil() {
		return "", nil
	}

	return schema.String(), nil
}
