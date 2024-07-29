package Debugging

import (
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
)

const (
	// PointerArrow is the arrow used to indicate the current position of the pointer.
	PointerArrow string = "^"
)

// PrintPointer generates a string representation of a pointer at a given distance.
//
// Parameters:
//   - distance: The distance of the pointer from the current position.
//
// Returns:
//   - string: A string representation of the pointer.
//
// Behaviors:
//   - Negative distance: The pointer is to the left of the current position.
//   - Positive distance: The pointer is to the right of the current position.
//   - 0: The pointer is at the current position.
func PrintPointer(distance int) string {
	var builder strings.Builder

	if distance < 0 {
		builder.WriteString(PointerArrow)
		builder.WriteString(strings.Repeat(" ", -distance-1))
	} else if distance == 0 {
		builder.WriteString(PointerArrow)
	} else {
		builder.WriteString(strings.Repeat(" ", distance-1))
		builder.WriteString(PointerArrow)
	}

	return builder.String()
}

// Backuper is an interface that represents a type that can be backed up and restored.
type Backuper[T any] interface {
	// Backup creates a backup of the object.
	//
	// Returns:
	//   - T: The backup of the object.
	Backup() T

	// Restore restores the object from a backup.
	//
	// Parameters:
	//   - backup: The backup of the object.
	//
	// Returns:
	//   - error: An error if the restoration fails.
	Restore(backup T) error
}

// DoWithBackup executes a function with a backup of the subject object.
//
// Parameters:
//   - subject: The object to create a backup of.
//   - f: The function to execute.
//
// Returns:
//   - error: An error if the function fails.
//
// Behaviors:
//   - If the function fails or does not accept the subject object, the
//     subject object is restored from the backup. Otherwise, the subject
//     object is left as is.
func DoWithBackup[T Backuper[E], E any](subject T, f uc.EvalOneFunc[T, bool]) error {
	backup := subject.Backup()

	accept, err := f(subject)
	if err != nil || !accept {
		subject.Restore(backup)
	}

	return err
}
