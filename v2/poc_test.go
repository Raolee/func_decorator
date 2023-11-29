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
	builder := NewFunctionBuilder[string, string]()

	var fn = func(ctx context.Context, req string) (string, error) {
		fmt.Println(req)
		return req + "/processed", nil
	}
	var requestInterceptor = func(ctx context.Context, req string) (string, error) {
		return req + "/req_processed", nil
	}
	var responseInterceptor = func(ctx context.Context, res string) (string, error) {
		return res + "/res_processed", nil
	}
	var exceptionInterceptor = func(ctx context.Context, req string, err error) error {
		return errors.Join(err, errors.New("custom raol error"))
	}

	builder.Func(fn)
	builder.RequestMiddleware(requestInterceptor, requestInterceptor)
	builder.ResponseMiddleware(responseInterceptor, responseInterceptor)
	builder.ExceptionMiddleware(exceptionInterceptor)
	builder.PanicHandling(true)
	function := builder.Build()

	reflectFuncCall(function.Call, context.Background(), "test")
}

func TestExecuteGenericFuncToAnotherGenericFuncUsingReflect(t *testing.T) {
	fromFuncBuilder := NewFunctionBuilder[int, int]()
	var fromFn = func(ctx context.Context, req int) (int, error) {
		fmt.Println("FromFunc ", req)
		return req, nil
	}
	var fromRequestInterceptor = func(ctx context.Context, req int) (int, error) {
		fmt.Println("FromFunc Req Interceptor ", req)
		return req + 1, nil
	}
	var fromResponseInterceptor = func(ctx context.Context, res int) (int, error) {
		fmt.Println("FromFunc Res Interceptor ", res)
		return res + 1, nil
	}
	var fromExceptionInterceptor = func(ctx context.Context, req int, err error) error {
		return errors.Join(err, errors.New("custom raol error"))
	}
	fromFuncBuilder.Func(fromFn)
	fromFuncBuilder.RequestMiddleware(fromRequestInterceptor, fromRequestInterceptor)
	fromFuncBuilder.ResponseMiddleware(fromResponseInterceptor, fromResponseInterceptor)
	fromFuncBuilder.ExceptionMiddleware(fromExceptionInterceptor)
	fromFuncBuilder.PanicHandling(true)
	fromFunction := fromFuncBuilder.Build()

	converter := func(ctx context.Context, req int) (nextReq string, err error) {
		fmt.Println("Adapt : ", req)
		return strconv.Itoa(req), nil
	}

	toFuncBuilder := NewFunctionBuilder[string, string]()
	var toFn = func(ctx context.Context, req string) (string, error) {
		fmt.Println("ToFunc : ", req)
		return req + "/processed", nil
	}
	var toRequestInterceptor = func(ctx context.Context, req string) (string, error) {
		fmt.Println("ToFunc Req Interceptor : ", req)
		return req + "/req_processed", nil
	}
	var toResponseInterceptor = func(ctx context.Context, res string) (string, error) {
		fmt.Println("ToFunc Res Interceptor : ", res)
		return res + "/res_processed", nil
	}
	var toExceptionInterceptor = func(ctx context.Context, req string, err error) error {
		return errors.Join(err, errors.New("custom raol error"))
	}
	toFuncBuilder.Func(toFn)
	toFuncBuilder.RequestMiddleware(toRequestInterceptor, toRequestInterceptor)
	toFuncBuilder.ResponseMiddleware(toResponseInterceptor, toResponseInterceptor)
	toFuncBuilder.ExceptionMiddleware(toExceptionInterceptor)
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
