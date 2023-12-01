package v2

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestConnectFunctionNode(t *testing.T) {
	registry := NewFunctionRegistry()

	// 함수 등록

	// int -> string
	func1, _ := NewAnyFunction(reflect.TypeOf(0), reflect.TypeOf(""), func(ctx context.Context, req any) (res any, err error) {
		v, ok := req.(int)
		if !ok {
			return nil, fmt.Errorf("req type is not '%s'. req = [%v]", reflect.Int.String(), req)
		}
		return any(strconv.Itoa(v)), nil
	})

	// string -> int
	func2, _ := NewAnyFunction(reflect.TypeOf(""), reflect.TypeOf(0), func(ctx context.Context, req any) (res any, err error) {
		v, ok := req.(string)
		if !ok {
			return nil, fmt.Errorf("req type is not '%s'. req = [%v]", reflect.String.String(), req)
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		return i, nil
	})

	// int -> int
	func3, _ := NewAnyFunction(reflect.TypeOf(0), reflect.TypeOf(0), func(ctx context.Context, req any) (res any, err error) {
		v, ok := req.(int)
		if !ok {
			return nil, fmt.Errorf("req type is not '%s'. req = [%v]", reflect.Int.String(), req)
		}
		return any(v), nil
	})

	registry.RegisterFunction("func1", func1)
	registry.RegisterFunction("func2", func2)
	registry.RegisterFunction("func3", func3)

	// 1. 기본 연결 테스트
	if err := registry.ConnectFunctionNode("func1", "func2"); err != nil {
		t.Errorf("Failed to connect func1 to func2: %v", err)
	}

	// 2. 동일 노드 연결 시도
	if err := registry.ConnectFunctionNode("func1", "func1"); err == nil {
		t.Errorf("Should fail when connecting a node to itself")
	}

	// 3. 없는 노드 연결 시도
	if err := registry.ConnectFunctionNode("func1", "func4"); err == nil {
		t.Errorf("Should fail when trying to connect to a non-existent node")
	}

	// 4. 이미 연결된 노드 재연결 시도
	if err := registry.ConnectFunctionNode("func1", "func2"); err == nil {
		t.Errorf("Should fail when trying to connect already connected nodes")
	}

	// 5. 타입 불일치 연결 시도
	if err := registry.ConnectFunctionNode("func1", "func3"); err == nil {
		t.Errorf("Should fail when trying to connect incompatible types without an adapter")
	}

	// 6. 어댑터를 사용한 타입 맞춤 연결
	// string -> int
	adapterFunc, _ := NewAnyFunction(reflect.TypeOf(""), reflect.TypeOf(0), func(ctx context.Context, req any) (res any, err error) {
		v, ok := req.(string)
		if !ok {
			return nil, fmt.Errorf("req type is not '%s'. req = [%v]", reflect.String.String(), req)
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		return i, nil
	})
	if err := registry.ConnectFunctionNode("func1", "func3", adapterFunc); err != nil {
		t.Errorf("Failed to connect func1 to func3 using an adapter: %v", err)
	}
}
