package main

import (
	"fmt"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/solver"
	"github.com/jorgenhanssen/go-genetic-mdvrp/src/visualizer"
)

func main() {
	gui, err := visualizer.New()
	if err != nil {
		panic(err)
	}

	go solveProblem("problems/p23", gui)

	gui.Run()
}

func solveProblem(path string, gui *visualizer.Instance) {
	depots, customers, err := LoadProblem(path)
	if err != nil {
		panic(err)
	}

	slvr, err := solver.NewSolver(solver.SolverConfig{
		Depots:    depots,
		Customers: customers,

		PopulationSize:  128,
		SelectionSize:   0.5,
		SelectionMethod: solver.Roulette,

		// 1/n chances:
		RandomChanceRouteSplit:              20,
		RandomChanceDepotRelocation:         50,
		RandomChanceEvaluateOuterDepotRoute: 100000,
	})
	if err != nil {
		panic(err)
	}

	slvr.PostIterationCallback = func(info solver.GenerationInfo) {
		fmt.Printf("%s (generation %d)\n", path, info.GenerationNumber)
		fmt.Printf("\tBest error:  %v\n", info.BestAgent.Fitness)
		fmt.Printf("\tTotal error: %v\n\n", info.PopulationFitness)
		gui.Draw(depots, customers, info.BestAgent)
	}

	_ = slvr.Solve(solver.EndCondition{})
}
