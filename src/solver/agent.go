package solver

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/jinzhu/copier"
	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
)

// Agent is a solution descrption to the problem.
// It contains a DNA (direct route description) and
// a fitness score.
type Agent struct {
	Dna     DNA
	Fitness Fitness
}

// NewAgent creates a new random agent and evaluates the agent.
func NewAgent(s *Solver) *Agent {
	agent := &Agent{
		Dna: NewDNA(s.Depots, s.Customers),
	}

	agent.Evaluate(s.Depots, s.Customers)

	return agent
}

// Evaluate evaluates the fitness of the agent.
// The fitness is stored in the agent as a property.
func (agent *Agent) Evaluate(depots entities.Depots, customers entities.Customers) {
	agent.Fitness.Clear()

	for _, route := range agent.Dna {
		if len(route.Path) == 0 {
			continue
		}

		// Add depot -> c_1 and depot -> c_n
		agent.Fitness.Distance += distance(depots[route.DepotID], customers[route.Path[0]])
		agent.Fitness.Distance += distance(depots[route.DepotID], customers[route.Path[len(route.Path)-1]])

		// Add c_1 -> c_2, c_2 -> c_3, ... , c_n-1 -> c_n.
		for i := 0; i < len(route.Path)-1; i++ {
			agent.Fitness.Distance += distance(customers[route.Path[i]], customers[route.Path[i+1]])
		}

		// Accumulate all demand for the route.
		demand := 0.0
		for _, cID := range route.Path {
			demand += customers[cID].Demand
		}

		// Add the positive difference between the max load and demand.
		// If there is no positive difference the over-demand is 0.
		agent.Fitness.OverDemand += math.Max(demand-depots[route.DepotID].MaxVehicleLoad, 0)
	}

	agent.Fitness.CalculateTotal()
}

func (a *Agent) Copy() (child *Agent) {
	child = &Agent{
		Fitness: Fitness{
			Total:      a.Fitness.Total,
			Distance:   a.Fitness.Distance,
			OverDemand: a.Fitness.OverDemand,
		},
	}

	for _, route := range a.Dna {
		var copyPath []int
		copier.Copy(&copyPath, &route.Path)
		child.Dna = append(child.Dna, &Route{
			DepotID: route.DepotID,
			Path:    copyPath,
		})
	}

	return
}

// InjectRoute injects a route into its best placement.
// The injected route is decomposed and fitted into the existing
// routes so that the best overall per-new-customer is achieved.
func (agent *Agent) InjectRoute(injectedRoute *Route, s *Solver) {
	agent.Dna.RemoveRouteNodes(injectedRoute)

	for _, cID := range injectedRoute.Path {
		bestScore := 999999999999.0
		bestRoute := agent.Dna[0]
		bestI := 0

		for _, route := range agent.Dna {
			if route.DepotID != injectedRoute.DepotID && rand.Intn(s.RandomChanceEvaluateOuterDepotRoute) != 0 {
				// In most cases, we do not bother checking routes that
				// do not belong to the injected route's depot.
				// i.e: each route is (often) closest to the depot it is connected to.
				// therefore, we should skip evaluating other routes "far away".
				// That being said, in some cases, this is not the cases.
				// Which is why we have a small chance of checking outer depot routes.
				continue
			}
			for i := 0; i < len(route.Path); i++ {
				route.Path = append(route.Path[:i+1], route.Path[i:]...)
				route.Path[i] = cID

				agent.Evaluate(s.Depots, s.Customers)
				if agent.Fitness.Total < bestScore {
					bestScore = agent.Fitness.Total
					bestRoute = route
					bestI = i
				}

				// remove inserted point
				route.Path = append(route.Path[:i], route.Path[i+1:]...)
			}

		}

		bestRoute.Path = append(bestRoute.Path[:bestI+1], bestRoute.Path[bestI:]...)
		bestRoute.Path[bestI] = cID
	}
}

// RandomMutation runs a mutation procedure on the agent.
// This allows chances for the following mutations:
// - splitting a route in two
// - re-locating a route's depot
// all mutations follow constraints.
func (agent *Agent) RandomMutation(s *Solver) {
	// FIXME: hardcoded chance
	for _, route := range agent.Dna {
		if len(route.Path) == 0 {
			continue
		}

		hasBeenSplit := false
		if rand.Intn(s.RandomChanceRouteSplit) == 0 {
			availableDepotID, err := agent.availableDepot(s, route.DepotID)
			if err != nil {
				continue
			}

			splitPoint := len(route.Path) / 2
			splitRoute := Route{
				DepotID: availableDepotID,
				Path:    route.Path[:splitPoint],
			}
			route.Path = route.Path[splitPoint:]
			agent.Dna = append(agent.Dna, &splitRoute)

			hasBeenSplit = true
		}

		// If we have split the path, we want to ensure that this path
		// is connected to its closest depot (if available).
		if hasBeenSplit || rand.Intn(s.RandomChanceDepotRelocation) == 1 {
			m := map[int]float64{}
			for i, depot := range s.Depots {
				if i == route.DepotID {
					continue
				}
				m[i] = 0
				for _, cID := range route.Path {
					m[i] += distance(depot, s.Customers[cID])
				}
			}

			for len(m) > 0 {
				lowestVal := 9999999999.0
				lowestKey := 0
				for k, v := range m {
					if v < lowestVal {
						lowestKey = k
					}
				}
				if agent.depotIsAvailable(s, lowestKey) {
					route.DepotID = lowestKey
					break
				}

				delete(m, lowestKey)
			}
		}
	}
}

// depotIsAvailable checks if the provided depot id
// has available routes.
func (agent *Agent) depotIsAvailable(s *Solver, id int) bool {
	max := s.Depots[id].MaxNumVehicles

	numOccurences := 1 // if we add a new one
	for _, route := range agent.Dna {
		if route.DepotID == id {
			numOccurences++
			if numOccurences > max {
				return false
			}
		}
	}

	return true
}

// availableDepot returns an arbitrary depot that can be given
// more routes. If there is no available depots, an error will
// be returned.
func (agent *Agent) availableDepot(s *Solver, biasID int) (int, error) {
	if agent.depotIsAvailable(s, biasID) {
		return biasID, nil
	}
	for i := range s.Depots {
		if i != biasID && agent.depotIsAvailable(s, i) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("No available depots")
}

// Agents is a collection of agents.
type Agents []*Agent

// SelectOne selects an agent from the collection
// using the provided selector as the selection method.
func (agents Agents) SelectOne(method Selector) (int, *Agent) {
	switch method {
	case Roulette:
		{
			sum := 0.0
			highest := 0.0
			for _, agent := range agents {
				if agent.Fitness.Total > highest {
					highest = agent.Fitness.Total
				}
				sum += agent.Fitness.Total
			}

			value := rand.Float64() * sum
			for i, agent := range agents {
				value -= (highest - agent.Fitness.Total)
				if value <= 0 {
					return i, agent
				}
			}
		}
	case Random: // random by default
	}

	index := rand.Intn(len(agents))
	return index, agents[index]
}
