package entities

// Location is an interface for entities that have
// a readable location.
type Location interface {
	GetPosition() (X, Y float64)
}
