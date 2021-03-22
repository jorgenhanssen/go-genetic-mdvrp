package solver

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/jinzhu/copier"
	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
)

type Fitness struct {
	Total      float64
	Distance   float64
	OverDemand float64
}

func (f *Fitness) CalculateTotal() {
	f.Total = f.Distance
	f.Total += 100 * math.Pow(f.OverDemand, 2)
}

func (f *Fitness) Add(f2 *Fitness) {
	f.Distance += f2.Distance
	f.OverDemand += f2.OverDemand
	f.CalculateTotal()
}

func (f Fitness) String() string {
	return fmt.Sprintf("Fitness(dist: %f, over-demand: %f, total: %f)", f.Distance, f.OverDemand, f.Total)
}

type Agent struct {
	Dna     DNA
	Fitness Fitness
}

func NewAgent(depots entities.Depots, customers entities.Customers) *Agent {
	agent := &Agent{
		Dna: NewDNA(depots, customers),
	}

	agent.Evaluate(depots, customers)

	return agent
}

func (agent *Agent) Evaluate(depots entities.Depots, customers entities.Customers) {
	agent.Fitness.Distance = 0
	agent.Fitness.OverDemand = 0
	for _, route := range agent.Dna {
		if len(route.Path) == 0 {
			continue
		}

		agent.Fitness.Distance += distance(depots[route.DepotID], customers[route.Path[0]])
		agent.Fitness.Distance += distance(depots[route.DepotID], customers[route.Path[len(route.Path)-1]])
		for i := 0; i < len(route.Path)-1; i++ {
			agent.Fitness.Distance += distance(customers[route.Path[i]], customers[route.Path[i+1]])
		}

		demand := 0.0
		for _, cID := range route.Path {
			demand += customers[cID].Demand
		}

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

func (agent *Agent) InjectRoute(injectedRoute *Route, s *Solver) {
	agent.Dna.RemoveRouteNodes(injectedRoute)

	for _, cID := range injectedRoute.Path {
		bestScore := 999999999999.0
		bestRoute := agent.Dna[0]
		bestI := 0

		for _, route := range agent.Dna {
			if route.DepotID != injectedRoute.DepotID && rand.Intn(s.RandomChanceEvaluateOuterDepotRoute) != 0 {
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

		// We want to relocate a split route's closest depot
		// 1/100 chance of re-locating closest depot
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
				lowestKey := -1
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

type Agents []*Agent

type Selector string

const (
	Roulette Selector = "Roulette"
	Random   Selector = "Random"
)

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

// TODO: MOVE OUT
func distance(a, b entities.Location) float64 {
	ax, ay := a.GetPosition()
	bx, by := b.GetPosition()
	return math.Sqrt(math.Pow(ax-bx, 2) + math.Pow(ay-by, 2))
}
