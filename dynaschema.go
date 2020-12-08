package dynaschema

import (
	"fmt"
	"strings"

	"github.com/cbh34680/dynajson"
	"github.com/xeipuuv/gojsonschema"
)

// JSONSchema ... struct
type JSONSchema interface {
	RawJSON() *dynajson.JSONElement
	String() string
	ValidateParameters(string, string, string, map[string]interface{}) (*gojsonschema.Result, error)
	ValidateJSONRequestBody(string, string, string) (*gojsonschema.Result, error)
}

// SchemaAbstract ... struct
type SchemaAbstract struct {
	objJSON *dynajson.JSONElement
}

// RawJSON ... func
func (me *SchemaAbstract) RawJSON() *dynajson.JSONElement {
	return me.objJSON
}

// String ... func
func (me *SchemaAbstract) String() string {
	return me.objJSON.String()
}

// StrMap2AnyMap ... func
func StrMap2AnyMap(arg map[string]string) map[string]interface{} {

	if arg == nil {
		return nil
	}

	ret := make(map[string]interface{}, len(arg))
	for k, v := range arg {

		ret[k] = v
	}
	return ret
}

func (me *SchemaAbstract) eachParams(argPath, argMethod, argIn string, callback func(pos int, spec *dynajson.JSONElement) (bool, error)) error {

	root := me.objJSON

	parameters := root.Select("paths", argPath, argMethod, "parameters")
	if parameters.IsNil() {
		return fmt.Errorf("parameters: Select return nil")
	}

	if !parameters.IsArray() {
		return fmt.Errorf("not parameters.IsArray()")
	}

	return parameters.EachArray(func(pos int, spec *dynajson.JSONElement) (bool, error) {

		if argIn != "" {
			if spec.Select("in").AsString() != argIn {

				return true, nil
			}
		}

		return callback(pos, spec)
	})
}

// ---------------------------------------------------------------------------

// New ... func
func New(root *dynajson.JSONElement) (JSONSchema, error) {

	ver := root.Select("openapi").AsString()
	if ver == "" {
		ver = root.Select("swagger").AsString()
		if ver == "" {
			return nil, fmt.Errorf("UnSupported Schema")
		}
	}

	err := expandRef(root)
	if err != nil {
		return nil, fmt.Errorf("expandRef: %w", err)
	}

	//fmt.Println(root)

	root.Readonly = true

	switch ver {
	case "2.0":
		return NewSchemaOpenAPI20(root), nil
	case "3.0.0":
		return NewSchemaOpenAPI300(root), nil
	}

	return nil, fmt.Errorf("UnSupported Schema-Version")
}

// NewByBytes ... func
func NewByBytes(data []byte) (JSONSchema, error) {

	root, err := dynajson.NewByBytes(data)
	if err != nil {
		return nil, fmt.Errorf("dynajson.NewByBytes: %w", err)
	}

	return New(root)
}

// NewByString ... func
func NewByString(data string) (JSONSchema, error) {

	root, err := dynajson.NewByString(data)
	if err != nil {
		return nil, fmt.Errorf("dynajson.NewByString: %w", err)
	}

	return New(root)
}

// NewByPath ... func
func NewByPath(argPath string) (JSONSchema, error) {

	root, err := dynajson.NewByPath(argPath)
	if err != nil {
		return nil, fmt.Errorf("%s: dynajson.NewByPath: %w", argPath, err)
	}

	return New(root)
}

// ---------------------------------------------------------------------------

type refInfoType struct {
	selKey []interface{}
	getKey []string
	refVal string
}

func expandRef(root *dynajson.JSONElement) error {

	for {
		var refInfo *refInfoType

		err := root.Walk(func(parents []interface{}, key interface{}, val interface{}) (bool, error) {

			if key != "$ref" {
				return true, nil
			}

			//
			// map のキーが "$ref" のものを探し、その参照先とともに refInfo に設定する
			// このとき、複数のものを一度に処理すると置き換え済に対する参照が発生してしまうので
			// 一度の Walk() により行う検出と置換は 1 つのみとする
			//
			valStr, ok := val.(string)
			if !ok {
				return false, fmt.Errorf("%T: key=[%[1]v] val=[%v]: key-type not string", key, val)
			}

			if !strings.HasPrefix(valStr, "#/") {
				return false, fmt.Errorf("%s: illegal $rev (val[0] != '#')", valStr)
			}

			getKey := strings.Split(valStr[2:], "/")
			if len(getKey) == 0 {
				return false, fmt.Errorf("%s: illegal $rev (len(val) == 0)", valStr)
			}

			refInfo = &refInfoType{
				selKey: parents,
				getKey: getKey,
				refVal: valStr,
			}

			//fmt.Printf("%v %v\n", parents, getKey)

			// 一度に 1 つしか検出しない
			return false, nil
		})

		if err != nil {
			return fmt.Errorf("Walk: %w", err)
		}

		if refInfo == nil {
			break
		}

		//fmt.Printf("-%v-\n", refInfo.selKey)
		where := root.Select(refInfo.selKey)
		//fmt.Println(where)

		if where.IsNil() {
			continue
		}

		update := root.Select(refInfo.getKey)
		//fmt.Println(update)

		if update.IsNil() {
			continue
		}

		if !where.IsMap() {
			return fmt.Errorf("%v: where.IsMap: illegal type (%s): %T", refInfo.selKey, refInfo.refVal, where.Raw())
		}

		if !update.IsMap() {
			return fmt.Errorf("%v: update.IsMap: illegal type (%s): %T", refInfo.getKey, refInfo.refVal, update.Raw())
		}

		err = where.Delete("$ref")
		if err != nil {
			return fmt.Errorf("where.Delete($ref): %w", err)
		}

		where.Put("#original-ref#", refInfo.refVal)

		update.EachMap(func(k string, v *dynajson.JSONElement) (bool, error) {

			where.Put(k, v.Raw())
			return true, nil
		})
	}

	return nil
}
