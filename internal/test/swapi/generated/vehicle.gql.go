// Code generated by graphql-gen-go. DO NOT EDIT.
package generated

type Vehicle struct {
	Name                 *string   `json:"name,omitempty"`
	Model                *string   `json:"model,omitempty"`
	VehicleClass         *string   `json:"vehicleClass,omitempty"`
	Manufacturers        []*string `json:"manufacturers,omitempty"`
	CostInCredits        *int64    `json:"costInCredits,omitempty"`
	Length               *float64  `json:"length,omitempty"`
	Crew                 *string   `json:"crew,omitempty"`
	Passengers           *string   `json:"passengers,omitempty"`
	MaxAtmospheringSpeed *int64    `json:"maxAtmospheringSpeed,omitempty"`
	CargoCapacity        *int64    `json:"cargoCapacity,omitempty"`
	Consumables          *string   `json:"consumables,omitempty"`
	Pilots               []*Person `json:"pilots,omitempty"`
	Films                []*Film   `json:"films,omitempty"`
	Created              *string   `json:"created,omitempty"`
	Edited               *string   `json:"edited,omitempty"`
	Id                   string    `json:"id,omitempty"`
}

func (r Vehicle) GetName() *string {
	return r.Name
}

func (r Vehicle) GetModel() *string {
	return r.Model
}

func (r Vehicle) GetVehicleClass() *string {
	return r.VehicleClass
}

func (r Vehicle) GetManufacturers() []*string {
	return r.Manufacturers
}

func (r Vehicle) GetCostInCredits() *int64 {
	return r.CostInCredits
}

func (r Vehicle) GetLength() *float64 {
	return r.Length
}

func (r Vehicle) GetCrew() *string {
	return r.Crew
}

func (r Vehicle) GetPassengers() *string {
	return r.Passengers
}

func (r Vehicle) GetMaxAtmospheringSpeed() *int64 {
	return r.MaxAtmospheringSpeed
}

func (r Vehicle) GetCargoCapacity() *int64 {
	return r.CargoCapacity
}

func (r Vehicle) GetConsumables() *string {
	return r.Consumables
}

func (r Vehicle) GetPilots() []*Person {
	return r.Pilots
}

func (r Vehicle) GetFilms() []*Film {
	return r.Films
}

func (r Vehicle) GetCreated() *string {
	return r.Created
}

func (r Vehicle) GetEdited() *string {
	return r.Edited
}

func (r Vehicle) GetId() string {
	return r.Id
}
