package visualizer

import (
	"math"
	"sync"
	"time"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"

	"github.com/jorgenhanssen/go-genetic-mdvrp/src/entities"
	"github.com/jorgenhanssen/go-genetic-mdvrp/src/solver"
)

type Instance struct {
	window *sdlcanvas.Window
	canvas *canvas.Canvas

	depots    entities.Depots
	customers entities.Customers
	bestAgent *solver.Agent

	// bounding box for nodes
	// X = max X, x = min x
	X, x, Y, y float64

	mu sync.Mutex

	stop chan bool
}

func New() (*Instance, error) {
	i := Instance{
		stop: make(chan bool),
	}
	wnd, cv, err := sdlcanvas.CreateWindow(1280, 720, "MDVRP")
	if err != nil {
		return nil, err
	}

	i.window = wnd
	i.canvas = cv

	return &i, nil
}

func (i *Instance) Draw(depots entities.Depots, customers entities.Customers, best *solver.Agent) {
	i.mu.Lock()
	i.customers = customers
	i.depots = depots
	i.bestAgent = best
	i.mu.Unlock()

	locations := []entities.Location{}
	for _, item := range customers {
		locations = append(locations, item)
	}
	for _, item := range depots {
		locations = append(locations, item)
	}

	X, x, Y, y := getBoundary(locations)

	i.X = X + 20
	i.x = x - 20
	i.Y = Y + 20
	i.y = y - 20
}

func (i *Instance) Run() {
	i.window.MainLoop(func() {
		if i.bestAgent == nil {
			time.Sleep(time.Millisecond * 500)
			return
		}

		i.mu.Lock()
		w, h := float64(i.canvas.Width()), float64(i.canvas.Height())
		customers := i.customers
		depots := i.depots
		bestAgent := i.bestAgent
		i.mu.Unlock()

		i.canvas.SetFillStyle("#15202e")
		i.canvas.FillRect(0, 0, w, h)

		// Draw customers
		i.canvas.SetStrokeStyle("#FFF5")
		i.canvas.SetLineWidth(2)
		for _, customer := range customers {
			i.canvas.BeginPath()
			// i.canvas.Arc((customer.X-i.x)*(w/(i.X-i.x)), (customer.Y-i.y)*(h/(i.Y-i.y)), 1, 0, math.Pi*2, false)
			x, y := i.scaledPosition(customer.GetPosition())
			i.canvas.Arc(x, y, 1, 0, math.Pi*2, false)
			i.canvas.Stroke()
		}

		// Draw depots
		i.canvas.SetStrokeStyle("#FFF")
		i.canvas.SetLineWidth(4)
		for _, depot := range depots {
			i.canvas.BeginPath()
			// i.canvas.Arc((depot.X-i.x)*(w/(i.X-i.x)), (depot.Y-i.y)*(h/(i.Y-i.y)), 3, 0, math.Pi*2, false)
			x, y := i.scaledPosition(depot.GetPosition())
			i.canvas.Arc(x, y, 1, 0, math.Pi*2, false)
			i.canvas.Stroke()
		}

		// Draw agents
		i.canvas.SetLineWidth(1)
		for routeID, route := range bestAgent.Dna {
			color, err := scalarToColor(float64(routeID) / float64(len(bestAgent.Dna)))
			if err != nil {
				panic(err)
			}
			i.canvas.SetStrokeStyle(color.String())

			depot := depots[route.DepotID]
			prevX, prevY := i.scaledPosition(depot.GetPosition())
			i.canvas.BeginPath()
			i.canvas.MoveTo(prevX, prevY)
			for _, cID := range route.Path {
				customer := customers[cID]
				i.canvas.LineTo(i.scaledPosition(customer.GetPosition()))
			}
			i.canvas.LineTo(i.scaledPosition(depot.GetPosition()))
			i.canvas.ClosePath()
			i.canvas.Stroke()
		}

		select {
		case <-i.stop:
			i.window.Close()
		default:
			return
		}
	})
}

func (i *Instance) Stop() {
	i.stop <- true
}

func (i *Instance) scaledPosition(_x, _y float64) (float64, float64) {
	w, h := float64(i.canvas.Width()), float64(i.canvas.Height())
	return (_x - i.x) * (w / (i.X - i.x)), (_y - i.y) * (h / (i.Y - i.y))
}

func getBoundary(locations []entities.Location) (X float64, x float64, Y float64, y float64) {
	X, Y = locations[0].GetPosition()
	x, y = X, Y

	for _, location := range locations {
		_x, _y := location.GetPosition()
		if X < _x {
			X = _x
		} else if x > _x {
			x = _x
		}
		if Y < _y {
			Y = _y
		} else if y > _y {
			y = _y
		}
	}

	return
}
