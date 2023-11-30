package v2

import (
	"context"
	"errors"
	"reflect"
)

// VoidType | AnyFunction 은 반드시 req, res 명시된 func 만 받기 때문에 void 형태의 Func 을 우회하여 받기 위해 정의
// ex) 아래 구조인 경우가 바로 Void func
// func (ctx context.Context, _ VoidType) (_ VoidType, err error)
type VoidType struct{}

// Signature 는 AnyFunction 의 Input, Output 의 Void 여부를 나타냄
type Signature int

const (
	NoRequestNoResponse         = Signature(0)
	NoRequestExistsResponse     = Signature(1)
	ExistsRequestNoResponse     = Signature(2)
	ExistsRequestExistsResponse = Signature(3)
)

// AnyFunction | 제네럴 하게 사용하는 AnyFunction Interface
type AnyFunction interface {
	Call(ctx context.Context, req any) (res any, err error)
	GetSignature() Signature
	GetRequestType() reflect.Type
	GetResponseType() reflect.Type
}

// 인터페이스 구현 구조체
type anyFunction struct {
	reqType reflect.Type
	resType reflect.Type
	sig     Signature
	fn      func(ctx context.Context, req any) (res any, err error)
}

// NewAnyFunction | AnyFunction 의 생성자
// reqType, resType, fn 은 필수(required)
func NewAnyFunction(reqType, resType reflect.Type, fn func(ctx context.Context, req any) (res any, err error)) (AnyFunction, error) {
	if reqType == nil || resType == nil {
		return nil, errors.New("reqType and resType must not be nil")
	}
	if fn == nil {
		return nil, errors.New("fn must not be nil")
	}
	var sig Signature
	switch {
	case EqualType(reqType, VoidType{}) && EqualType(resType, VoidType{}):
		sig = NoRequestNoResponse
	case EqualType(reqType, VoidType{}):
		sig = NoRequestExistsResponse
	case EqualType(resType, VoidType{}):
		sig = ExistsRequestNoResponse
	default:
		sig = ExistsRequestExistsResponse
	}
	return &anyFunction{
		reqType: reqType,
		resType: resType,
		sig:     sig,
		fn:      fn,
	}, nil
}

func (a *anyFunction) Call(ctx context.Context, req any) (res any, err error) {
	return a.fn(ctx, req)
}

func (a *anyFunction) GetSignature() Signature {
	return a.sig
}

func (a *anyFunction) GetRequestType() reflect.Type {
	return a.reqType
}

func (a *anyFunction) GetResponseType() reflect.Type {
	return a.resType
}
