package calculation

import (
	"errors"
	"math"
	"testing"
)

func number(t *testing.T, value float64) Number {
	t.Helper()
	result, err := NewNumber(value)
	if err != nil { t.Fatal(err) }
	return result
}

func numberPtr(t *testing.T, value float64) *Number {
	result := number(t, value)
	return &result
}

func TestParseOperation(t *testing.T) {
	for _, operation := range []Operation{Add, Subtract, Multiply, Divide, Power, Sqrt, Percent} {
		got, err := ParseOperation(string(operation))
		if err != nil || got != operation { t.Errorf("ParseOperation(%q) = %q, %v", operation, got, err) }
	}
	if _, err := ParseOperation("modulo"); !errors.Is(err, ErrUnknownOperation) { t.Errorf("unknown operation error = %v", err) }
}

func TestNewNumberRejectsNonFiniteValues(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NewNumber(value); !errors.Is(err, ErrNonFiniteNumber) { t.Errorf("NewNumber(%v) error = %v", value, err) }
	}
}

func TestCalculationExecutesOperations(t *testing.T) {
	tests := []struct { name string; operation Operation; left, right float64; hasRight bool; want float64 }{
		{"addition", Add, 2, 3, true, 5}, {"subtraction", Subtract, 7, 4, true, 3},
		{"multiplication", Multiply, 2.5, 4, true, 10}, {"division", Divide, 9, 2, true, 4.5},
		{"power", Power, 2, 3, true, 8}, {"negative power", Power, 2, -2, true, .25},
		{"square root", Sqrt, 81, 0, false, 9}, {"percentage", Percent, 25, 200, true, 12.5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var right *Number
			if test.hasRight { right = numberPtr(t, test.right) }
			calculation, err := NewCalculation(test.operation, number(t, test.left), right)
			if err != nil { t.Fatal(err) }
			result, err := calculation.Execute()
			if err != nil { t.Fatal(err) }
			if result.Value() != test.want { t.Errorf("result = %v, want %v", result.Value(), test.want) }
		})
	}
}

func TestCalculationEnforcesInvariants(t *testing.T) {
	tests := []struct { name string; operation Operation; left, right float64; hasRight bool; want error }{
		{"missing right operand", Add, 1, 0, false, ErrMissingOperand}, {"divide by zero", Divide, 1, 0, true, ErrDivisionByZero},
		{"negative root", Sqrt, -1, 0, false, ErrNegativeRoot}, {"zero to zero", Power, 0, 0, true, ErrZeroToZero},
		{"zero to negative", Power, 0, -1, true, ErrZeroNegative}, {"percentage zero base", Percent, 2, 0, true, ErrZeroBase},
		{"unexpected square root operand", Sqrt, 4, 2, true, ErrUnexpectedOperand}, {"unknown operation", Operation("modulo"), 4, 2, true, ErrUnknownOperation},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var right *Number
			if test.hasRight { right = numberPtr(t, test.right) }
			_, err := NewCalculation(test.operation, number(t, test.left), right)
			if !errors.Is(err, test.want) { t.Errorf("error = %v, want %v", err, test.want) }
		})
	}
}

func TestCalculationRejectsNonFiniteResult(t *testing.T) {
	left := number(t, math.MaxFloat64)
	right := numberPtr(t, 2)
	calculation, err := NewCalculation(Multiply, left, right)
	if err != nil { t.Fatal(err) }
	if _, err := calculation.Execute(); !errors.Is(err, ErrNonFiniteResult) { t.Errorf("error = %v", err) }
}
