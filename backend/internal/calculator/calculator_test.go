package calculator

import (
	"errors"
	"math"
	"testing"
)

func floatPtr(value float64) *float64 { return &value }

func TestCalculateOperations(t *testing.T) {
	tests := []struct { name string; op Operation; a, b *float64; want float64 }{
		{"addition", Add, floatPtr(2), floatPtr(3), 5}, {"subtraction", Subtract, floatPtr(7), floatPtr(4), 3},
		{"multiplication", Multiply, floatPtr(2.5), floatPtr(4), 10}, {"division", Divide, floatPtr(9), floatPtr(2), 4.5},
		{"power", Power, floatPtr(2), floatPtr(3), 8}, {"negative power", Power, floatPtr(2), floatPtr(-2), .25},
		{"square root", Sqrt, floatPtr(81), nil, 9}, {"percentage", Percent, floatPtr(25), floatPtr(200), 12.5},
	}
	for _, test := range tests { t.Run(test.name, func(t *testing.T) { got, err := Calculate(test.op, test.a, test.b); if err != nil { t.Fatal(err) }; if got != test.want { t.Errorf("got %v, want %v", got, test.want) } }) }
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct { name string; op Operation; a, b *float64; want error }{
		{"missing a", Add, nil, floatPtr(1), ErrMissingOperand}, {"missing b", Add, floatPtr(1), nil, ErrMissingOperand},
		{"divide by zero", Divide, floatPtr(1), floatPtr(0), ErrDivisionByZero}, {"negative root", Sqrt, floatPtr(-1), nil, ErrNegativeRoot},
		{"zero to zero", Power, floatPtr(0), floatPtr(0), ErrZeroToZero}, {"zero to negative", Power, floatPtr(0), floatPtr(-1), ErrZeroNegative},
		{"percentage zero base", Percent, floatPtr(2), floatPtr(0), ErrZeroBase}, {"unknown operation", Operation("modulo"), floatPtr(2), floatPtr(3), ErrUnknownOperation},
		{"overflow", Power, floatPtr(10), floatPtr(309), ErrNonFiniteResult},
	}
	for _, test := range tests { t.Run(test.name, func(t *testing.T) { _, err := Calculate(test.op, test.a, test.b); if !errors.Is(err, test.want) { t.Errorf("error = %v, want %v", err, test.want) } }) }
}

func TestCalculateRejectsNonFiniteInputs(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} { if _, err := Calculate(Add, &value, floatPtr(1)); err == nil { t.Errorf("Calculate(%v) expected error", value) } }
}
