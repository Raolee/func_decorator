package main

import (
	"context"
	"errors"
	"fmt"
	v3 "func_decorator/v3"
	"func_decorator/v3/cmd/example_func"
)

func main() {
	// FunctionRegistry 및 함수 등록
	registry := v3.NewSimpleFunctionRegistry()
	registry.RegisterFunction("add", example_func.AddInt)
	registry.RegisterFunction("multiply", example_func.MultiplyInt)

	// TaskBuilder 로 func 조립
	addFn, _ := registry.GetFunction("add")
	multiplyFn, _ := registry.GetFunction("multiply")

	tb := v3.NewTaskBuilder(v3.Composite)
	task := tb.AddFunction(addFn).
		AttachConverter(
			func(ctx context.Context, a any) (any, error) {
				errMsg := "converter is invalid input for AddInt output, MultiplyInt input"
				inputs, ok := a.([]any)
				if !ok {
					return nil, errors.New(errMsg)
				}
				aio, ok := inputs[0].(example_func.AddIntOutput)
				if !ok {
					return nil, errors.New(errMsg)
				}
				aii, ok := inputs[1].(example_func.AddIntInput)
				if !ok {
					return nil, errors.New(errMsg)
				}

				return example_func.MultiplyIntInput{Num1: aii.Num1, Num2: aio.Result}, nil

			}).
		AddLastFunction(multiplyFn).
		Build()

	// Task 실행
	ctx := context.Background()
	result, err := task.Execute(ctx, example_func.AddIntInput{
		Num1: 10,
		Num2: 20,
	}) // 예시 입력
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		return
	}

	// AddInt = 10 + 20 = 30
	// MultiplyInt = 10 * 30 = 300
	fmt.Println("Result: ", result)
}
