// git tag v0.1.18

package General

import (
	"fmt"
	"os"
)

func WaitForUserConfirmation() {
	fmt.Println("Press enter to proceed...")
	fmt.Scanln()
}

func WaitForExitConfirmation(exitCode int) {
	fmt.Println("Press enter to confirm exit...")
	fmt.Scanln()

	os.Exit(exitCode)
}

func FindMaximumValue[T any](comparisonFunc func(T, T) int, inputValues ...T) T {
	if len(inputValues) == 0 {
		panic("Cannot find maximum value in an empty set")
	}

	maxValue := inputValues[0]
	for _, currentValue := range inputValues[1:] {
		if comparisonFunc(currentValue, maxValue) > 0 {
			maxValue = currentValue
		}
	}

	return maxValue
}

func FindMinimumValue[T any](comparisonFunc func(T, T) int, inputValues ...T) T {
	if len(inputValues) == 0 {
		panic("Cannot find minimum value in an empty set")
	}

	minValue := inputValues[0]
	for _, currentValue := range inputValues[1:] {
		if comparisonFunc(currentValue, minValue) < 0 {
			minValue = currentValue
		}
	}

	return minValue
}
