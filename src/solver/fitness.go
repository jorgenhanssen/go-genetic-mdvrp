package solver

import (
	"fmt"
	"math"
)

// Fitness contains fitness information.
// Fitness does in this case equal error, meaning
// that it is inversely corelated to
// good-performing fitness.
type Fitness struct {
	Total      float64
	Distance   float64
	OverDemand float64
}

func (f *Fitness) Clear() {
	f.Total = 0
	f.Distance = 0
	f.OverDemand = 0
}

// CalculateTotal calculates the total error for the fitness
// given distance and over-demand.
func (f *Fitness) CalculateTotal() {
	f.Total = f.Distance
	f.Total += 100 * math.Pow(f.OverDemand, 2)
}

// Add adds a secondary fitness to this fitness.
func (f *Fitness) Add(f2 *Fitness) {
	f.Distance += f2.Distance
	f.OverDemand += f2.OverDemand
	f.CalculateTotal()
}

// String returns a print-friendly string of
// this fitness.
func (f Fitness) String() string {
	return fmt.Sprintf("Fitness(dist: %f, over-demand: %f, total: %f)", f.Distance, f.OverDemand, f.Total)
}
