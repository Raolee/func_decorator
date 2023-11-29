package v2

import (
	"context"
)

type DecoratedFunctionBuilder[REQ any, RES any] interface {
	Func(fn func(ctx context.Context, req REQ) (RES, error)) DecoratedFunctionBuilder[REQ, RES]
	RequestDecorators(fns ...func(ctx context.Context, req REQ) (REQ, error)) DecoratedFunctionBuilder[REQ, RES]
	ResponseDecorators(fns ...func(ctx context.Context, res RES) (RES, error)) DecoratedFunctionBuilder[REQ, RES]
	ExceptionDecorators(fns ...func(ctx context.Context, req REQ, err error) error) DecoratedFunctionBuilder[REQ, RES]
	PanicHandling(accept bool) DecoratedFunctionBuilder[REQ, RES]
	Build() *DecoratedFunction[REQ, RES]
}

type decoratedFunctionBuilder[REQ any, RES any] struct {
	function *DecoratedFunction[REQ, RES]
}

func NewDecoratedFunctionBuilder[REQ any, RES any]() DecoratedFunctionBuilder[REQ, RES] {
	return &decoratedFunctionBuilder[REQ, RES]{function: &DecoratedFunction[REQ, RES]{}}
}

func (f *decoratedFunctionBuilder[REQ, RES]) Func(fn func(ctx context.Context, req REQ) (RES, error)) DecoratedFunctionBuilder[REQ, RES] {
	f.function.fn = fn
	return f
}

func (f *decoratedFunctionBuilder[REQ, RES]) RequestDecorators(fns ...func(ctx context.Context, req REQ) (REQ, error)) DecoratedFunctionBuilder[REQ, RES] {
	f.function.requestDecorators = fns
	return f
}

func (f *decoratedFunctionBuilder[REQ, RES]) ResponseDecorators(fns ...func(ctx context.Context, res RES) (RES, error)) DecoratedFunctionBuilder[REQ, RES] {
	f.function.responseDecorators = fns
	return f
}

func (f *decoratedFunctionBuilder[REQ, RES]) ExceptionDecorators(fns ...func(ctx context.Context, req REQ, err error) error) DecoratedFunctionBuilder[REQ, RES] {
	f.function.exceptionDecorators = fns
	return f
}

func (f *decoratedFunctionBuilder[REQ, RES]) PanicHandling(accept bool) DecoratedFunctionBuilder[REQ, RES] {
	f.function.panicHandling = accept
	return f
}

func (f *decoratedFunctionBuilder[REQ, RES]) Build() *DecoratedFunction[REQ, RES] {
	return f.function
}
