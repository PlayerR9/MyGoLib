package MathExt

import (
	"math/big"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// PrimeFactorization is a function that performs prime factorization on an
// input number.
//
// Parameters:
//   - inputNumber: The number to factorize.
//
// Returns:
//   - map[int]int: A map where the keys are the prime factors and the values
//     are their respective powers. Nil if the input number is 0.
//
// Behaviors:
//   - The input number is converted to a positive number.
//   - The prime factors are sorted in ascending order.
//   - -1 and 1 are represented as [1: 1].
//   - The resulting map does not contain any prime factor with a value of 1.
func PrimeFactorization(inputNumber int) map[int]int {
	if inputNumber == 0 {
		return nil
	}

	if inputNumber == 1 || inputNumber == -1 {
		return map[int]int{1: 1}
	}

	if inputNumber < 0 {
		inputNumber = -inputNumber
	}

	primeFactors := make(map[int]int)
	currentPrimeFactor := 2

	for inputNumber > 1 {
		// Find the next factor such that it is prime
		for {
			isFactorFound := false

			for factor := range primeFactors {
				if currentPrimeFactor%factor != 0 {
					continue
				}

				isFactorFound = true
				break
			}

			if !isFactorFound {
				break
			}

			currentPrimeFactor++
		}

		factorCount := 0

		for inputNumber > 1 && (inputNumber%currentPrimeFactor) == 0 {
			factorCount++
			inputNumber /= currentPrimeFactor
		}

		if factorCount != 0 {
			primeFactors[currentPrimeFactor] = factorCount
		}

		currentPrimeFactor++
	}

	return primeFactors
}

// GreatestCommonDivisor is a function that calculates the greatest common divisor
// (GCD) of two integers using the Euclidean algorithm.
//
// Parameters:
//   - a, b: The two integers to find the GCD of.
//
// Returns:
//   - int: The GCD of the two input numbers.
func GreatestCommonDivisor(a, b int) int {
	// If one of the numbers is 0, return the other number
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	// Ensure that 'a' is always the larger number
	if a < b {
		a, b = b, a
	}

	// Use Euclidean algorithm to find GCD
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

// BigFloatDivision is a function that performs division on two integers and returns
// the result as a big.Float.
//
// Parameters:
//   - numerator: The numerator of the division.
//   - denominator: The denominator of the division.
//
// Returns:
//   - *big.Float: The result of the division as a big.Float.
//   - error: An error of type *ErrInvalidParameter if the numerator is negative or
//     the denominator is less than or equal to 0.
func BigFloatDivision(numerator, denominator int) (*big.Float, error) {
	if numerator < 0 {
		return new(big.Float), uc.NewErrInvalidParameter(
			"numerator",
			uc.NewErrGTE(0),
		)
	}

	if denominator <= 0 {
		return new(big.Float), uc.NewErrInvalidParameter(
			"denominator",
			uc.NewErrGT(0),
		)
	}

	// return R(n1) / R(n2)
	convergence := new(big.Float).Quo(
		new(big.Float).SetInt(IntToBigInt(numerator)),
		new(big.Float).SetInt(IntToBigInt(denominator)),
	)

	return convergence, nil
}
