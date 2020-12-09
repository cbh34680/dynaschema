package dynaschema

import (
	"fmt"

	"github.com/cbh34680/dynajson"
)

// SchemaOpenAPI30 ... func
type SchemaOpenAPI30 struct {
	SchemaAbstract
}

// NewSchemaOpenAPI30 ... func
func NewSchemaOpenAPI30(argJSON *dynajson.JSONElement) JSONSchema {
	ret := SchemaOpenAPI30{}
	ret.SchemaAbstract.objJSON = argJSON
	return &ret
}

// ---------------------------------------------------------------------------

// GetServers ... func
func (me *SchemaOpenAPI30) GetServers() ([]string, error) {

	root := me.objJSON

	servers := root.Select("servers")
	if servers.IsNil() {
		return nil, fmt.Errorf("servers undefined")
	}

	serversLen := servers.Count()

	if serversLen == 0 {
		return nil, fmt.Errorf("empty servers")
	}

	ret := make([]string, serversLen)

	err := servers.EachArray(func(pos int, elem *dynajson.JSONElement) (bool, error) {

		url := elem.Select("url")
		if url.IsNil() {
			return false, fmt.Errorf("servers(%d): url undefined", pos)
		}

		ret[pos] = url.AsString()

		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("servers.EachArray: %w", err)
	}

	return ret, nil
}

// FindParameter ... func
func (me *SchemaOpenAPI30) FindParameter(argPath, argMethod, argIn string) (string, error) {

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
func (me *SchemaOpenAPI30) FindBody(argPath, argMethod, argContent string) (string, error) {

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
