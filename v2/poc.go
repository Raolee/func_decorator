package v2

import (
	"context"
	"errors"
	"fmt"
	v1 "func_decorator/v1"
	"reflect"
)

type ConnectorBuilder interface {
	FromFunction(fn any) ConnectorBuilder
	AdaptFunction(fn any) ConnectorBuilder
	ToFunction(fn any) ConnectorBuilder
	Test() error
	Build() (Connector, error)
}

type connectorBuilder struct {
	fromFunc         any
	fromInputTypes   []reflect.Type
	fromOutputTypes  []reflect.Type
	adaptFunc        any
	adaptInputTypes  []reflect.Type
	adaptOutputTypes []reflect.Type
	toFunc           any
	toInputTypes     []reflect.Type
	toOutputTypes    []reflect.Type
}

func NewConnectorBuilder() ConnectorBuilder {
	return &connectorBuilder{}
}

// FromFunction | 시작할 FromFunction 을 주입합니다.
// 주입된 FromFunction 의 output 은 (any, error) 형태여야 합니다.
func (c *connectorBuilder) FromFunction(fn any) ConnectorBuilder {
	c.fromFunc = fn
	c.fromInputTypes, c.fromOutputTypes = v1.ExtractFuncType(fn)
	return c
}

// FromFunc 으로 주입된 func 이 규칙을 지키고 있는지 테스트 합니다.
func (c *connectorBuilder) testFromFunc() error {
	// fromFunc 은 nil 일 수 없음
	if c.fromFunc == nil {
		return errors.New("fromFunc can't nil")
	}
	// fromFunc 은 Func type 이어야 함
	if reflect.ValueOf(c.fromFunc).Kind() != reflect.Func {
		return errors.New("fromFunc isn't function type (reflect.Func)")
	}

	// input types 의 첫 번째는 context.Context 여야 함
	if !c.fromInputTypes[0].Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return errors.New("fromFunc first input type must be 'context.Context'")
	}

	// output types 는 (any, error) 처럼 2가지만 리턴해야 함
	if len(c.fromOutputTypes) != 2 {
		return errors.New("fromFunc output types length must be '2' (any, error)")
	}

	// output types 의 마지막은 error 타입이어야 함
	if !c.fromOutputTypes[len(c.fromOutputTypes)-1].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.New("last element in fromFunc output types must be 'error' interface")
	}
	return nil
}

// AdaptFunction | FromFunc 과 ToFunc 의 in/out types 을 맞춰주는 Adapt DecoratedFunction 을 주입합니다.
// 본 메서드는 필수가 아닙니다. (AdaptFunction 이 nil 일 경우, FromFunc 의 output 을 ToFunc 의 input 으로 바로 넣게됨)
// AdaptFunction 의 input 은 FromFunc 의 output 과 형식이 같아야 합니다.
// 또한, AdaptFunction 의 output 은 ToFunc 의 input 과 형식이 같아야 합니다.
func (c *connectorBuilder) AdaptFunction(fn any) ConnectorBuilder {
	c.adaptFunc = fn
	c.adaptInputTypes, c.adaptOutputTypes = v1.ExtractFuncType(fn)
	return c
}

// 주입된 AdaptFunction 이 규칙을 지키고 있는지 테스트 합니다.
func (c *connectorBuilder) testAdaptFunc() error {

	// adapt func 이 nil 이면,
	// fromFunc 의 output type 과 toFunc 의 input type 이 같아야 함
	if c.adaptFunc == nil && c.fromFunc != nil && c.toFunc != nil {
		if !reflect.DeepEqual(c.fromOutputTypes[0], c.toInputTypes[0]) {
			return errors.New("function mismatch: output types of 'fromFunc' does not match input types of 'toFunc'")
		}
		return nil
	}
	// adaptFunc 은 Func type 이어야 함
	if reflect.ValueOf(c.adaptFunc).Kind() != reflect.Func {
		return errors.New("adaptFunc isn't function type (reflect.Func)")
	}

	// input types 의 첫 번째는 context.Context 여야 함
	if !c.adaptInputTypes[0].Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return errors.New("adaptFunc first args type must be 'context.Context'")
	}

	// fromFunc 의 output types 은 adaptFunc 의 input types 와 같아야 함
	if c.fromOutputTypes[0].String() != c.adaptInputTypes[1].String() {
		return errors.New(fmt.Sprintf("function mismatch: output type[0](%s) of 'fromFunc' does not match input type[1](%s) of 'adaptFunc'", c.fromOutputTypes[0].String(), c.adaptInputTypes[1].String()))
	}
	// adaptFunc 의 output types 은 toFunc 의 input types 와 같아야 함
	if c.adaptOutputTypes[0].String() != c.toInputTypes[1].String() {
		return errors.New(fmt.Sprintf("function mismatch: output type[0](%s) of 'adaptFunc' does not match input type[1](%s) of 'toFunc'", c.adaptOutputTypes[0].String(), c.toInputTypes[1].String()))
	}
	return nil

}

