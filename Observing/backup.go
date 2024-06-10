package Observing

// Backuper is an interface that provides methods to create a backup of an element
// and restore the element from the backup.
type Backuper[T any] interface {
	// Backup creates a backup of the element.
	//
	// Returns:
	//   - Backuper: The backup of the element.
	Backup() T

	// Restore restores the element from the backup.
	//
	// Parameters:
	//   - backup: The backup to restore from.
	Restore(backup T)
}

// DoWithBackup calls the given function with the element and creates a backup of the
// element before calling the function. If the function returns an error, the element is
// restored from the backup. If the function returns false, the element is restored from
// the backup.
//
// Parameters:
//   - elem: The element to call the function with.
//   - fn: The function to call with the element.
//
// Returns:
//   - error: An error if the function returns an error.
func DoWithBackup[T Backuper[E], E any](elem T, fn func(T) (bool, error)) error {
	if fn == nil {
		return nil // Nothing to do
	}

	backup := elem.Backup()

	accept, err := fn(elem)
	if err != nil {
		elem.Restore(backup)

		return err
	}

	if !accept {
		elem.Restore(backup)
	}

	return nil
}
