package solver

import (
	"fmt"
	"runtime"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
	"github.com/jorgenhanssen/go-genetic-mdvrp/src/threading"
)

type EndCondition struct {
	Distance float64
}

type SolverConfig struct {
	PopulationSize int
	SelectionSize  float64
	Depots         entities.Depots
	Customers      entities.Customers
	NumCPUs        int
}

func (cfg *SolverConfig) ValidateAndSetDefaults() error {
	if len(cfg.Depots) == 0 {
		return fmt.Errorf("No depots provided")
	}
	if len(cfg.Customers) == 0 {
		return fmt.Errorf("No customers provided")
	}

	if cfg.PopulationSize == 0 {
		cfg.SelectionSize = 200
	}
	if cfg.NumCPUs == 0 {
		cfg.NumCPUs = runtime.NumCPU()
	}
	if cfg.SelectionSize == 0 {
		cfg.SelectionSize = 0.3
	}

	return nil
}

type Solver struct {
	SolverConfig
	threads *threading.Instance

	agents     Agents
	generation int

	PostIterationCallback func(best *Agent)
}

func NewSolver(cfg SolverConfig) (*Solver, error) {
	if err := cfg.ValidateAndSetDefaults(); err != nil {
		return nil, err
	}

	return &Solver{
		SolverConfig:          cfg,
		PostIterationCallback: func(best *Agent) {},
		threads:               threading.New(threading.Config{NumThreads: cfg.NumCPUs}),
	}, nil

}

func (s *Solver) Solve(endCondition EndCondition) {
	s.solve(endCondition)
}

func (s *Solver) solve(endCondition EndCondition) {
	s.initializeAgents()

	for ; ; s.generation++ {
		numNewAgents := int(float64(s.PopulationSize) * s.SelectionSize)

		s.threads.Run(func(tid int) error {
			for i := tid; i < numNewAgents; i += s.threads.NumThreads {
				p1i, p1 := s.agents.SelectOne(Roulette)
				_, p2 := s.agents.SelectOne(Roulette)

				c1 := s.mate(p1, p2)
				if c1.Fitness.Total < p1.Fitness.Total {
					s.agents[p1i] = c1
				}
			}

			return nil
		})

		s.onIterationEnd()

		// if s.generation > 20 {
		// 	return
		// }
	}

}

func (s *Solver) mate(a, b *Agent) (child *Agent) {
	route := b.Dna.GetRandomRoute()

	child = a.Copy()
	child.InjectRoute(route, s)
	child.RandomMutation(s)

	// child

	// child.Evaluate(s.Depots, s.Customers)

	return
}

func (s *Solver) initializeAgents() {
	s.threads.Run(func(tid int) error {
		agents := Agents{}
		for i := 0; i < s.PopulationSize/s.NumCPUs; i++ {
			agents = append(agents, NewAgent(s.Depots, s.Customers))
		}

		s.threads.Lock()
		s.agents = append(s.agents, agents...)
		s.threads.Unlock()

		return nil
	})
}

func (s *Solver) onIterationEnd() {
	bestAgent := s.agents[0]
	for _, agent := range s.agents {
		if agent.Fitness.Total < bestAgent.Fitness.Total {
			bestAgent = agent
		}
	}

	s.PostIterationCallback(bestAgent)
}
