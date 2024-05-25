package Tray

// Trayable is an interface that represents a type that can be converted to a Tray.
type Trayable[T any] interface {
	// ToTray converts the Trayable to a Tray.
	//
	// Returns:
	//   - Trayer[T]: The Tray representation of the Trayable.
	ToTray() Trayer[T]
}

// Trayer is an interface that represents any type of Tray.
type Trayer[T any] interface {
	// Move moves the arrow by n positions.
	//
	// Parameters:
	//   - n: The number of positions to move the arrow.
	//
	// Returns:
	//   - int: The number of positions the arrow's movement was limited by.
	//
	// Behaviors:
	//   - Negative n: Move the arrow to the left by n positions.
	//   - Positive n: Move the arrow to the right by n positions.
	//   - 0: Do not move the arrow.
	Move(n int) int

	// Write writes the given element to the tape at the arrow position.
	//
	// Parameters:
	//   - elem: The element to write.
	//
	// Returns:
	//   - error: An error of type *ers.Empty[*GeneralTray] if the tape is empty.
	Write(elem T) error

	// Read reads the element at the arrow position.
	//
	// Returns:
	//   - T: The element at the arrow position.
	//   - error: An error of type *ers.Empty[*GeneralTray] if the tape is empty.
	Read() (T, error)

	// Delete deletes n elements to the right of the arrow position.
	//
	// Parameters:
	//   - n: The number of elements to delete.
	//
	// Returns:
	//   - error: An error of type *ers.ErrInvalidParameter if n is less than 0.
	Delete(n int) error

	// Insert inserts the given elements to the tape at the arrow position.
	//
	// Parameters:
	//   - elems: The elements to insert.
	Insert(elems ...T)

	// ExtendTapeOnLeft extends the tape on the left with the given elements.
	//
	// Parameters:
	//   - elems: The elements to add.
	ExtendTapeOnLeft(elems ...T)

	// ExtendTapeOnRight extends the tape on the right with the given elements.
	//
	// Parameters:
	//   - elems: The elements to add.
	ExtendTapeOnRight(elems ...T)

	// ArrowStart moves the arrow to the start of the tape.
	ArrowStart()

	// ArrowEnd moves the arrow to the end of the tape.
	ArrowEnd()

	/*

		// ShiftLeftOfArrow shifts the elements on the right of the
		// arrow to the left by n positions.
		//
		// Parameters:
		//   - n: The number of positions to shift the elements.
		ShiftLeftOfArrow(n int)

		// ShiftRightOfArrow shifts the elements on the left of the
		// arrow to the right by n positions.
		//
		// Parameters:
		//   - n: The number of positions to shift the elements.
		ShiftRightOfArrow(n int)

	*/
}
