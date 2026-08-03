package calculation

import "math"

// Number is the calculator's finite-number value object. Keeping the raw
// float private prevents invalid values from entering the domain model.
type Number struct {
	value float64
}

func NewNumber(value float64) (Number, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return Number{}, ErrNonFiniteNumber
	}
	return Number{value: value}, nil
}

func (number Number) Value() float64 {
	return number.value
}
