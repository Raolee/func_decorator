package v1

import (
	"context"
	"fmt"
	"reflect"
)

func ExtractFuncType(fn any) ([]reflect.Type, []reflect.Type) {
	if fn == nil {
		return nil, nil
	}

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return nil, nil
	}

	numIn := fnType.NumIn()
	inTypes := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		inTypes[i] = fnType.In(i)
	}

	numOut := fnType.NumOut()
	outTypes := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		outTypes[i] = fnType.Out(i)
	}

	return inTypes, outTypes
}

func ConvertReflectValuesToAnySlice(rvs []reflect.Value) []any {
	anys := make([]any, len(rvs))
	for i, rv := range rvs {
		anys[i] = convertToAny(rv)
	}
	return anys
}

func ConvertAnySliceToReflectValues(anySlice []any) []reflect.Value {
	rvs := make([]reflect.Value, len(anySlice))

	for i, a := range anySlice {
		rvs[i] = convertToReflectValue(a)
	}
	return rvs
}

func convertToAny(rv reflect.Value) any {
	// reflect.Value 자체가 유효한지 먼저 체크
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil // Ptr 선언만 되어 있고 값이 비어있음
		}
		return rv.Interface() // Ptr 자체의 주소를 리턴
	}

	return rv.Interface() // Value가 가진 값을 리턴
}

func convertToReflectValue(value any) reflect.Value {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice:
		sliceValue := reflect.ValueOf(value)
		rv := reflect.MakeSlice(sliceValue.Type(), sliceValue.Len(), sliceValue.Cap())
		for i := 0; i < sliceValue.Len(); i++ {
			rv.Index(i).Set(convertToReflectValue(sliceValue.Index(i).Interface()))
		}
		return rv
	case reflect.Ptr:
		if _, ok := value.(context.Context); ok {
			return reflect.ValueOf(value) // 컨텍스트는... 따로 처리해야한다 -_-
		}
		return convertToReflectValue(v.Elem().Interface())
	default:
		return reflect.ValueOf(value)
	}
}

func ConvertReflectValuesToReflectTypes(rvs []reflect.Value) []reflect.Type {
	rts := make([]reflect.Type, len(rvs))
	for i, rv := range rvs {
		rts[i] = rv.Type()
	}
	return rts
}

func printReflects(rs []reflect.Type) {
	for i := range rs {
		fmt.Println(rs[i].Name())
	}
}
