package solver

import (
	"fmt"
	"math/rand"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
)

// DNA describes a solution.
type DNA []*Route

// NewDNA creates a new random DNA where a depot's routes
// consist of customers closest to the depot.
func NewDNA(depots entities.Depots, customers entities.Customers) (dna DNA) {
	depotCustomers := make(map[int]entities.Customers)
	for depotID := range depots {
		depotCustomers[depotID] = make(entities.Customers)
	}
	for cID, customer := range customers {
		closestDepotID := 0
		closestDepotDistance := 999999999.0
		for dID, depot := range depots {
			dist := distance(depot, customer)
			if dist < closestDepotDistance {
				closestDepotDistance = dist
				closestDepotID = dID
			}
		}
		depotCustomers[closestDepotID][cID] = customer
	}

	for depotID, remainingCustomers := range depotCustomers {
		depotRoutes := []*Route{}
		for j := 0; j < depots[depotID].MaxNumVehicles; j++ {
			depotRoutes = append(depotRoutes, &Route{DepotID: depotID})
		}

		for i := 0; len(remainingCustomers) != 0; i = (i + 1) % len(depotRoutes) {
			cID, customer := remainingCustomers.RandomSelect()
			delete(remainingCustomers, cID)
			nucleotide := depotRoutes[i]
			nucleotide.Path = append(nucleotide.Path, customer.ID)
		}

		dna = append(dna, depotRoutes...)
	}

	return
}

// String returns a print-friendly description of the dna.
func (dna DNA) String() string {
	text := ""
	for _, route := range dna {
		text += fmt.Sprintf("%v\n", route)
	}
	return text
}

// GetRandomRoute returns a random route in the dna.
func (dna DNA) GetRandomRoute() *Route {
	return dna[rand.Int63n(int64(len(dna)))]
}

// RemoveRouteNodes removes all customers from the dna that
// is also found in the provided route.
func (dna DNA) RemoveRouteNodes(route *Route) {
	for _, cID := range route.Path {
		for _, _route := range dna {
			index := -1
			for i, _cID := range _route.Path {
				if cID == _cID {
					index = i
					break
				}
			}
			if index >= 0 {
				_route.Path = append(_route.Path[:index], _route.Path[index+1:]...)
				break
			}
		}
	}
}

//  Route is a route that connects customers.
// This is the path vehicles will travel.
type Route struct {
	DepotID int
	Path    []int
}

// GetRandomRoute returns a random route in the route.
func (route Route) String() string {
	text := ""

	text += fmt.Sprintf("(%d", route.DepotID)
	for _, cID := range route.Path {
		text += fmt.Sprintf(" %d ", cID)
	}
	text += fmt.Sprint(route.DepotID)
	text += ")"

	return text
}
