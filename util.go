package func_decorator

import (
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
		anys[i] = rv.Interface()
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