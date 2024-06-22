package object

// Fixer is an interface for objects that can be fixed.
type Fixer interface {
	// Fix fixes the object.
	//
	// Returns:
	//   - error: An error if the object could not be fixed.
	Fix() error
}

// Fix fixes an object if it exists.
//
// Parameters:
//   - elem: The object to fix.
//   - mustExists: A flag indicating if the object must exist.
//
// Returns:
//   - error: An error if the object could not be fixed.
//
// Behaviors:
//   - Returns an error if the object must exist but does not.
//   - Returns nil if the object is nil and does not need to exist.
func Fix(elem Fixer, mustExists bool) error {
	if elem != nil {
		err := elem.Fix()
		if err != nil {
			return err
		}
	} else if mustExists {
		return NewErrValueMustExists()
	}

	return nil
}
