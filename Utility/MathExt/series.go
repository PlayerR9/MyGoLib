package MathExt

import (
	"errors"
	"fmt"
	"math/big"

	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Serieser is an interface for series.
type Serieser interface {
	// Term returns the nth term of the series.
	//
	// Parameters:
	//   - n: The term number.
	//
	// Returns:
	//   - *big.Int: The nth term of the series.
	//   - error: An error if the term cannot be calculated.
	Term(n int) (*big.Int, error)
}

// ConvergenceResult is a struct that holds the values of the convergence of a series.
type ConvergenceResult struct {
	// values is a slice of the convergence values.
	values []*big.Float
}

// ApproximateConvergence approximates the convergence of a series.
// It calculates the average of the last n values.
//
// Parameters:
//   - n: The number of values to calculate the average.
//
// Returns:
//   - *big.Float: The average of the last n values.
//   - error: An error if the calculation fails.
func (cr *ConvergenceResult) ApproximateConvergence(n int) (*big.Float, error) {
	if n <= 0 {
		return nil, ue.NewErrInvalidParameter("n", ue.NewErrGT(0))
	} else if len(cr.values) < n {
		return nil, ue.NewErrInvalidParameter("n", errors.New("not enough values to calculate average"))
	}

	sum := new(big.Float)
	for i := len(cr.values) - n; i < len(cr.values); i++ {
		sum.Add(sum, cr.values[i])
	}

	average := new(big.Float).Quo(sum, new(big.Float).SetFloat64(float64(n)))
	return average, nil
}

// CalculateConvergence calculates the convergence of a series.
// It calculates the quotient of the ith term and the (i+delta)th term.
//
// Parameters:
//   - series: The series to calculate the convergence.
//   - upperLimit: The upper limit of the series to calculate the convergence.
//   - delta: The difference between the terms to calculate the convergence.
//
// Returns:
//   - *ConvergenceResult: The convergence result.
//   - error: An error if the calculation fails.
func CalculateConvergence(series Serieser, upperLimit int, delta int) (*ConvergenceResult, error) {
	result := &ConvergenceResult{
		values: make([]*big.Float, 0),
	}

	for i := 0; i < upperLimit-delta; i++ {
		ithTerm, err := series.Term(i)
		if err != nil {
			return result, fmt.Errorf("cannot get %dth term: %s", i, err.Error())
		}

		ithPlusDeltaTerm, err := series.Term(i + delta)
		if err != nil {
			return result, fmt.Errorf("cannot get %dth term: %s", i+delta, err.Error())
		}

		quotient := new(big.Float).Quo(
			new(big.Float).SetInt(ithPlusDeltaTerm),
			new(big.Float).SetInt(ithTerm),
		)

		result.values = append(result.values, quotient)
	}

	return result, nil
}

// LinearRegression is a struct that holds the equation of a linear regression.
type LinearRegression struct {
	// A is the coefficient of the linear regression.
	A *big.Float

	// B is the exponent of the linear regression.
	B *big.Float
}

// String implements the fmt.Stringer interface.
//
// Format: y = a * x^b
func (lr *LinearRegression) String() string {
	return fmt.Sprintf("y = %s * x ^ %s", lr.A.String(), lr.B.String())
}

// NewLinearRegression creates a new LinearRegression.
//
// Returns:
//   - LinearRegression: The new LinearRegression.
func NewLinearRegression() *LinearRegression {
	return &LinearRegression{
		A: new(big.Float).SetPrec(1000),
		B: new(big.Float).SetPrec(1000),
	}
}

// FindEquation is a method of ConvergenceResult that finds the equation of the series
// that best fits the convergence values.
//
// The equation is of the form y = a * x^b.
//
// Returns:
//   - float64: The value of a.
//   - float64: The value of b.
//   - error: An error if the calculation fails.
func (l *LinearRegression) FindEquation(cr *ConvergenceResult) error {
	if cr == nil {
		return ue.NewErrNilParameter("cr")
	}

	if len(cr.values) < 2 {
		return fmt.Errorf("not enough values to calculate equation")
	}

	sumX := new(big.Float)
	sumY := new(big.Float)
	sumXY := new(big.Float)
	sumX2 := new(big.Float)
	n := big.NewFloat(float64(len(cr.values)))

	for i, v := range cr.values {
		x := big.NewFloat(float64(i))
		y := new(big.Float).Set(v)

		sumX.Add(sumX, x)
		sumY.Add(sumY, y)
		sumXY.Add(sumXY, new(big.Float).Mul(x, y))
		sumX2.Add(sumX2, new(big.Float).Mul(x, x))
	}

	// a = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	l.A = new(big.Float).Sub(
		new(big.Float).Mul(n, sumXY),
		new(big.Float).Mul(sumX, sumY),
	)
	l.A = l.A.Quo(l.A, new(big.Float).Sub(
		new(big.Float).Mul(n, sumX2),
		new(big.Float).Mul(sumX, sumX),
	))
	l.A = new(big.Float).SetPrec(1000).Set(l.A)

	// b = (sumY - a*sumX) / n
	l.B = new(big.Float).Sub(
		sumY,
		new(big.Float).Mul(l.A, sumX),
	)
	l.B = l.B.Quo(l.B, n)
	l.B = new(big.Float).SetPrec(1000).Set(l.B)

	return nil
}
