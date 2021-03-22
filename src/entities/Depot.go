package entities

import (
	"fmt"
	"math/rand"
	"reflect"
)

// Depot describes a depot.
// Depots contain vehicles that visit
// customers.
type Depot struct {
	// Position of the depot.
	X, Y float64

	// MaxNumVehicles is how many vehicles/routes
	// can be dispatched from this depot.
	MaxNumVehicles int

	// MaxRouteDuration is how long a vehicle/route
	// can take when dispatched from this depot.
	// This is linked with customers service duration.
	MaxRouteDuration float64

	// MaxVehicleLoad is how much load a vehicle/route
	// can take when dispatched from this depot.
	// This is linked with customer's demands.
	MaxVehicleLoad float64
}

// String returns the stringified depot.
func (d *Depot) String() string {
	return fmt.Sprintf(`Depot(
  Coords: [%.2f, %.2f]
  MaxNumVehicles: %d
  MaxRouteDuration: %.3f
  MaxVehicleLoad: %.3f
)`, d.X, d.Y, d.MaxNumVehicles, d.MaxRouteDuration, d.MaxVehicleLoad)
}

// GetPosition returns the x and y position of the depot.
func (d *Depot) GetPosition() (X, Y float64) {
	return d.X, d.Y
}

// Depots is a map of multiple depots.
// The key is the depot's ID.
type Depots map[int]*Depot

// RandomSelect returns a random depot and its ID.
func (ds Depots) RandomSelect() (k int, v *Depot) {
	mapKeys := reflect.ValueOf(ds).MapKeys()
	selectedKey := mapKeys[rand.Intn(len(mapKeys))].Interface().(int)
	return selectedKey, ds[selectedKey]
}

// String prints a collection of depots
func (ds Depots) String() string {
	text := ""
	for _, depot := range ds {
		text += fmt.Sprintf("%v\n", depot)
	}
	return text
}
