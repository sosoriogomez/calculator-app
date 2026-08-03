package calculation

import "errors"

var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrMissingOperand = errors.New("missing operand")
	ErrUnexpectedOperand = errors.New("unexpected operand")
	ErrDivisionByZero = errors.New("division by zero")
	ErrNegativeRoot = errors.New("cannot take the square root of a negative number")
	ErrZeroToZero = errors.New("0^0 is undefined")
	ErrZeroNegative = errors.New("zero cannot be raised to a negative exponent")
	ErrNonFiniteNumber = errors.New("number must be finite")
	ErrNonFiniteResult = errors.New("result is outside the supported numeric range")
	ErrZeroBase = errors.New("percentage base cannot be zero")
)
