package v1

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type IsolationFuncType func(ctx context.Context, args ...any) (context.Context, error)
type ComposableFuncType func(ctx context.Context) (context.Context, error)

type Function[T any] struct {
	fn                    T
	fnInputTypes          []reflect.Type
	fnOutputTypes         []reflect.Type
	isolatedBeforeFuncs   []IsolationFuncType
	composableBeforeFuncs []ComposableFuncType
	isolatedAfterFuncs    []IsolationFuncType
	composableAfterFuncs  []ComposableFuncType
}

func (f *Function[T]) Call(ctx context.Context, args ...any) ([]any, error) {
	var err error

	// 격리된 before 호출
	for _, isolatedBeforeFunc := range f.isolatedBeforeFuncs {
		ctx, err = isolatedBeforeFunc(ctx, args...)
		if err != nil {
			return nil, err
		}
	}

	defer func(c *context.Context) {
		// 격리된 after 호출
		for _, isolatedAfterFunc := range f.isolatedAfterFuncs {
			_, _ = isolatedAfterFunc(*c, args...)
		}
	}(&ctx)

	// 조합된 before 호출
	for _, composableBeforeFunc := range f.composableBeforeFuncs {
		ctx, err = composableBeforeFunc(ctx)
		if err != nil {
			return nil, err
		}
	}

	// 메인 func 호출
	results, execError := executeFuncUsingReflect(f.fn, ctx, args...)
	if execError != nil {
		return nil, execError // 핸들링되지 못한 error 를 잡음
	}

	// 조합된 after 호출
	for _, composableAfterFunc := range f.composableAfterFuncs {
		ctx, err = composableAfterFunc(ctx)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func executeFuncUsingReflect[T any](fn T, ctx context.Context, args ...any) (results []any, err error) {
	defer func(e *error) {
		if r := recover(); r != nil {
			innerErr := errors.New(fmt.Sprintf("%v", r))
			*e = innerErr
		}
	}(&err)
	// reflect 를 이용한 함수 실행 준비
	fnValue := reflect.ValueOf(fn)                                                                                // fn 을 reflect.Value 로 바꿈
	reflectInputs := append([]reflect.Value{ConvertToReflectValue(ctx)}, ConvertAnySliceToReflectValues(args)...) // 맨 앞에 ctx를 넣음
	reflectOutputs := fnValue.Call(reflectInputs)

	// 메인 func 결과 중 맨 마지막 output 을 에러로 가정하고 체크
	if len(reflectOutputs) > 0 {
		// 사실 function 을 만들 때는 output 중에서 맨 마지막을 error 로 명시하기 때문에 맨 뒤만 보면 됨
		if e, ok := reflectOutputs[len(reflectOutputs)-1].Interface().(error); ok && e != nil {
			return nil, e
		}

		// 마지막 error를 제거하고 리턴
		return ConvertReflectValuesToAnySlice(reflectOutputs[:len(reflectOutputs)-1]), nil
	}
	return nil, nil
}

//type Function struct {
//	fn                    any
//	fnInputTypes          []reflect.Type
//	fnOutputTypes         []reflect.Type
//	isolatedBeforeFuncs   []IsolationFuncType
//	composableBeforeFuncs []any
//	isolatedAfterFuncs    []IsolationFuncType
//	composableAfterFuncs  []any
//}
//
//func (f *Function) Call(ctx context.Context, args ...any) ([]any, error) {
//
//	// 격리된 before 호출
//	for _, isolatedBeforeFunc := range f.isolatedBeforeFuncs {
//		isolatedBeforeFunc(ctx, args...)
//	}
//
//	// 조합된 before 호출
//	for _, composableBeforeFunc := range f.composableBeforeFuncs {
//		_, _ = executeFuncUsingReflect(composableBeforeFunc, args...)
//	}
//
//	// 메인 func 호출
//	results, execError := executeFuncUsingReflect(f.fn, args...)
//	if execError != nil {
//		return nil, execError // 핸들링되지 못한 error 를 잡음
//	}
//
//	// 격리된 after 호출
//	for _, isolatedAfterFunc := range f.isolatedAfterFuncs {
//		isolatedAfterFunc(ctx, args...)
//	}
//
//	// 조합된 after 호출
//	for _, composableAfterFunc := range f.composableAfterFuncs {
//		_, _ = executeFuncUsingReflect(composableAfterFunc, args...)
//	}
//
//	return results, nil
//}
//
//func executeFuncUsingReflect(fn any, args ...any) (results []any, err error) {
//	defer func(e *error) {
//		if r := recover(); r != nil {
//			innerErr := errors.New(fmt.Sprintf("%v", r))
//			*e = innerErr
//		}
//	}(&err)
//	// reflect 를 이용한 함수 실행 준비
//	fnValue := reflect.ValueOf(fn) // fn 을 reflect.Value 로 바꿈
//	reflectInputs := ConvertAnySliceToReflectValues(args)
//	reflectOutputs := fnValue.Call(reflectInputs)
//
//	// 메인 func 결과 중 맨 마지막 output 을 에러로 가정하고 체크
//	if len(reflectOutputs) > 0 {
//		// 사실 function 을 만들 때는 output 중에서 맨 마지막을 error 로 명시하기 때문에 맨 뒤만 보면 됨
//		if err, ok := reflectOutputs[len(reflectOutputs)-1].Interface().(error); ok && err != nil {
//			return nil, err
//		}
//
//		// 마지막 error를 제거하고 리턴
//		return ConvertReflectValuesToAnySlice(reflectOutputs[:len(reflectOutputs)-1]), nil
//	}
//	return nil, nil
//}
