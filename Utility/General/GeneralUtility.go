package Utility

import (
	"fmt"
	"os"
)

// PressEnterToContinue prints "Press enter to continue..." to the console and waits for the user to press enter.
func PressEnterToContinue() {
	fmt.Println("Press enter to continue...")
	fmt.Scanln()
}

// PressEnterToExit prints "Press enter to exit..." to the console and waits for the user to press enter. Then, it exits the program with the given exit code.
//
// Parameters:
//   - exit_code: The exit code to exit the program with.
func PressEnterToExit(exit_code int) {
	fmt.Println("Press enter to exit...")
	fmt.Scanln()

	os.Exit(exit_code)
}

// BinarySearch searches for an element in a sorted slice of integers using the binary search algorithm.
//
// Parameters:
//   - elements: The slice of integers to search in.
//   - e: The element to search for.
//
// Returns:
//   - int: The index of the element in the slice, or -1 if the element is not in the slice.
func BinarySearch(elements []int, e int) int {
	sx := 0
	dx := len(elements) - 1
	pos := -1

	for sx < dx && pos == -1 {
		m := (sx + dx) / 2

		if elements[m] == e {
			pos = m
		} else if elements[m] < e {
			sx = m + 1
		} else {
			dx = m
		}
	}

	return pos
}
