package entities

import (
	"fmt"
	"math/rand"
	"reflect"
)

type Depot struct {
	X, Y             float64
	MaxNumVehicles   int
	MaxRouteDuration float64
	MaxVehicleLoad   float64
}

func (d *Depot) String() string {
	return fmt.Sprintf(`Depot(
  Coords: [%.2f, %.2f]
  MaxNumVehicles: %d
  MaxRouteDuration: %.3f
  MaxVehicleLoad: %.3f
)`, d.X, d.Y, d.MaxNumVehicles, d.MaxRouteDuration, d.MaxVehicleLoad)
}

func (d *Depot) GetPosition() (X, Y float64) {
	return d.X, d.Y
}

type Depots map[int]*Depot

func (ds Depots) RandomSelect() (k int, v *Depot) {
	mapKeys := reflect.ValueOf(ds).MapKeys()
	selectedKey := mapKeys[rand.Intn(len(mapKeys))].Interface().(int)
	return selectedKey, ds[selectedKey]
}

func (ds Depots) String() string {
	text := ""
	for _, depot := range ds {
		text += fmt.Sprintf("%v\n", depot)
	}
	return text
}
