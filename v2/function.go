package v2

import (
	"context"
	"errors"
	"fmt"
)

type Function[REQ any, RES any] struct {
	panicHandling       bool
	requestMiddleware   []func(ctx context.Context, req REQ) (REQ, error)
	fn                  func(ctx context.Context, req REQ) (RES, error)
	responseMiddleware  []func(ctx context.Context, res RES) (RES, error)
	exceptionMiddleware []func(ctx context.Context, req REQ, err error) error
}

func (f *Function[REQ, RES]) Call(ctx context.Context, req REQ) (res RES, err error) {
	if f.panicHandling {
		defer func(e *error) {
			if r := recover(); r != nil {
				innerErr := errors.New(fmt.Sprintf("%s", r)) // TODO : stack trace 찍게 해야 함
				*e = innerErr
			}
		}(&err)
	}
	res, err = f.call(ctx, req)

	if err != nil && len(f.exceptionMiddleware) > 0 {
		for _, exceptionInterceptor := range f.exceptionMiddleware {
			err = exceptionInterceptor(ctx, req, err)
		}
	}

	return res, err
}

func (f *Function[REQ, RES]) call(ctx context.Context, req REQ) (RES, error) {

	var err error
	// request interceptor 호출
	for _, requestInterceptor := range f.requestMiddleware {
		req, err = requestInterceptor(ctx, req)
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

	// request interceptor 호출
	for _, responseInterceptor := range f.responseMiddleware {
		res, err = responseInterceptor(ctx, res)
		if err != nil {
			return zeroValue[RES](), err
		}
	}

	return res, nil
}
