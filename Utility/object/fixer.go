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

// Cleaner is an interface for objects that can be cleaned.
type Cleaner interface {
	// Cleanup cleans the object.
	Cleanup()
}

/*
func Cleanup(elem any) any {
	if elem == nil {
		return nil
	}

	switch v := elem.(type) {
	case Cleaner:
		v.Cleanup()
	case []any:
		for i := 0; i < len(v); i++ {
			Cleanup(v[i])

			v[i] = nil
		}

		return v[:0]
	case map[any]any:
		for k := range v {
			Cleanup(v[k])

			delete(v, k)
		}

		return nil
	default:
		zero_val := reflect.Zero(reflect.TypeOf(elem))

		return zero_val.Interface()
	}
	/*
	case *int:
		*v = 0
	case *int8:
		*v = 0
	case *int16:
		*v = 0
	case *int32:
		*v = 0
	case *int64:
		*v = 0
	case *uint:
		*v = 0
	case *uint8:
		*v = 0
	case *uint16:
		*v = 0
	case *uint32:
		*v = 0
	case *uint64:
		*v = 0
	case *float32:
		*v = 0
	case *float64:
		*v = 0
	case *bool:
		*v = false
	case *string:
		*v = ""
	case *[]any:
		*v = nil
	case *map[any]any:
		*v = nil
	case *struct{}:
		*v = struct{}{}
	}

}
*/
