package v3

import (
	"context"
	"errors"
)

// FunctionType Function Registry 용 FunctionType을 정의합니다.
type FunctionType func(context.Context, any) (any, error)

type FunctionRegistry interface {
	RegisterFunction(name string, fn FunctionType)
	GetFunction(name string) (FunctionType, error)
}

type SimpleFunctionRegistry struct {
	functions map[string]FunctionType
}

func NewSimpleFunctionRegistry() FunctionRegistry {
	return &SimpleFunctionRegistry{
		functions: make(map[string]FunctionType),
	}
}

func (r *SimpleFunctionRegistry) RegisterFunction(name string, fn FunctionType) {
	r.functions[name] = fn
}

func (r *SimpleFunctionRegistry) GetFunction(name string) (FunctionType, error) {
	fn, ok := r.functions[name]
	if !ok {
		return nil, errors.New("function not found")
	}
	return fn, nil
}
