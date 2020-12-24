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
}

// FindBody ... func
func (me *SchemaOpenAPI30) FindBody(argPath, argMethod, argContent string) (string, error) {

	root := me.objJSON

	requestBody := root.Select("paths", argPath, argMethod, "requestBody")
	if requestBody.IsNil() {
		return "", fmt.Errorf("%s, %s: not found", argPath, argMethod)
	}

	if !requestBody.IsMap() {
		return "", fmt.Errorf("not requestBody.IsMap()")
	}

	schema := requestBody.Select("content", argContent, "schema")
	if schema.IsNil() {
		return "", fmt.Errorf("%s: not found", argContent)
	}

	return schema.String(), nil
}
