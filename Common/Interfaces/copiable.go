package Interfaces

// Copier is an interface that provides a method to create a deep copy of an object.
type Copier interface {
	// Copy creates a shallow copy of the object.
	//
	// Returns:
	//
	//   - Copier: A shallow copy of the object.
	Copy() Copier
}
