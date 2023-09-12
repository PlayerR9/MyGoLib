package MathExt

import "strconv"

// Max returns the maximum of a slice of numbers. Panics if the slice is empty. For example, Max(1, 2, 3, 4) returns 4.
//
// Parameters:
// 	- numbers: The slice of numbers to get the maximum of.
//
// Returns:
// 	- int: The maximum of the slice of numbers.
func Max(numbers ...int) int {
	if len(numbers) == 0 {
		panic("Cannot get max of no numbers")
	}

	max := numbers[0]

	for _, number := range numbers[1:] {
		if number > max {
			max = number
		}
	}

	return max
}

// Min returns the minimum of a slice of numbers. Panics if the slice is empty. For example, Min(1, 2, 3, 4) returns 1.
//
// Parameters:
// 	- numbers: The slice of numbers to get the minimum of.
//
// Returns:
// 	- int: The minimum of the slice of numbers.
func Min(numbers ...int) int {
	if len(numbers) == 0 {
		panic("Cannot get min of no numbers")
	}

	min := numbers[0]

	for _, number := range numbers[1:] {
		if number < min {
			min = number
		}
	}

	return min
}

// Compare returns whether n1 is greater than, equal to, or less than n2. For example, Compare(1, 2) returns -1.
//
// Parameters:
// 	- n1: The first number to compare.
// 	- n2: The second number to compare.
//
// Returns:
// 	- int: > 0 if n1 > n2, 0 if n1 == n2, < 0 if n1 < n2.
func Compare(n1, n2 int) int {
	return n1 - n2
}

func is_in_map(n int, m map[int]int) bool {
	for k := range m {
		if n%k != 0 {
			return true
		}
	}

	return false
}

// PrimeFactorization returns a slice of all the factors of a number (excluding 1 and itself, unless it is a prime).
// Panics if num is 0. For example, PrimeFactorization(12) returns []int{2, 2, 3}. Negative numbers are converted to positive.
//
// Parameters:
// 	- num: The number to factorize.
//
// Returns:
// 	- map[int]int: A map of the factors of the number, with the key being the factor and the value being the number of times it is a factor.
//
// Information:
// 	- The factors are not sorted.
//  - The number 1 and -1 are factorized to map[int]int{1: 1} (the only time 1 appears in the map is as a key with a value of 1)
func PrimeFactorization(num int) map[int]int {
	// Check for 0
	if num == 0 {
		panic("Cannot factorize 0")
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
		for is_in_map(current_factor, factors) {
			current_factor++
		}

		// Find the number of times current_factor is a factor
		factor_count := 0 // The number of times the current factor is a factor

		for num > 1 && (num%current_factor) == 0 {
			factor_count++
			num /= current_factor
		}

		if factor_count != 0 {
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
// 	- n1: The first number.
// 	- n2: The second number.
//
// Returns:
// 	- int: The greatest common divisor of n1 and n2.
func GetGCD(n1, n2 int) int {
	if n1 == 0 {
		panic("Cannot get GCD of 0 and " + strconv.Itoa(n2))
	}

	if n2 == 0 {
		panic("Cannot get GCD of " + strconv.Itoa(n1) + " and 0")
	}

	if Compare(n1, n2) < 0 {
		n1, n2 = n2, n1
	}

	factors := PrimeFactorization(n1)

	gcd := 1

	for i := 0; i < len(factors); i++ {
		if n2%factors[i] == 0 {
			gcd *= factors[i]
		}
	}

	return gcd
}
