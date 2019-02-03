// Code generated by graphql-gen-go. DO NOT EDIT.
package generated

import (
	"encoding/json"
	"io"
)

func ReadPersonByID(r io.Reader) (PersonByID, error) {
	var ret personByID
	if err := json.NewDecoder(r).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

type PersonByID interface {
	GetPersonByIDPerson() PersonByIDPerson
}

type PersonByIDPerson interface {
	GetName() string
	GetHairColor() string
	GetPersonByIDPersonSpecies() []PersonByIDPersonSpecies
}

type PersonByIDPersonSpecies interface {
	GetName() string
}

type personByID struct {
	Person *Person `json:"person"`
}

func (r personByID) GetPersonByIDPerson() PersonByIDPerson {
	return r.Person
}
