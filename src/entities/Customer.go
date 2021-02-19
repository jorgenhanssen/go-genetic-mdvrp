package entities

import "fmt"

type Customer struct {
	ID int
	X float64
	Y float64
	ServiceDuration float64
	Demand float64
}

func (c *Customer) String() string {
	return fmt.Sprintf(`Customer(
  ID: %d
  Coords: [%.2f, %.2f]
  ServiceDuration: %f
  Demand: %.3f
)`, c.ID, c.X, c.Y, c.ServiceDuration, c.Demand)
}

type Customers []*Customer

func (cs Customers) String() string {
	text := ""
	for _, customer := range cs {
		text += fmt.Sprintf("%v\n", customer)
	}
	return text
}