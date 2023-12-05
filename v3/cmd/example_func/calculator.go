package example_func

import (
	"context"
	"errors"
)

type AddIntInput struct {
	Num1, Num2 int
}
type AddIntOutput struct {
	Result int
}

// AddInt
// input : AddIntInput
// output : AddIntOutput, AddIntInput
func AddInt(ctx context.Context, input any) (any, error) {
	if aii, ok := input.(AddIntInput); ok {
		r := aii.Num1 + aii.Num2
		return []any{AddIntOutput{Result: r}, aii}, nil
	}
	return nil, errors.New("AddInt - invalid input")
}

type MultiplyIntInput struct {
	Num1, Num2 int
}
type MultiplyIntOutput struct {
	Result int
}

// MultiplyInt
// input : MultiplyIntInput
// output : MultiplyIntOutput, MultiplyIntInput
func MultiplyInt(ctx context.Context, input any) (any, error) {
	if mii, ok := input.(MultiplyIntInput); ok {
		r := mii.Num1 * mii.Num2
		return []any{MultiplyIntOutput{Result: r}, mii}, nil
	}
	return nil, errors.New("MultiplyInt - invalid input")
}
