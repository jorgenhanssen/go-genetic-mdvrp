package entities

import (
	"fmt"
	"math/rand"
	"reflect"
)

type Customer struct {
	X, Y            float64
	ID              int
	ServiceDuration float64
	Demand          float64
}

func (c *Customer) String() string {
	return fmt.Sprintf(`Customer(
  ID: %d
  Coords: [%.2f, %.2f]
  ServiceDuration: %f
  Demand: %.3f
)`, c.ID, c.X, c.Y, c.ServiceDuration, c.Demand)
}

func (c *Customer) GetPosition() (X, Y float64) {
	return c.X, c.Y
}

type Customers map[int]*Customer

func (cs Customers) RandomSelect() (k int, v *Customer) {
	mapKeys := reflect.ValueOf(cs).MapKeys()
	selectedKey := mapKeys[rand.Intn(len(mapKeys))].Interface().(int)
	return selectedKey, cs[selectedKey]
}

func (cs Customers) String() string {
	text := ""
	for _, customer := range cs {
		text += fmt.Sprintf("%v\n", customer)
	}
	return text
}
