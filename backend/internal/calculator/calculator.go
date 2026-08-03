package calculator

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrMissingOperand   = errors.New("missing operand")
	ErrDivisionByZero   = errors.New("division by zero")
	ErrNegativeRoot     = errors.New("cannot take the square root of a negative number")
	ErrZeroToZero       = errors.New("0^0 is undefined")
	ErrZeroNegative     = errors.New("zero cannot be raised to a negative exponent")
	ErrNonFiniteResult  = errors.New("result is outside the supported numeric range")
	ErrZeroBase         = errors.New("percentage base cannot be zero")
)

type Operation string

const (
	Add Operation = "add"
	Subtract Operation = "subtract"
	Multiply Operation = "multiply"
	Divide Operation = "divide"
	Power Operation = "power"
	Sqrt Operation = "sqrt"
	Percent Operation = "percentage"
)

func Calculate(operation Operation, a, b *float64) (float64, error) {
	if a == nil {
		return 0, ErrMissingOperand
	}
	var result float64
	switch operation {
	case Add:
		if b == nil {
			return 0, ErrMissingOperand
		}
		result = *a + *b
	case Subtract:
		if b == nil {
			return 0, ErrMissingOperand
		}
		result = *a - *b
	case Multiply:
		if b == nil {
			return 0, ErrMissingOperand
		}
		result = *a * *b
	case Divide:
		if b == nil {
			return 0, ErrMissingOperand
		}
		if *b == 0 {
			return 0, ErrDivisionByZero
		}
		result = *a / *b
	case Power:
		if b == nil {
			return 0, ErrMissingOperand
		}
		if *a == 0 && *b == 0 {
			return 0, ErrZeroToZero
		}
		if *a == 0 && *b < 0 {
			return 0, ErrZeroNegative
		}
		result = math.Pow(*a, *b)
	case Sqrt:
		if *a < 0 {
			return 0, ErrNegativeRoot
		}
		result = math.Sqrt(*a)
	case Percent:
		if b == nil {
			return 0, ErrMissingOperand
		}
		if *b == 0 {
			return 0, ErrZeroBase
		}
		result = (*a / *b) * 100
	default:
		return 0, fmt.Errorf("%w: %s", ErrUnknownOperation, operation)
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}
