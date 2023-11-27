package func_decorator

import (
	"errors"
	"fmt"
	"reflect"
)

type FunctionBuilder[T any] interface {
	Func(fn T) FunctionBuilder[T]
	BeforeIsolation(fns ...IsolationFuncType) FunctionBuilder[T]
	BeforeComposition(fns ...ComposableFuncType) FunctionBuilder[T]
	AfterIsolation(fns ...IsolationFuncType) FunctionBuilder[T]
	AfterComposition(fns ...ComposableFuncType) FunctionBuilder[T]
	Test() error
	Build() (*Function[T], error)
}

type functionBuilder[T any] struct {
	function *Function[T]
}

func NewFunctionBuilder[T any]() FunctionBuilder[T] {
	return &functionBuilder[T]{function: &Function[T]{}}
}

func (fb *functionBuilder[T]) Func(fn T) FunctionBuilder[T] {
	fb.function.fn = fn
	fb.function.fnInputTypes, fb.function.fnOutputTypes = ExtractFuncType(fn)
	return fb
}

func (fb *functionBuilder[T]) testFunc() error {
	//if fb.function.fn == nil {
	//	return errors.New("fn is nil")
	//}
	if reflect.ValueOf(fb.function.fn).Kind() != reflect.Func {
		return errors.New("fn isn't function type (reflect.Func)")
	}
	return nil
}

func (fb *functionBuilder[T]) BeforeIsolation(fns ...IsolationFuncType) FunctionBuilder[T] {
	fb.function.isolatedBeforeFuncs = fns
	return fb
}

func (fb *functionBuilder[T]) testBeforeIsolation() error {
	// func slice 가 비어거나 null이 없어야 함
	if len(fb.function.isolatedBeforeFuncs) > 0 {
		for i, beforeFunc := range fb.function.isolatedBeforeFuncs {
			if beforeFunc == nil {
				return errors.New(fmt.Sprintf("isolatedBeforeFuncs[%d] is nil", i))
			}
		}
	}
	return nil
}

func (fb *functionBuilder[T]) BeforeComposition(fns ...ComposableFuncType) FunctionBuilder[T] {
	fb.function.composableBeforeFuncs = fns
	return fb
}

func (fb *functionBuilder[T]) testBeforeComposition() error {
	inputTypes, _ := ExtractFuncType(fb.function.fn)
	if len(fb.function.composableBeforeFuncs) > 0 {
		for i, composableBefore := range fb.function.composableBeforeFuncs {

			// composableBefore 는 nil 일 수 없음
			//if composableBefore == nil {
			//	return errors.New(fmt.Sprintf("composableBefores[%d] is nil", i))
			//}

			// composableBefore 는 function 이어야 함
			if reflect.ValueOf(composableBefore).Kind() != reflect.Func {
				return errors.New(fmt.Sprintf("composableBefores[%d] isn't function type (reflect.Func)", i))
			}

			// composableBefore 의 input types는 은 f.function.fn 의 inputTypes 와 같아야 함
			composableBeforeInputTypes, _ := ExtractFuncType(composableBefore)

			if !reflect.DeepEqual(inputTypes, composableBeforeInputTypes) {
				return errors.New(fmt.Sprintf("composableBefores[%d] output types (%v) is not equal function input types (%v)", i, composableBeforeInputTypes, inputTypes))
			}
		}
	}

	return nil
}

func (fb *functionBuilder[T]) AfterIsolation(fns ...IsolationFuncType) FunctionBuilder[T] {
	fb.function.isolatedAfterFuncs = fns
	return fb
}

func (fb *functionBuilder[T]) testAfterIsolation() error {
	// func slice 가 비어거나 null이 없어야 함
	if len(fb.function.isolatedAfterFuncs) > 0 {
		for i, afterFunc := range fb.function.isolatedAfterFuncs {
			if afterFunc == nil {
				return errors.New(fmt.Sprintf("isolatedAfterFuncs[%d] is nil", i))
			}
		}
	}
	return nil
}

func (fb *functionBuilder[T]) AfterComposition(fns ...ComposableFuncType) FunctionBuilder[T] {
	fb.function.composableAfterFuncs = fns
	return fb
}

func (fb *functionBuilder[T]) testAfterComposition() error {
	inputTypes, _ := ExtractFuncType(fb.function.fn)
	if len(fb.function.composableAfterFuncs) > 0 {
		for i, composableAfter := range fb.function.composableAfterFuncs {

			// composableBefore 는 nil 일 수 없음
			//if composableAfter == nil {
			//	return errors.New(fmt.Sprintf("composableAfters[%d] is nil", i))
			//}

			// composableBefore 는 function 이어야 함
			if reflect.ValueOf(composableAfter).Kind() != reflect.Func {
				return errors.New(fmt.Sprintf("composableAfters[%d] isn't function type (reflect.Func)", i))
			}

			// composableBefore 의 input types는 은 f.function.fn 의 inputTypes 와 같아야 함
			composableAfterInputTypes, _ := ExtractFuncType(composableAfter)

			if !reflect.DeepEqual(inputTypes, composableAfterInputTypes) {
				return errors.New(fmt.Sprintf("composableAfters[%d] output types (%v) is not equal function input types (%v)", i, composableAfterInputTypes, inputTypes))
			}
		}
	}

	return nil
}

func (fb *functionBuilder[T]) Test() error {
	err := fb.testBeforeIsolation()
	if err != nil {
		return err
	}
	err = fb.testBeforeComposition()
	if err != nil {
		return err
	}
	err = fb.testFunc()
	if err != nil {
		return err
	}
	err = fb.testAfterIsolation()
	if err != nil {
		return err
	}
	err = fb.testAfterComposition()
	if err != nil {
		return err
	}
	return nil
}

func (fb *functionBuilder[T]) Build() (*Function[T], error) {
	if err := fb.Test(); err != nil {
		return nil, err
	}
	return fb.function, nil
}
