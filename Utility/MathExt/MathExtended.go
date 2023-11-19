package MathExt

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
