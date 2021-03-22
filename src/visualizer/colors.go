package visualizer

import (
	"math"

	"gopkg.in/go-playground/colors.v1"
)

func scalarToColor(value float64) (colors.Color, error) {
	var a = (1 - value) / 0.2
	var X = math.Floor(a)
	var Y = uint8(math.Floor(255 * (a - X)))

	var r, g, b uint8

	switch X {
	case 0:
		r = 255
		g = Y
		b = 0
		break
	case 1:
		r = 255 - Y
		g = 255
		b = 0
		break
	case 2:
		r = 0
		g = 255
		b = Y
		break
	case 3:
		r = 0
		g = 255 - Y
		b = 255
		break
	case 4:
		r = Y
		g = 0
		b = 255
		break
	case 5:
		r = 255
		g = 0
		b = 255
		break
	}

	return colors.RGB(r, g, b)
}
