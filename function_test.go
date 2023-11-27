package func_decorator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// 아래 테스트 코드에서 사용할 error 정의
var testErr = errors.New("test error")

// 아래 테스트 코드에서 사용할 func type 정의
type TestFuncType func(ctx context.Context, args ...any) ([]any, error)

func NoErrorTestFunc(ctx context.Context, args ...any) ([]any, error) {
	return []any{args[0]}, nil
}
func OccurErrorTestFunc(ctx context.Context, args ...any) ([]any, error) {
	return nil, testErr
}

func IsolatedBeforeTestFunc(ctx context.Context, args ...any) (context.Context, error) {
	value := ctx.Value("test")
	fmt.Println("before : ", value)
	return ctx, nil
}
func IsolatedAfterTestFunc(ctx context.Context, args ...any) (context.Context, error) {
	value := ctx.Value("test")
	fmt.Println("after : ", value)
	return ctx, nil
}

func TestFunction_Call(t *testing.T) {

	testCases := []struct {
		name                  string
		fn                    TestFuncType
		ctx                   context.Context
		args                  []any
		isolatedBeforeFuncs   []IsolationFuncType
		composableBeforeFuncs []ComposableFuncType
		isolatedAfterFuncs    []IsolationFuncType
		composableAfterFuncs  []ComposableFuncType
		expectedResults       []any
		expectedError         error
	}{
		{
			name: "only func - no error",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return []any{args[0]}, nil
			},
			ctx:                   context.WithValue(context.Background(), "test", "my ctx"),
			args:                  []any{"test1"},
			isolatedBeforeFuncs:   []IsolationFuncType{},
			composableBeforeFuncs: []ComposableFuncType{},
			isolatedAfterFuncs:    []IsolationFuncType{},
			composableAfterFuncs:  []ComposableFuncType{},
			expectedResults:       []any{[]any{"test1"}},
			expectedError:         nil,
		},
		{
			name: "only func - occur error",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return nil, testErr
			},
			ctx:                   context.WithValue(context.Background(), "test", "my ctx"),
			args:                  []any{"test decorating"},
			isolatedBeforeFuncs:   []IsolationFuncType{},
			composableBeforeFuncs: []ComposableFuncType{},
			isolatedAfterFuncs:    []IsolationFuncType{},
			composableAfterFuncs:  []ComposableFuncType{},
			expectedResults:       nil,
			expectedError:         testErr,
		},
		{
			name: "with isolated befor and after funcs - no error",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return []any{args[0]}, nil
			},
			ctx:  context.WithValue(context.Background(), "test", "my ctx"),
			args: []any{"test1"},
			isolatedBeforeFuncs: []IsolationFuncType{
				IsolatedBeforeTestFunc,
				IsolatedBeforeTestFunc,
			},
			composableBeforeFuncs: []ComposableFuncType{},
			isolatedAfterFuncs: []IsolationFuncType{
				IsolatedAfterTestFunc,
				IsolatedAfterTestFunc,
			},
			composableAfterFuncs: []ComposableFuncType{},
			expectedResults:      []any{[]any{"test1"}},
			expectedError:        nil,
		},
		{
			name: "with composable befor and after funcs - no error",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return []any{args[0]}, nil
			},
			ctx:                 context.WithValue(context.Background(), "test", "my ctx"),
			args:                []any{"test1"},
			isolatedBeforeFuncs: []IsolationFuncType{},
			composableBeforeFuncs: []ComposableFuncType{
				func(ctx context.Context) (context.Context, error) {
					newCtx := context.WithValue(ctx, "beforecompsable", "raol")
					return newCtx, nil
				},
			},
			isolatedAfterFuncs: []IsolationFuncType{
				func(ctx context.Context, args ...any) (context.Context, error) {
					fmt.Println(ctx.Value("test"))
					fmt.Println(ctx.Value("beforecompsable"))
					fmt.Println(ctx.Value("aftercompsable"))
					return ctx, nil
				},
			},
			composableAfterFuncs: []ComposableFuncType{
				func(ctx context.Context) (context.Context, error) {
					newCtx := context.WithValue(ctx, "aftercompsable", "raol")
					return newCtx, nil
				},
			},
			expectedResults: []any{[]any{"test1"}},
			expectedError:   nil,
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := &Function[TestFuncType]{
				fn:                    tc.fn,
				fnInputTypes:          []reflect.Type{},
				fnOutputTypes:         []reflect.Type{},
				isolatedBeforeFuncs:   tc.isolatedBeforeFuncs,
				composableBeforeFuncs: tc.composableBeforeFuncs,
				isolatedAfterFuncs:    tc.isolatedAfterFuncs,
				composableAfterFuncs:  tc.composableAfterFuncs,
			}

			results, err := f.Call(tc.ctx, tc.args...)

			if !reflect.DeepEqual(results, tc.expectedResults) {
				t.Errorf("Expected results %v, but got %v", tc.expectedResults, results)
			}

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
			}
		})
	}
}
