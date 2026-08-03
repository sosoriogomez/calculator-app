package calculation

import (
	"context"
	"errors"
	"testing"
)

func floatPtr(value float64) *float64 { return &value }

func TestServiceCalculates(t *testing.T) {
	service := NewService()
	tests := []struct { operation string; a, b *float64; want float64 }{
		{"add", floatPtr(2), floatPtr(3), 5}, {"subtract", floatPtr(7), floatPtr(4), 3},
		{"multiply", floatPtr(2.5), floatPtr(4), 10}, {"divide", floatPtr(9), floatPtr(2), 4.5},
		{"power", floatPtr(2), floatPtr(3), 8}, {"sqrt", floatPtr(81), nil, 9}, {"percentage", floatPtr(25), floatPtr(200), 12.5},
	}
	for _, test := range tests {
		result, err := service.Calculate(context.Background(), Command{Operation: test.operation, A: *test.a, B: test.b})
		if err != nil { t.Errorf("%s returned error %v", test.operation, err) } else if result.Value != test.want { t.Errorf("%s = %v, want %v", test.operation, result.Value, test.want) }
	}
}

func TestServiceClassifiesInvalidCommands(t *testing.T) {
	service := NewService()
	_, err := service.Calculate(context.Background(), Command{Operation: "modulo", A: 1})
	if !errors.Is(err, ErrInvalidCommand) { t.Errorf("error = %v", err) }
	_, err = service.Calculate(context.Background(), Command{Operation: "add", A: 1})
	if !errors.Is(err, ErrInvalidCommand) { t.Errorf("missing operand error = %v", err) }
	_, err = service.Calculate(context.Background(), Command{Operation: "add", A: 1, B: floatPtr(0)})
	if err != nil { t.Errorf("valid command error = %v", err) }
}

func TestServicePropagatesDomainErrorsAndContext(t *testing.T) {
	service := NewService()
	_, err := service.Calculate(context.Background(), Command{Operation: "divide", A: 1, B: floatPtr(0)})
	if err == nil || err.Error() != "division by zero" { t.Errorf("domain error = %v", err) }
	ctx, cancel := context.WithCancel(context.Background()); cancel()
	_, err = service.Calculate(ctx, Command{Operation: "add", A: 1, B: floatPtr(2)})
	if !errors.Is(err, context.Canceled) { t.Errorf("context error = %v", err) }
}
