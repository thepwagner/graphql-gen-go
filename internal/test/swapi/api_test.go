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
	// Pre:
	// - have a GraphQL schema (this uses SWAPI)
	// - run graphql-gen-go against it

	// Model objects created for all types in the schema
	// Create and link some together for test data
	newHope := &generated.Film{Title: "A New Hope"}
	phantom := &generated.Film{Title: "Phantom Menace"}
	droid := &generated.Species{Name: "Droid"}
	human := &generated.Species{Name: "Human"}
	r2 := generated.Person{
		Name:      "R2D2",
		HairColor: "Chrome",
		Species:   []*generated.Species{droid},
		Films:     []*generated.Film{newHope, phantom},
	}
	hanSolo := generated.Person{
		Name:      "Han Solo",
		HairColor: "Brown",
		Species:   []*generated.Species{human},
		Films:     []*generated.Film{newHope},
	}

	// Mock a JSON response to GraphQL query:
	// TODO: generate boilerplate to make doing this in a httptest server ezpz
	// TODO: generate client to make performing this against an http server ezpz
	personJSON, err := json.Marshal(map[string]generated.Person{"person": r2})
	require.NoError(t, err)

	// Use generated deser to read the mocked response:
	personByID, err := generated.ReadPersonByID(bytes.NewReader(personJSON))
	require.NoError(t, err)

	// Data is backed by model *generated.Person, wrapped in an interface that limits scope to requested fields:
	person := personByID.GetPersonByIDPerson()
	assert.Equal(t, "R2D2", person.GetName())
	assert.Equal(t, "Chrome", person.GetHairColor())

	// Problem: Golang slices are not covariant
	// This is handled by generating custom accessors on *generated.Person for each view
	assert.Equal(t, []generated.PersonByIDPersonSpecies{droid}, person.GetPersonByIDPersonSpecies())

	// AllPerson
	personsJSON, err := json.Marshal(map[string][]generated.Person{"allPeople": {r2, hanSolo}})
	require.NoError(t, err)
	allPeople, err := generated.ReadAllPeople(bytes.NewReader(personsJSON))
	require.NoError(t, err)
	assert.Len(t, allPeople.GetAllPeopleAllPeople(), 2)
	for _, person := range allPeople.GetAllPeopleAllPeople() {
		if person.GetName() == "R2D2" {
			for _, film := range person.GetAllPeopleAllPeopleFilms() {
				assert.Contains(t, []string{"A New Hope", "Phantom Menace"}, film.GetTitle())
			}
		} else if person.GetName() == "Han Solo" {
			for _, film := range person.GetAllPeopleAllPeopleFilms() {
				assert.Contains(t, []string{"A New Hope"}, film.GetTitle())
			}
		} else {
			assert.Fail(t, "unknown person")
		}
	}

}
