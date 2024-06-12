package Tray

// TrayBackup is a struct that represents a tape.
type TrayBackup[T any] struct {
	// tape is a slice of elements on the tape.
	tape []T

	// arrow is the position of the arrow on the tape.
	arrow int
}
