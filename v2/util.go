package v2

import (
	"reflect"
)

// GetGenericType | 제네릭 타입을 reflect.Type 으로 리턴합니다.
func GetGenericType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
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
