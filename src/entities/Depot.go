package entities

import "fmt"

type Depot struct {
	MaxNumVehicles int
	MaxRouteDuration float64
	MaxVehicleLoad float64
	X float64
	Y float64
}

func (d *Depot) String() string {
	return fmt.Sprintf(`Depot(
  Coords: [%.2f, %.2f]
  MaxNumVehicles: %d
  MaxRouteDuration: %.3f
  MaxVehicleLoad: %.3f
)`, d.X, d.Y, d.MaxNumVehicles, d.MaxRouteDuration, d.MaxVehicleLoad)
}

type Depots []*Depot

func (ds Depots) String() string {
	text := ""
	for _, depot := range ds {
		text += fmt.Sprintf("%v\n", depot)
	}
	return text
}