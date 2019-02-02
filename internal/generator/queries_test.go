package generator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQueries(t *testing.T) {
	gen := NewTestGenerator(t)

	err := gen.GenerateQueries(`
query PersonByID($id: ID!) {
 person(personID: $id) {
   name
   hairColor
   species {
     name
   }
 }
}`)
	assert.NoError(t, err)
}
