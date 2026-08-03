package calculation

import (
	"context"
	"errors"
	"fmt"

	domain "calculator-app/backend/internal/domain/calculation"
)

var ErrInvalidCommand = errors.New("invalid calculation command")

type Command struct {
	Operation string
	A float64
	B *float64
}

type Result struct {
	Operation string
	Value float64
}

type Service struct{}

func NewService() *Service { return &Service{} }

func (service *Service) Calculate(ctx context.Context, command Command) (Result, error) {
	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
	}
	operation, err := domain.ParseOperation(command.Operation)
	if err != nil { return Result{}, invalidCommand(err) }
	left, err := domain.NewNumber(command.A)
	if err != nil { return Result{}, invalidCommand(err) }
	var right *domain.Number
	if command.B != nil {
		value, numberErr := domain.NewNumber(*command.B)
		if numberErr != nil { return Result{}, invalidCommand(numberErr) }
		right = &value
	}
	calculation, err := domain.NewCalculation(operation, left, right)
	if err != nil {
		if errors.Is(err, domain.ErrMissingOperand) || errors.Is(err, domain.ErrUnexpectedOperand) { return Result{}, invalidCommand(err) }
		return Result{}, err
	}
	value, err := calculation.Execute()
	if err != nil { return Result{}, err }
	return Result{Operation: command.Operation, Value: value.Value()}, nil
}

func invalidCommand(err error) error { return fmt.Errorf("%w: %v", ErrInvalidCommand, err) }
