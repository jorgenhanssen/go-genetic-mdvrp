package main

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
)

// LoadProblem reads and loads depots and customers related
// to a problem found in a file in the specified filePath.
func LoadProblem(filePath string) (depots entities.Depots, customers entities.Customers, err error) {
	depots = make(entities.Depots)
	customers = make(entities.Customers)

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	maxNumVehiclesPerDepot := math.MaxInt32
	numCustomers := math.MaxInt32
	numDepots := math.MaxInt32

	scanner := bufio.NewScanner(file)
	for line := 0; scanner.Scan(); line++ {
		text := scanner.Text()

		// This is read last
		if line > numDepots+numCustomers {
			depotIndex := line - (numDepots + numCustomers) - 1
			depot := depots[depotIndex]

			dummy := 0
			if _, err = fmt.Sscanf(text, "%d %f %f", &dummy, &depot.X, &depot.Y); err != nil {
				return
			}
			continue
		}

		if line > numDepots {
			customer := entities.Customer{}
			if _, err = fmt.Sscanf(text, "%d %f %f %f %f",
				&customer.ID,
				&customer.X,
				&customer.Y,
				&customer.ServiceDuration,
				&customer.Demand,
			); err != nil {
				return
			}
			customers[customer.ID] = &customer
			continue
		}

		if line > 0 {
			depot := &entities.Depot{
				MaxNumVehicles: maxNumVehiclesPerDepot,
			}
			if _, err = fmt.Sscanf(text, "%f %f", &depot.MaxRouteDuration, &depot.MaxVehicleLoad); err != nil {
				return
			}
			depots[len(depots)] = depot
			continue
		}

		// This is read first
		if _, err = fmt.Sscanf(text, "%d %d %d", &maxNumVehiclesPerDepot, &numCustomers, &numDepots); err != nil {
			return
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return depots, customers, nil
}
