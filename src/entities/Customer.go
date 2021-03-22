package entities

import (
	"fmt"
	"math/rand"
	"reflect"
)

// Customer describes a customer objects.
// Customers are nodes in the graph that
// vehicles from depots visit.
type Customer struct {
	// Position of the customer.
	X, Y float64

	// ID of the customer.
	ID int

	// ServiceDuration is how long a visit will
	// take for this customer.
	// This is linked with depots' max route duration.
	ServiceDuration float64

	// Demand is how much demand is required
	// for this customer.
	// This is linked with depots' max vehicle load.
	Demand float64
}

// String returns the stringified customer.
func (c *Customer) String() string {
	return fmt.Sprintf(`Customer(
  ID: %d
  Coords: [%.2f, %.2f]
  ServiceDuration: %f
  Demand: %.3f
)`, c.ID, c.X, c.Y, c.ServiceDuration, c.Demand)
}

// GetPosition returns the x and y position of the customer.
func (c *Customer) GetPosition() (X, Y float64) {
	return c.X, c.Y
}

// Customers is a map of multiple customers.
// The key is the customer's ID.
type Customers map[int]*Customer

// RandomSelect returns a random customer and their ID.
func (cs Customers) RandomSelect() (k int, v *Customer) {
	mapKeys := reflect.ValueOf(cs).MapKeys()
	selectedKey := mapKeys[rand.Intn(len(mapKeys))].Interface().(int)
	return selectedKey, cs[selectedKey]
}

// String prints a collection of customers
func (cs Customers) String() string {
	text := ""
	for _, customer := range cs {
		text += fmt.Sprintf("%v\n", customer)
	}
	return text
}
