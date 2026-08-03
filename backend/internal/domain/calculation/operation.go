package calculation

import "fmt"

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

func ParseOperation(value string) (Operation, error) {
	operation := Operation(value)
	switch operation {
	case Add, Subtract, Multiply, Divide, Power, Sqrt, Percent:
		return operation, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownOperation, value)
	}
}

func (operation Operation) requiresRightOperand() bool {
	return operation != Sqrt
}
