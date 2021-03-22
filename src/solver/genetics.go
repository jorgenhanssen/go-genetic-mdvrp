package solver

import (
	"fmt"
	"math/rand"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
)

type DNA []*Route

func NewDNA(depots entities.Depots, customers entities.Customers) (dna DNA) {
	for i, depot := range depots {
		for j := 0; j < depot.MaxNumVehicles; j++ {
			dna = append(dna, &Route{DepotID: i})
		}
	}

	remainingCustomers := make(entities.Customers)
	for k, v := range customers {
		remainingCustomers[k] = v
	}

	for i := 0; len(remainingCustomers) != 0; i = (i + 1) % len(dna) {
		cID, customer := remainingCustomers.RandomSelect()
		delete(remainingCustomers, cID)
		nucleotide := dna[i]
		nucleotide.Path = append(nucleotide.Path, customer.ID)
	}

	return
}

func (dna DNA) String() string {
	text := ""
	for _, route := range dna {
		text += fmt.Sprintf("%v\n", route)
	}
	return text
}

func (dna DNA) GetRandomRoute() *Route {
	return dna[rand.Int63n(int64(len(dna)))]
}

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

func (dna DNA) depotIsAvailable(s *Solver, id int) bool {
	max := s.Depots[id].MaxNumVehicles

	numOccurences := 1 // if we add a new one
	for _, route := range dna {
		if route.DepotID == id {
			numOccurences++
			if numOccurences > max {
				return false
			}
		}
	}

	return true
}

func (dna DNA) availableDepot(s *Solver, biasID int) (int, error) {
	if dna.depotIsAvailable(s, biasID) {
		return biasID, nil
	}
	for i := range s.Depots {
		if i != biasID && dna.depotIsAvailable(s, i) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("No available depots")
}

type Route struct {
	DepotID int
	Path    []int
}

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
