// Code generated by graphql-gen-go. DO NOT EDIT.
package generated

type Person struct {
	Name      string
	BirthYear *string
	EyeColor  *string
	Gender    *string
	HairColor string
	Height    *int64
	Mass      *int64
	SkinColor *string
	Homeworld *Planet
	Films     []*Film
	Species   []*Species
	Starships []*Starship
	Vehicles  []*Vehicle
	Created   *string
	Edited    *string
	Id        string
}

func (r Person) GetName() string {
	return r.Name
}

func (r Person) GetBirthYear() *string {
	return r.BirthYear
}

func (r Person) GetEyeColor() *string {
	return r.EyeColor
}

func (r Person) GetGender() *string {
	return r.Gender
}

func (r Person) GetHairColor() string {
	return r.HairColor
}

func (r Person) GetHeight() *int64 {
	return r.Height
}

func (r Person) GetMass() *int64 {
	return r.Mass
}

func (r Person) GetSkinColor() *string {
	return r.SkinColor
}

func (r Person) GetHomeworld() *Planet {
	return r.Homeworld
}

func (r Person) GetFilms() []*Film {
	return r.Films
}

func (r Person) GetSpecies() []*Species {
	return r.Species
}

func (r Person) GetStarships() []*Starship {
	return r.Starships
}

func (r Person) GetVehicles() []*Vehicle {
	return r.Vehicles
}

func (r Person) GetCreated() *string {
	return r.Created
}

func (r Person) GetEdited() *string {
	return r.Edited
}

func (r Person) GetId() string {
	return r.Id
}

func (r Person) GetPersonByIDPersonSpecies() []PersonByIDPersonSpecies {
	ret := make([]PersonByIDPersonSpecies, len(r.Species))
	for i, o := range r.Species {
		ret[i] = o
	}
	return ret
}