// ToFunction | FromFunc 이후 호출할 ToFunc 을 주입합니다.
// 주입된 ToFunc 의 output types 는 (any, error) 형태여야 합니다.
func (c *connectorBuilder) ToFunction(fn any) ConnectorBuilder {
	c.toFunc = fn
	c.toInputTypes, c.toOutputTypes = v1.ExtractFuncType(fn)
	return c
}

// 주입된 ToFunction 이 규칙을 지키고 있는지 테스트 합니다.
func (c *connectorBuilder) testToFunc() error {
	if c.toFunc == nil {
		return errors.New("toFunc can't nil")
	}
	if reflect.ValueOf(c.toFunc).Kind() != reflect.Func {
		return errors.New("toFunc isn't function type (reflect.Func)")
	}

	// input types 의 첫 번째는 context.Context 여야 함
	if !c.toInputTypes[0].Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return errors.New("toFunc first args type must be 'context.Context'")
	}

	// output types 의 마지막은 error 타입이어야 함
	if !c.fromOutputTypes[len(c.fromOutputTypes)-1].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.New("last element in toFunc output types must be 'error' interface")
	}
	return nil
}

func (c *connectorBuilder) Test() error {
	if err := c.testFromFunc(); err != nil {
		return err
	}
	if err := c.testToFunc(); err != nil {
		return err
	}
	if err := c.testAdaptFunc(); err != nil {
		return err
	}
	return nil
}

func (c *connectorBuilder) Build() (Connector, error) {
	err := c.Test()
	if err != nil {
		return nil, err
	}
	return &connector{
		fromFunc:  c.fromFunc,
		adaptFunc: c.adaptFunc,
		toFunc:    c.toFunc,
	}, nil
}

type Connector interface {
	Invoke(args ...any) (any, error)
}

type connector struct {
	fromFunc  any
	adaptFunc any
	toFunc    any
}

func (c *connector) Invoke(args ...any) (any, error) {
	// args[0]의 context.Context Value
	ctxArgReflectValue := v1.ConvertToReflectValue(args[0])

	fromFuncValue := reflect.ValueOf(c.fromFunc)
	fromInputValues := v1.ConvertAnySliceToReflectValues(args)
	fromOutputValues := fromFuncValue.Call(fromInputValues)
	if err := hasError(fromOutputValues); err != nil {
		return nil, err
	}

	// args[0] = context.Context 를 계속 물고 다념
	var toFuncInputValues = append([]reflect.Value{ctxArgReflectValue}, fromOutputValues[:len(fromOutputValues)-1]...)
	if c.adaptFunc != nil {
		adaptFuncValue := reflect.ValueOf(c.adaptFunc)
		// adapt func 으로 DecoratedFunction Type 매칭 시킴
		adaptOutputValues := adaptFuncValue.Call(toFuncInputValues)
		if err := hasError(adaptOutputValues); err != nil {
			return nil, err
		}

		toFuncInputValues = append([]reflect.Value{ctxArgReflectValue}, adaptOutputValues[:len(adaptOutputValues)-1]...)
	}

	toFuncValue := reflect.ValueOf(c.toFunc)
	toFuncOutputValues := toFuncValue.Call(toFuncInputValues)
	if err := hasError(toFuncOutputValues); err != nil {
		return nil, err
	}
	return toFuncOutputValues[0], nil
}

func hasError(result []reflect.Value) error {
	if len(result) > 0 {
		// 사실 function 을 만들 때는 output 중에서 맨 마지막을 error 로 명시하기 때문에 맨 뒤만 보면 됨
		if err, ok := result[len(result)-1].Interface().(error); ok && err != nil {
			return err
		}
	}
	return nil
}

func reflectFuncCall(fn any, args ...any) {
	funcValue := reflect.ValueOf(fn)
	inputValues := v1.ConvertAnySliceToReflectValues(args)
	funcValue.Call(inputValues)
}
