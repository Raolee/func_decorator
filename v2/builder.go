package v2

import (
	"context"
)

type FunctionBuilder[REQ any, RES any] interface {
	Func(fn func(ctx context.Context, req REQ) (RES, error)) FunctionBuilder[REQ, RES]
	RequestMiddleware(fns ...func(ctx context.Context, req REQ) (REQ, error)) FunctionBuilder[REQ, RES]
	ResponseMiddleware(fns ...func(ctx context.Context, res RES) (RES, error)) FunctionBuilder[REQ, RES]
	ExceptionMiddleware(fns ...func(ctx context.Context, req REQ, err error) error) FunctionBuilder[REQ, RES]
	PanicHandling(accept bool) FunctionBuilder[REQ, RES]
	Build() *Function[REQ, RES]
}

type functionBuilder[REQ any, RES any] struct {
	function *Function[REQ, RES]
}

func NewFunctionBuilder[REQ any, RES any]() FunctionBuilder[REQ, RES] {
	return &functionBuilder[REQ, RES]{function: &Function[REQ, RES]{}}
}

func (f *functionBuilder[REQ, RES]) Func(fn func(ctx context.Context, req REQ) (RES, error)) FunctionBuilder[REQ, RES] {
	f.function.fn = fn
	return f
}

func (f *functionBuilder[REQ, RES]) RequestMiddleware(fns ...func(ctx context.Context, req REQ) (REQ, error)) FunctionBuilder[REQ, RES] {
	f.function.requestMiddleware = fns
	return f
}

func (f *functionBuilder[REQ, RES]) ResponseMiddleware(fns ...func(ctx context.Context, res RES) (RES, error)) FunctionBuilder[REQ, RES] {
	f.function.responseMiddleware = fns
	return f
}

func (f *functionBuilder[REQ, RES]) ExceptionMiddleware(fns ...func(ctx context.Context, req REQ, err error) error) FunctionBuilder[REQ, RES] {
	f.function.exceptionMiddleware = fns
	return f
}

func (f *functionBuilder[REQ, RES]) PanicHandling(accept bool) FunctionBuilder[REQ, RES] {
	f.function.panicHandling = accept
	return f
}

func (f *functionBuilder[REQ, RES]) Build() *Function[REQ, RES] {
	return f.function
}
