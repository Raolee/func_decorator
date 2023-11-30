package v2

import (
	"context"
	"reflect"
)

// GetGenericType | 제네릭 타입을 reflect.Type 으로 리턴합니다.
func GetGenericType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// EqualType | t 와 v 의 타입이 같은지 확인 합니다.
func EqualType(t reflect.Type, v any) bool {
	// v의 reflect.Type 을 가져옴
	valueType := reflect.TypeOf(v)

	// nil 값에 대한 처리
	if t == nil || valueType == nil {
		return t == valueType
	}

	// 포인터와 관련된 타입일 경우, Elem()을 사용하여 실제 타입 비교
	if t.Kind() == reflect.Ptr || valueType.Kind() == reflect.Ptr {
		return reflect.PtrTo(t).AssignableTo(reflect.PtrTo(valueType)) ||
			reflect.PtrTo(valueType).AssignableTo(reflect.PtrTo(t))
	}

	// slice, struct, interface 등 다른 케이스들
	return t.AssignableTo(valueType) || valueType.AssignableTo(t)
}

// zeroValue | 제네릭 타입의 zero value 를 반환합니다.
func zeroValue[T any]() T {
	zeroType := reflect.TypeOf((*T)(nil)).Elem()

	// T가 인터페이스인 경우는 따로 처리
	// 인터페이스의 zero value 는 인터페이스가 가진 모든 method 의 zero value 를 가진 형태이므로 nil이 아니기 때문임
	if zeroType.Kind() == reflect.Interface {
		var nilInterface T
		return nilInterface
	}

	return reflect.Zero(zeroType).Interface().(T)
}

func zeroValueWithType(t reflect.Type) any {
	return reflect.Zero(t).Interface()
}

func SetNodeFlowInContext(ctx context.Context, next string) context.Context {
	if value := ctx.Value("nodeFlow"); value == nil {
		return context.WithValue(ctx, "nodeFlow", next)
	} else if nodeFlow, ok := value.(string); ok {
		return context.WithValue(ctx, "nodeFlow", nodeFlow+"/"+next)
	} else {
		return context.WithValue(ctx, "nodeFlow", next)
	}
}

func GetNodeFlowInContext(ctx context.Context) string {
	value := ctx.Value("nodeFlow")
	if value != nil {
		if flow, ok := value.(string); ok {
			return flow
		}
	}
	return ""
}
