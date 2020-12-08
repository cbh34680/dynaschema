package dynaschema

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	TestValidBody(t)
}

func TestValidBody(t *testing.T) {

	assert := assert.New(t)

	for _, verdir := range []string{"v2.0", "v3.0"} {

		schemaPath := filepath.Join("testdata", verdir, "petstore-expanded.json")

		schema, err := NewByPath(schemaPath)
		assert.Nil(err)

		result, err := schema.ValidateJSONRequestBody("/pets", "post", `{"name":"Tama, Tanaka"}`)
		assert.Nil(err)
		assert.True(result.Valid()) // True

		result, err = schema.ValidateJSONRequestBody("/pets", "post", `{"NAME":"Tama, Tanaka"}`)
		assert.Nil(err)
		assert.False(result.Valid()) // False

	}
}

func TestValidQuery(t *testing.T) {

	assert := assert.New(t)

	for _, verdir := range []string{"v2.0", "v3.0"} {

		schemaPath := filepath.Join("testdata", verdir, "petstore-expanded.json")

		schema, err := NewByPath(schemaPath)
		assert.Nil(err)

		result, err := schema.ValidateJSONRequestBody("/pets", "get", `{}`)
		assert.Nil(err)
		assert.True(result.Valid()) // True

		result, err = schema.ValidateJSONRequestBody("/pets", "get", `{}`)
		assert.Nil(err)
		assert.False(result.Valid()) // False

	}
}
