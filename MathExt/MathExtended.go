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

	for _, number := range numbers {
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

	for _, number := range numbers {
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

// PrimeFactorization returns a slice of all the factors of a number (including 1) in ascending order.
// Panics if num is 0. For example, PrimeFactorization(12) returns []int{1, 2, 2, 3}. Negative numbers are converted to positive.
//
// Parameters:
// 	- num: The number to factorize.
//
// Returns:
// 	- []int: A slice of all the factors of num in ascending order.
func PrimeFactorization(num int) []int {
	if num == 0 {
		panic("Cannot factorize 0")
	}

	if num < 0 {
		num = -num
	}

	factors := []int{1}

	factor := 2

	for num >= 2 && factor < num/2 {
		if num%factor == 0 {
			factors = append(factors, factor)
			num /= factor
		} else {
			factor++
		}
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
