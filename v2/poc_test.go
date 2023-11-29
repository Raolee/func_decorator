package v2

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
)

// Generic Func 을 reflect로 실행해도 문제가 없는지 검증
func TestExecuteGenericFuncUsingReflect(t *testing.T) {
	builder := NewDecoratedFunctionBuilder[string, string]()

	var fn = func(ctx context.Context, req string) (string, error) {
		fmt.Println(req)
		return req + "/processed", nil
	}
	var reqDecorator = func(ctx context.Context, req string) (string, error) {
		return req + "/req_processed", nil
	}
	var resDecorator = func(ctx context.Context, res string) (string, error) {
		return res + "/res_processed", nil
	}
	var exDecorator = func(ctx context.Context, req string, err error) error {
		return errors.Join(err, errors.New("custom raol error"))
	}

	builder.Func(fn)
	builder.RequestDecorators(reqDecorator, reqDecorator)
	builder.ResponseDecorators(resDecorator, resDecorator)
	builder.ExceptionDecorators(exDecorator)
	builder.PanicHandling(true)
	function := builder.Build()

	reflectFuncCall(function.Call, context.Background(), "test")
}

func TestExecuteGenericFuncToAnotherGenericFuncUsingReflect(t *testing.T) {
	fromFuncBuilder := NewDecoratedFunctionBuilder[int, int]()
	var fromFn = func(ctx context.Context, req int) (int, error) {
		fmt.Println("FromFunc ", req)
		return req, nil
	}
	var fromReqDecorator = func(ctx context.Context, req int) (int, error) {
		fmt.Println("FromFunc Req Decorator ", req)
		return req + 1, nil
	}
	var fromResDecorator = func(ctx context.Context, res int) (int, error) {
		fmt.Println("FromFunc Res Decorator ", res)
		return res + 1, nil
	}
	var fromExDecorator = func(ctx context.Context, req int, err error) error {
		return errors.Join(err, errors.New("custom raol error"))
	}
	fromFuncBuilder.Func(fromFn)
	fromFuncBuilder.RequestDecorators(fromReqDecorator, fromReqDecorator)
	fromFuncBuilder.ResponseDecorators(fromResDecorator, fromResDecorator)
	fromFuncBuilder.ExceptionDecorators(fromExDecorator)
	fromFuncBuilder.PanicHandling(true)
	fromFunction := fromFuncBuilder.Build()

	converter := func(ctx context.Context, req int) (nextReq string, err error) {
		fmt.Println("Adapt : ", req)
		return strconv.Itoa(req), nil
	}

	toFuncBuilder := NewDecoratedFunctionBuilder[string, string]()
	var toFn = func(ctx context.Context, req string) (string, error) {
		fmt.Println("ToFunc : ", req)
		return req + "/processed", nil
	}
	var toReqDecorator = func(ctx context.Context, req string) (string, error) {
		fmt.Println("ToFunc Req Decorator : ", req)
		return req + "/req_processed", nil
	}
	var toResDecorator = func(ctx context.Context, res string) (string, error) {
		fmt.Println("ToFunc Res Decorator : ", res)
		return res + "/res_processed", nil
	}
	var toExDecorator = func(ctx context.Context, req string, err error) error {
		return errors.Join(err, errors.New("custom raol error"))
	}
	toFuncBuilder.Func(toFn)
	toFuncBuilder.RequestDecorators(toReqDecorator, toReqDecorator)
	toFuncBuilder.ResponseDecorators(toResDecorator, toResDecorator)
	toFuncBuilder.ExceptionDecorators(toExDecorator)
	toFuncBuilder.PanicHandling(true)
	toFunction := toFuncBuilder.Build()

	builder, err := NewConnectorBuilder().
		FromFunction(fromFunction.Call).
		AdaptFunction(converter).
		ToFunction(toFunction.Call).
		Build()

	if err != nil {
		t.Errorf("Unexcept error %s ", err)
	}
	builder.Invoke(context.Background(), 10)
}
