package func_decorator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestFunction_Call(t *testing.T) {

	testErr := errors.New("test error")
	type TestFuncType func(ctx context.Context, args ...any) ([]any, error)

	testCases := []struct {
		name                  string
		fn                    TestFuncType
		args                  []any
		isolatedBeforeFuncs   []IsolationFuncType
		composableBeforeFuncs []TestFuncType
		isolatedAfterFuncs    []IsolationFuncType
		composableAfterFuncs  []TestFuncType
		expectedResults       []any
		expectedError         error
	}{
		{
			name: "no error",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return []any{args[0]}, nil
			},
			args:                  []any{context.Background(), "test1"},
			isolatedBeforeFuncs:   []IsolationFuncType{},
			composableBeforeFuncs: []TestFuncType{},
			isolatedAfterFuncs:    []IsolationFuncType{},
			composableAfterFuncs:  []TestFuncType{},
			expectedResults:       []any{[]any{"test1"}},
			expectedError:         nil,
		},
		{
			name: "occur error",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return nil, testErr
			},
			args:                  []any{context.Background(), "test2"},
			isolatedBeforeFuncs:   []IsolationFuncType{},
			composableBeforeFuncs: []TestFuncType{},
			isolatedAfterFuncs:    []IsolationFuncType{},
			composableAfterFuncs:  []TestFuncType{},
			expectedResults:       nil,
			expectedError:         testErr,
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

			results, err := f.Call(context.Background(), tc.args...)

			if !reflect.DeepEqual(results, tc.expectedResults) {
				t.Errorf("Expected results %v, but got %v", tc.expectedResults, results)
			}

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
			}
		})
	}
}
