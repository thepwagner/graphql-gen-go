package generator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thepwagner/graphql-gen-go/internal/generator"
	"github.com/thepwagner/graphql-gen-go/internal/test"
)

func TestGenerate_SWAPI(t *testing.T) {
	schema := test.LoadTestSchema(t, test.StarWars)
	output := test.SetupOutputDir(t, test.StarWars)

	const personByID = `
query PersonByID($id: ID!) {
 person(personID: $id) {
   name
   hairColor
   species {
     name
   }
 }
}`

	err := generator.Generate(output, schema, personByID)
	assert.NoError(t, err)
}
