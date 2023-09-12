package MathExt

import (
	"log"
	"os"
	"strconv"
)

// GLOBAL VARIABLES
var (
	// DebugMode is a boolean that is used to enable or disable debug mode. When debug mode is enabled, the package will print debug messages.
	// **Note:** Debug mode is disabled by default.
	DebugMode bool = false

	debugger *log.Logger = log.New(os.Stdout, "[MathExt] ", log.LstdFlags) // Debugger
)

// PrimeFactorization returns a slice of all the factors of a number (excluding 1 and itself, unless it is a prime).
// Panics if num is 0. For example, PrimeFactorization(12) returns []int{2, 2, 3}. Negative numbers are converted to positive.
//
// Parameters:
//   - num: The number to factorize.
//
// Returns:
//   - map[int]int: A map of the factors of the number, with the key being the factor and the value being the number of times it is a factor.
//
// Information:
//   - The factors are not sorted.
//   - The number 1 and -1 are factorized to map[int]int{1: 1} (the only time 1 appears in the map is as a key with a value of 1)
func PrimeFactorization(num int) map[int]int {
	if num == 0 {
		// Cannot factorize 0, so panic.
		if DebugMode {
			debugger.Panic("cannot factorize 0")
		} else {
			panic("Cannot factorize 0")
		}
	}

	// isInMap returns true if n is not a factor of any of the keys in m, otherwise it returns false.
	isInMap := func(n int, m map[int]int) bool {
		for k := range m {
			if n%k != 0 {
				return true
			}
		}

		return false
	}

	// Check for 1 and -1
	if num == 1 || num == -1 {
		return map[int]int{1: 1}
	}

	// Convert to positive if negative
	if num < 0 {
		num = -num
	}

	// Factorize
	factors := make(map[int]int) // The map of factors
	current_factor := 2          // The current factor

	for num > 1 {
		// Skip current factors that are not prime
		for isInMap(current_factor, factors) {
			current_factor++
		}

		// Find the number of times current_factor is a factor
		factor_count := 0 // The number of times the current factor is a factor

		for num > 1 && (num%current_factor) == 0 {
			factor_count++
			num /= current_factor
		}

		if factor_count != 0 {
			// Add the factor to the map if it is a factor
			factors[current_factor] = factor_count
		}

		current_factor++
	}

	return factors
}

// GetGCD returns the greatest common divisor of two numbers. Panics if either number is 0. For example, GetGCD(12, 18) returns 6.
// Negative numbers are converted to positive.
//
// Parameters:
//   - n1: The first number.
//   - n2: The second number.
//
// Returns:
//   - int: The greatest common divisor of n1 and n2.
func GetGCD(n1, n2 int) int {
	if n1 == 0 {
		panic("Cannot get GCD of 0 and " + strconv.Itoa(n2))
	}

	if n2 == 0 {
		panic("Cannot get GCD of " + strconv.Itoa(n1) + " and 0")
	}

	// Invert if n1 < n2. This is done to reduce the number of iterations.
	if n1 < n2 {
		n1, n2 = n2, n1
	}

	factors := PrimeFactorization(n1) // The factors of n1

	gcd := 1 // The greatest common divisor

	// For each factor of n1, check if it is a factor of n2.
	for factor, count := range factors {
		for n2%factor == 0 && count > 0 {
			// If it is a factor of n2, multiply it to the gcd until it is not a factor of n2.
			gcd *= factor
			count--
		}
	}

	return gcd
}
