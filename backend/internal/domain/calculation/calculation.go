package calculation

import "math"

// Calculation is the aggregate representing one requested operation and its
// operands. Its constructor enforces invariants before execution is allowed.
type Calculation struct {
	operation Operation
	left Number
	right *Number
}

func NewCalculation(operation Operation, left Number, right *Number) (Calculation, error) {
	switch operation {
	case Add, Subtract, Multiply, Divide, Power, Sqrt, Percent:
	default:
		return Calculation{}, ErrUnknownOperation
	}
	if operation.requiresRightOperand() && right == nil {
		return Calculation{}, ErrMissingOperand
	}
	if operation == Sqrt && right != nil {
		return Calculation{}, ErrUnexpectedOperand
	}
	if operation == Divide && right.Value() == 0 {
		return Calculation{}, ErrDivisionByZero
	}
	if operation == Percent && right.Value() == 0 {
		return Calculation{}, ErrZeroBase
	}
	if operation == Sqrt && left.Value() < 0 {
		return Calculation{}, ErrNegativeRoot
	}
	if operation == Power {
		if left.Value() == 0 && right.Value() == 0 {
			return Calculation{}, ErrZeroToZero
		}
		if left.Value() == 0 && right.Value() < 0 {
			return Calculation{}, ErrZeroNegative
		}
	}
	return Calculation{operation: operation, left: left, right: right}, nil
}

func (calculation Calculation) Execute() (Number, error) {
	left := calculation.left.Value()
	var result float64
	if calculation.right != nil {
		right := calculation.right.Value()
		switch calculation.operation {
		case Add:
			result = left + right
		case Subtract:
			result = left - right
		case Multiply:
			result = left * right
		case Divide:
			result = left / right
		case Power:
			result = math.Pow(left, right)
		case Percent:
			result = (left / right) * 100
		}
	} else {
		result = math.Sqrt(left)
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return Number{}, ErrNonFiniteResult
	}
	return NewNumber(result)
}
