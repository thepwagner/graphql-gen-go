package swapi_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thepwagner/graphql-gen-go/internal/test/swapi/generated"
)

func TestAPI(t *testing.T) {
	droid := &generated.Species{Name: "Droid"}
	jawa := &generated.Species{Name: "Jawa"}
	r2 := generated.Person{
		Name:      "R2D2",
		HairColor: "Chrome",
		Species:   []*generated.Species{droid, jawa},
	}
	personJSON, err := json.Marshal(map[string]generated.Person{"person": r2})
	require.NoError(t, err)

	personByID, err := generated.ReadPersonByID(bytes.NewReader(personJSON))
	require.NoError(t, err)

	person := personByID.GetPersonByIDPerson()
	assert.Equal(t, "R2D2", person.GetName())
	assert.Equal(t, "Chrome", person.GetHairColor())
	assert.Equal(t, []generated.PersonByIDPersonSpecies{droid, jawa}, person.GetPersonByIDPersonSpecies())
}
