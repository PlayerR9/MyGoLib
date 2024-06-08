package MathExt

import (
	"math"
	"math/big"

	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// IsValidNumber checks if the given number is valid for the given base.
//
// Parameters:
//   - n: The number to check.
//   - base: The base of the number.
//
// Returns:
//   - bool: True if the number is valid for the given base, false otherwise.
func IsValidNumber(n []int, base int) bool {
	if base < 1 {
		return false
	} else if base == 1 {
		return true
	}

	for _, digit := range n {
		if digit < 0 || digit >= base {
			return false
		}
	}

	return true
}

// DecToBase converts a decimal number to a number of the given base.
// The number's Least Significant Digit (LSD) is at index 0.
//
// Parameters:
//   - n: The decimal number to convert.
//   - base: The base of the result number.
//
// Returns:
//   - []int: The number in the given base.
//   - error: An error if the base is invalid.
func DecToBase(n, base int) ([]int, error) {
	if base <= 0 {
		return nil, ers.NewErrInvalidParameter("base", ers.NewErrGT(0))
	}

	if n < 0 {
		n = -n
	}

	// Immediate cases
	if base == 1 {
		return make([]int, n), nil
	} else if n < base {
		return []int{n}, nil
	}

	logBase := math.Log(float64(base))
	result := make([]int, int(math.Log(float64(n))/logBase+1))

	for i := range result {
		result[i] = n % base
		n /= base
	}

	return result, nil
}

// Add adds two numbers of the same base. Both numbers are Least Significant Digit (LSD) first.
//
// Parameters:
//   - n1: The first number to add.
//   - n2: The second number to add.
//   - base: The base of the numbers.
//
// Returns:
//   - []int: The sum of the two numbers.
//   - error: An error if the base is invalid.
func Add(n1, n2 []int, base int) ([]int, error) {
	if base <= 0 {
		return nil, ers.NewErrInvalidParameter("base", ers.NewErrGT(0))
	}

	if base == 1 {
		return make([]int, len(n1)+len(n2)), nil
	}

	maxLen := len(n1)
	if len(n2) > maxLen {
		maxLen = len(n2)
	}

	// Add the two binary numbers.
	result := make([]int, maxLen)
	carry := 0

	// Add the digits for the common length of n1 and n2
	for i := 0; i < len(n1) && i < len(n2); i++ {
		result[i] = n1[i] + n2[i] + carry
		carry = result[i] / base
		result[i] %= base
	}

	// Add the remaining digits of the longer number
	for i := len(n2); i < len(n1); i++ {
		result[i] = n1[i] + carry
		carry = result[i] / base
		result[i] %= base
	}

	for i := len(n1); i < len(n2); i++ {
		result[i] = n2[i] + carry
		carry = result[i] / base
		result[i] %= base
	}

	if carry > 0 {
		result = append(result, carry)
	}

	return result, nil
}

// Subtract subtracts two numbers of the same base. Both numbers are Least Significant Digit (LSD) first.
//
// Parameters:
//   - n1: The number to subtract from.
//   - n2: The number to subtract.
//   - base: The base of the numbers.
//
// Returns:
//   - []int: The result of the subtraction.
//   - error: An error if the base is invalid or the subtraction resulted in a negative number.
func Subtract(n1, n2 []int, base int) ([]int, error) {
	if base <= 0 {
		return nil, ers.NewErrInvalidParameter("base", ers.NewErrGT(0))
	}

	if base == 1 {
		return make([]int, len(n1)), nil
	}

	// Subtract the two binary numbers.
	result := make([]int, len(n1))
	borrow := 0

	// Subtract the digits for the common length of n1 and n2
	for i := 0; i < len(n1) && i < len(n2); i++ {
		result[i] = n1[i] - n2[i] - borrow

		if result[i] < 0 {
			result[i] += base
			borrow = 1
		} else {
			borrow = 0
		}
	}

	// Subtract the remaining digits of the longer number
	for i := len(n2); i < len(n1); i++ {
		result[i] = n1[i] - borrow

		if result[i] < 0 {
			result[i] += base
			borrow = 1
		} else {
			borrow = 0
		}
	}

	if borrow > 0 {
		return nil, NewErrSubtractionUnderflow()
	}

	// Remove leading zeros
	i := len(result) - 1
	for ; i >= 0 && result[i] == 0; i-- {
	}
	result = result[:i+1]

	if len(result) == 0 {
		result = []int{0}
	}

	return result, nil
}

// BaseToDec converts a number of the given base to a decimal number.
// The number's Least Significant Digit (LSD) is at index 0.
//
// Parameters:
//   - n: The number to convert.
//   - base: The base of the number.
//
// Returns:
//   - int: The decimal number.
//   - error: An error if the base is invalid or the number is invalid for the given base.
func BaseToDec(n []int, base int) (int, error) {
	if base <= 0 {
		return 0, ers.NewErrInvalidParameter("base", ers.NewErrGT(0))
	}

	if base == 1 {
		return len(n), nil
	}

	result := 0

	for i, digit := range n {
		if digit < 0 || digit >= base {
			return 0, ers.NewErrOutOfBounds(digit, 0, base)
		}

		result += digit * int(math.Pow(float64(base), float64(i)))
	}

	return result, nil
}

func IntToBigInt(n int) *big.Int {
	return new(big.Int).SetInt64(int64(n))
}
