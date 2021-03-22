package entities

type Location interface {
	GetPosition() (X, Y float64)
}
