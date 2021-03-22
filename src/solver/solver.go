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

type GenerationInfo struct {
	BestAgent         *Agent
	GenerationNumber  int
	PopulationFitness Fitness
}

type SolverConfig struct {
	Depots          entities.Depots
	Customers       entities.Customers
	PopulationSize  int
	SelectionSize   float64
	NumCPUs         int
	SelectionMethod Selector

	RandomChanceRouteSplit              int
	RandomChanceDepotRelocation         int
	RandomChanceEvaluateOuterDepotRoute int
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
	if cfg.SelectionMethod == "" {
		cfg.SelectionMethod = Roulette
	}

	if cfg.RandomChanceRouteSplit == 0 {
		cfg.RandomChanceRouteSplit = 9999999999
	}
	if cfg.RandomChanceDepotRelocation == 0 {
		cfg.RandomChanceDepotRelocation = 9999999999
	}
	if cfg.RandomChanceEvaluateOuterDepotRoute == 0 {
		cfg.RandomChanceEvaluateOuterDepotRoute = 9999999999
	}

	return nil
}

type Solver struct {
	SolverConfig
	threads *threading.Instance

	agents     Agents
	generation int

	PostIterationCallback func(info GenerationInfo)
}

func NewSolver(cfg SolverConfig) (*Solver, error) {
	if err := cfg.ValidateAndSetDefaults(); err != nil {
		return nil, err
	}

	return &Solver{
		SolverConfig:          cfg,
		PostIterationCallback: func(info GenerationInfo) {},
		threads:               threading.New(threading.Config{NumThreads: cfg.NumCPUs}),
	}, nil

}

func (s *Solver) Solve(endCondition EndCondition) func() {
	abort := make(chan bool)
	go s.solve(endCondition, abort)

	return func() {
		abort <- true
	}
}

func (s *Solver) solve(endCondition EndCondition, abort chan bool) {
	s.initializeAgents()

	for ; ; s.generation++ {
		numNewAgents := int(float64(s.PopulationSize) * s.SelectionSize)

		s.threads.Run(func(tid int) error {
			for i := tid; i < numNewAgents; i += s.threads.NumThreads {
				p1i, p1 := s.agents.SelectOne(Roulette)
				p2i, p2 := s.agents.SelectOne(Roulette)

				c1 := s.mate(p1, p2)
				if c1.Fitness.Total < p1.Fitness.Total {
					s.agents[p1i] = c1
				}

				c2 := s.mate(p2, p1)
				if c2.Fitness.Total < p2.Fitness.Total {
					s.agents[p2i] = c2
				}
			}

			return nil
		})

		s.onIterationEnd()

		select {
		case <-abort:
			return
		default:
			continue
		}
	}

}

func (s *Solver) mate(a, b *Agent) (child *Agent) {
	route := b.Dna.GetRandomRoute()

	child = a.Copy()
	child.InjectRoute(route, s)
	child.RandomMutation(s)

	child.Evaluate(s.Depots, s.Customers)

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
	info := GenerationInfo{
		BestAgent:        s.agents[0],
		GenerationNumber: s.generation,
	}

	for _, agent := range s.agents {
		info.PopulationFitness.Add(&agent.Fitness)
		if agent.Fitness.Total < info.BestAgent.Fitness.Total {
			info.BestAgent = agent
		}
	}

	s.PostIterationCallback(info)
}
