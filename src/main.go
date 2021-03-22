package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/solver"
	"github.com/jorgenhanssen/go-genetic-mdvrp/src/visualizer"
)

func main() {
	gui, err := visualizer.New()
	if err != nil {
		panic(err)
	}

	go solveProblemsInFolder("problems", gui)

	gui.Run()
}

func solveProblemsInFolder(problemDirPath string, gui *visualizer.Instance) {
	files, err := ioutil.ReadDir(problemDirPath)
	if err != nil {
		panic(err)
	}

	for i, file := range files {
		if i+1 < 23 {
			continue
		}
		if file.IsDir() {
			continue
		}

		fmt.Printf("Solving %s\n", file.Name())
		depots, customers, err := LoadProblem(fmt.Sprintf("%s/%s", problemDirPath, file.Name()))
		if err != nil {
			panic(err)
		}

		slvr, err := solver.NewSolver(solver.SolverConfig{
			PopulationSize: 256,
			Depots:         depots,
			Customers:      customers,
			// SelectionSize:  0.5,
			// NumCPUs:   1,
		})
		if err != nil {
			panic(err)
		}

		slvr.PostIterationCallback = func(best *solver.Agent) {
			fmt.Printf("Best fitness: (dist:%.3f od:%.3f)\n", best.Fitness.Distance, best.Fitness.OverDemand)
			gui.Draw(depots, customers, best)
		}

		slvr.Solve(solver.EndCondition{
			Distance: 0,
		})

		// time.Sleep(time.Millisecond * 1000)
		break
	}

	fmt.Println("Problems solved!")
	gui.Stop()
}
