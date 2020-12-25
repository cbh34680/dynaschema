package dynaschema

import (
	"fmt"
	"strings"

	"github.com/cbh34680/dynajson"
)

// NewOption ... int
type NewOption int

const (
	// NoOption ... 0
	NoOption NewOption = 0
	// ExpandRef ... 1
	ExpandRef NewOption = 1 << iota
	// SetMinlenIfRequired ... 2
	SetMinlenIfRequired
)

// JSONSchema ... struct
type JSONSchema interface {
	RawJSON() *dynajson.JSONElement
	String() string
	/*
		SetFlag(string, bool)
		GetServers() ([]string, error)
	*/
	FindParameter(string, string, string) (string, error)
	FindBody(string, string, string) (string, error)
}

// SchemaAbstract ... struct
type SchemaAbstract struct {
	objJSON *dynajson.JSONElement
	//flagMap map[string]bool
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
		return fmt.Errorf("%s, %s: not found", argPath, argMethod)
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

type findParameterHelperCallback func(*dynajson.JSONElement) map[string]interface{}

func (me *SchemaAbstract) findParameterHelper(argPath, argMethod, argIn string, callback findParameterHelperCallback) (string, error) {

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

		property := callback(spec)

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

// ---------------------------------------------------------------------------

// New ... func
func New(root *dynajson.JSONElement, opt NewOption) (JSONSchema, error) {

	ver := root.Select("openapi").AsString()
	if ver == "" {
		ver = root.Select("swagger").AsString()
		if ver == "" {
			return nil, fmt.Errorf("UnSupported Schema")
		}
	}

	if (opt & SetMinlenIfRequired) != 0 {

		err := setMinlenIfRequired(root)
		if err != nil {
			return nil, fmt.Errorf("setMinlenIfRequired: %w", err)
		}
	}

	if (opt & ExpandRef) != 0 {

		err := expandRef(root)
		if err != nil {
			return nil, fmt.Errorf("expandRef: %w", err)
		}
	}

	//fmt.Println(root)

	root.Readonly = true

	switch {
	case strings.HasPrefix(ver, "2.0"):
		return NewSchemaSwagger20(root), nil
	case strings.HasPrefix(ver, "3.0"):
		return NewSchemaOpenAPI30(root), nil
	}

	return nil, fmt.Errorf("UnSupported Schema-Version")
}

// NewByBytes ... func
func NewByBytes(data []byte, opt NewOption) (JSONSchema, error) {

	root, err := dynajson.NewByBytes(data)
	if err != nil {
		return nil, fmt.Errorf("dynajson.NewByBytes: %w", err)
	}

	return New(root, opt)
}

// NewByString ... func
func NewByString(data string, opt NewOption) (JSONSchema, error) {

	root, err := dynajson.NewByString(data)
	if err != nil {
		return nil, fmt.Errorf("dynajson.NewByString: %w", err)
	}

	return New(root, opt)
}

// NewByPath ... func
func NewByPath(argPath string, opt NewOption) (JSONSchema, error) {

	root, err := dynajson.NewByPath(argPath)
	if err != nil {
		return nil, fmt.Errorf("%s: dynajson.NewByPath: %w", argPath, err)
	}

	return New(root, opt)
}

// ---------------------------------------------------------------------------

func setMinlenIfRequired(root *dynajson.JSONElement) error {

	err := root.Walk(func(parents []interface{}, key interface{}, val interface{}) (bool, error) {

		//
		// required が true であり、schema 中の type が string であり、minLength が設定されて
		// いない場合は minLength=1 を付与する
		//
		// --> /exec?key=&val=10 のときに、key を有効な値として判定させないため
		//
		anyMap, ok := val.(map[string]interface{})
		if !ok {
			return true, nil
		}

		pRequired, ok := anyMap["required"]
		if !ok {
			return true, nil
		}

		bRequired, ok := pRequired.(bool)
		if !ok {
			return true, nil
		}

		if !bRequired {
			return true, nil
		}

		pSchema, ok := anyMap["schema"]
		if !ok {
			return true, nil
		}

		mSchema, ok := pSchema.(map[string]interface{})
		if !ok {
			return true, nil
		}

		pType, ok := mSchema["type"]
		if !ok {
			return true, nil
		}

		sType, ok := pType.(string)
		if !ok {
			return true, nil
		}

		if sType != "string" {
			return true, nil
		}

		if _, ok := mSchema["minLength"]; ok {
			return true, nil
		}

		mSchema["minLength"] = 1

		return true, nil
	})

	if err != nil {
		return fmt.Errorf("Walk: %w", err)
	}

	return nil
}

func expandRef(root *dynajson.JSONElement) error {

	type refInfoType struct {
		selKey []interface{}
		getKey []string
		refVal string
	}

	for {
		var refInfo *refInfoType

		err := root.Walk(func(parents []interface{}, key interface{}, val interface{}) (bool, error) {

			//
			// map のキーが "$ref" のものを探し、その参照先とともに refInfo に設定する
			// このとき、複数のものを一度に処理すると置き換え済に対する参照が発生してしまうので
			// 一度の Walk() により行う検出と置換は 1 つのみとする
			//
			if key != "$ref" {
				return true, nil
			}

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

		where := root.Select(refInfo.selKey)
		if where.IsNil() {
			continue
		}

		update := root.Select(refInfo.getKey)
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
		//where.Put("#original-ref#", refInfo.refVal)

		update.EachMap(func(k string, v *dynajson.JSONElement) (bool, error) {

			where.Put(k, v.Raw())
			return true, nil
		})
	}

	return nil
}
