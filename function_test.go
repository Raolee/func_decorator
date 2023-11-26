package func_decorator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestFunction_Call(t *testing.T) {
	testCases := []struct {
		name                  string
		fn                    any
		args                  []any
		isolatedBeforeFuncs   []func(ctx context.Context, args ...any)
		composableBeforeFuncs []any
		isolatedAfterFuncs    []func(ctx context.Context, args ...any)
		composableAfterFuncs  []any
		expectedResults       []any
		expectedError         error
	}{
		{
			name: "Test case 1",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return []any{args[0]}, nil
			},
			args:                  []any{"test1"},
			isolatedBeforeFuncs:   []func(ctx context.Context, args ...any){},
			composableBeforeFuncs: []any{},
			isolatedAfterFuncs:    []func(ctx context.Context, args ...any){},
			composableAfterFuncs:  []any{},
			expectedResults:       []any{"test1"},
			expectedError:         nil,
		},
		{
			name: "Test case 2",
			fn: func(ctx context.Context, args ...any) ([]any, error) {
				return nil, errors.New("test error")
			},
			args:                  []any{"test2"},
			isolatedBeforeFuncs:   []func(ctx context.Context, args ...any){},
			composableBeforeFuncs: []any{},
			isolatedAfterFuncs:    []func(ctx context.Context, args ...any){},
			composableAfterFuncs:  []any{},
			expectedResults:       nil,
			expectedError:         errors.New("test error"),
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := &Function{
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

			if !reflect.DeepEqual(err, tc.expectedError) {
				t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
			}
		})
	}
}
