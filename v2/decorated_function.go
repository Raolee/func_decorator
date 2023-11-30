package v2

import (
	"context"
	"errors"
	"fmt"
)

type DecoratedFunction[REQ any, RES any] struct {
	panicHandling       bool
	requestDecorators   []func(ctx context.Context, req REQ) (REQ, error)
	fn                  func(ctx context.Context, req REQ) (RES, error)
	responseDecorators  []func(ctx context.Context, res RES) (RES, error)
	exceptionDecorators []func(ctx context.Context, req REQ, err error) error
}

func (f *DecoratedFunction[REQ, RES]) Call(ctx context.Context, req REQ) (res RES, err error) {
	// panic 핸들링 처리
	if f.panicHandling {
		defer func(e *error) {
			if r := recover(); r != nil {
				innerErr := errors.New(fmt.Sprintf("%s", r)) // TODO : stack trace 찍게 해야 함
				*e = innerErr
			}
		}(&err)
	}

	res, err = f.call(ctx, req)

	// 예외 데코레이터 처리
	if err != nil && len(f.exceptionDecorators) > 0 {
		for _, exDecorator := range f.exceptionDecorators {
			err = exDecorator(ctx, req, err)
		}
	}

	return res, err
}

func (f *DecoratedFunction[REQ, RES]) call(ctx context.Context, req REQ) (RES, error) {

	var err error
	// request 데코레이터 수행
	for _, reqDecorator := range f.requestDecorators {
		req, err = reqDecorator(ctx, req)
		if err != nil {
			return zeroValue[RES](), err
		}
	}

	// 본 func 호출
	var res RES
	res, err = f.fn(ctx, req)
	if err != nil {
		return zeroValue[RES](), err
	}

	// response 데코레이터 수행
	for _, resDecorator := range f.responseDecorators {
		res, err = resDecorator(ctx, res)
		if err != nil {
			return zeroValue[RES](), err
		}
	}

	return res, nil
}

func (f *DecoratedFunction[REQ, RES]) Any() (AnyFunction, error) {
	return NewAnyFunction(
		GetGenericType[REQ](),
		GetGenericType[RES](),
		func(ctx context.Context, req any) (res any, err error) {
			reqTyped, ok := req.(REQ)
			if !ok {
				return zeroValue[RES](), errors.New("REQ type convert failed")
			}
			return f.Call(ctx, reqTyped)
		})
}
