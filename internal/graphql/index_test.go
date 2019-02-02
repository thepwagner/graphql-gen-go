package graphql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thepwagner/magenny/internal/graphql"
	"github.com/thepwagner/magenny/internal/test"
)

func TestIndexSchema_SWAPI(t *testing.T) {
	schema := test.LoadTestSchema(t, test.StarWars)

	indexed := graphql.IndexSchema(schema)
	assert.Equal(t, "RootQuery", indexed.Query.Name)

	// A type:
	if person := indexed.TypeIndex["Person"]; assert.NotNil(t, person) {
		assert.Len(t, person.FieldIndex, 16)
		if personName := person.FieldIndex["name"]; assert.NotNil(t, personName) {
			assert.NotNil(t, personName.Definition)

			//   name: String!
			fieldType := personName.FieldType
			assert.True(t, fieldType.Primitive)
			assert.True(t, fieldType.NonNull)
			assert.False(t, fieldType.List)
			assert.Equal(t, "string", fieldType.Go)
		}
	}

	// An interface:
	if node := indexed.TypeIndex["Node"]; assert.NotNil(t, node) {
		assert.Len(t, node.FieldIndex, 1)
		if nodeID := node.FieldIndex["id"]; assert.NotNil(t, nodeID) {
			assert.NotNil(t, nodeID.Definition)
		}
	}
}
