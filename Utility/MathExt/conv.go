package MathExt

import (
	"fmt"
	"math"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

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

// LSD result[0]
func DecToBase(n, base int) ([]int, error) {
	if base < 1 {
		return nil, ers.NewErrInvalidParameter("base").
			Wrap(NewErrInvalidBase())
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

// Both n1 and n2 must be LSD first and of the same base.
func Add(n1, n2 []int, base int) ([]int, error) {
	if base < 1 {
		return nil, ers.NewErrInvalidParameter("base").
			Wrap(NewErrInvalidBase())
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

// LSD n[0]
func Subtract(n1, n2 []int, base int) ([]int, error) {
	if base < 1 {
		return nil, ers.NewErrInvalidParameter("base").
			Wrap(NewErrInvalidBase())
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
		return nil, fmt.Errorf("subtraction resulted in a negative number")
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

// LSD n[0]
func BaseToDec(n []int, base int) (int, error) {
	if base < 1 {
		return 0, ers.NewErrInvalidParameter("base").
			Wrap(NewErrInvalidBase())
	}

	if base == 1 {
		return len(n), nil
	}

	result := 0

	for i, digit := range n {
		if digit < 0 || digit >= base {
			return 0, fmt.Errorf("invalid digit at %d: %v", i, ers.NewErrOutOfBound(digit, 0, base))
		}

		result += digit * int(math.Pow(float64(base), float64(i)))
	}

	return result, nil
}
