package MathExt

// PrimeFactorization is a function that performs prime factorization on an
// input number.
// It takes an integer as input and returns a map where the keys are the
// prime factors and the values are their respective powers.
//
// The function first checks if the input number is 0, and if so, it panics
// because 0 cannot be factorized.
// If the input number is 1 or -1, it returns a map with 1 as the only factor.
// If the input number is negative, it converts it to a positive number.
//
// The function then initializes an empty map to store the prime factors and
// their powers, and a variable to keep track of the current prime factor.
// It then enters a loop that continues until the input number is reduced to
// 1.
//
// When the loop ends, the function returns the map of prime factors.
func PrimeFactorization(inputNumber int) map[int]int {
	if inputNumber == 0 {
		panic("Cannot factorize 0")
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
// (GCD) of two integers.
// It takes two integers, a and b, as input and returns the GCD as an integer.
//
// The function first checks if either of the input numbers is 0. If a is 0, it
// returns b. If b is 0, it returns a.
// This is because the GCD of 0 and any number is that number.
//
// The function then ensures that a is always the larger number. If a is less than
// b, it swaps the values of a and b.
//
// The function then uses the Euclidean algorithm to find the GCD. The Euclidean
// algorithm is a method for finding the GCD of two numbers
// by repeatedly replacing the larger number with the remainder of the division
// of the larger number by the smaller number, until the remainder is 0.
// In this function, this is done in a loop that continues until b is 0. In each
// iteration of the loop, a is replaced with b and b is replaced with the remainder
// of a divided by b.
//
// When the loop ends, a is the GCD of the original input numbers, so the function
// returns a.
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
