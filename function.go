package func_decorator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Function struct {
	fn                    any
	fnInputTypes          []reflect.Type
	fnOutputTypes         []reflect.Type
	isolatedBeforeFuncs   []func(ctx context.Context, args ...any)
	composableBeforeFuncs []any
	isolatedAfterFuncs    []func(ctx context.Context, args ...any)
	composableAfterFuncs  []any
}

func (f *Function) Call(ctx context.Context, args ...any) ([]any, error) {

	// 격리된 before 호출
	for _, isolatedBeforeFunc := range f.isolatedBeforeFuncs {
		isolatedBeforeFunc(ctx, args...)
	}

	// 조합된 before 호출
	for _, composableBeforeFunc := range f.composableBeforeFuncs {
		_, _ = executeFuncUsingReflect(composableBeforeFunc, args...)
	}

	// 메인 func 호출
	results, err := executeFuncUsingReflect(f.fn, args...)
	if err != nil {
		return nil, err // 핸들링되지 못한 error 를 잡음
	}

	// 격리된 after 호출
	for _, isolatedAfterFunc := range f.isolatedAfterFuncs {
		isolatedAfterFunc(ctx, args...)
	}

	// 조합된 after 호출
	for _, composableAfterFunc := range f.composableAfterFuncs {
		_, _ = executeFuncUsingReflect(composableAfterFunc, args...)
	}

	return results, nil
}

func executeFuncUsingReflect(fn any, args ...any) (results []any, err error) {
	defer func(e *error) {
		if r := recover(); r != nil {
			innerErr := errors.New(fmt.Sprintf("%v", r))
			e = &innerErr
		}
	}(&err)
	// reflect 를 이용한 함수 실행 준비
	fnValue := reflect.ValueOf(fn) // fn 을 reflect.Value 로 바꿈
	reflectInputs := ConvertAnySliceToReflectValues(args)
	reflectOutputs := fnValue.Call(reflectInputs)
	return ConvertReflectValuesToAnySlice(reflectOutputs), nil
}
