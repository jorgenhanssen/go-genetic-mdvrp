package solver

import (
	"math"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
)

// Selector is the parent-selection method.
type Selector string

const (
	Roulette Selector = "Roulette"
	Random   Selector = "Random"
)

// GeneratioInfo contains information about the current generation.
type GenerationInfo struct {
	BestAgent         *Agent
	GenerationNumber  int
	PopulationFitness Fitness
}

// distance calculates the distance between two entity location
func distance(a, b entities.Location) float64 {
	ax, ay := a.GetPosition()
	bx, by := b.GetPosition()
	return math.Sqrt(math.Pow(ax-bx, 2) + math.Pow(ay-by, 2))
}
