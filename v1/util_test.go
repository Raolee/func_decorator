package v1

import (
	"fmt"
	"reflect"
	"testing"
)

type DummyInterface interface {
	dummy()
}

type DummyStruct struct {
}

func (d *DummyStruct) dummy() {
}

func TestExtractFuncType(t *testing.T) {
	testCases := []struct {
		name             string
		fn               any
		expectedInTypes  []reflect.Type
		expectedOutTypes []reflect.Type
	}{
		{
			name: "int, string inputs, int string outputs",
			fn: func(a int, b string) (int, string) {
				return a, b
			},
			expectedInTypes:  []reflect.Type{reflect.TypeOf(0), reflect.TypeOf("")}, // input : int, string
			expectedOutTypes: []reflect.Type{reflect.TypeOf(0), reflect.TypeOf("")}, // output : int, string
		},
		{
			name:             "no inputs, no outputs",
			fn:               func() {},
			expectedInTypes:  []reflect.Type{},
			expectedOutTypes: []reflect.Type{},
		},
		{
			name:             "no func type",
			fn:               "not_a_function",
			expectedInTypes:  nil,
			expectedOutTypes: nil,
		},
		{
			name:             "nil func",
			fn:               nil,
			expectedInTypes:  nil,
			expectedOutTypes: nil,
		},
		{
			name: "no inputs, error output",
			fn: func() error {
				return nil
			},
			expectedInTypes:  []reflect.Type{},
			expectedOutTypes: []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()},
		},
		{
			name:             "struct inputs, struct,error outputs",
			fn:               func(r DummyStruct) (DummyStruct, error) { return DummyStruct{}, nil },
			expectedInTypes:  []reflect.Type{reflect.TypeOf(DummyStruct{})},
			expectedOutTypes: []reflect.Type{reflect.TypeOf(DummyStruct{}), reflect.TypeOf((*error)(nil)).Elem()},
		},
		{
			name:             "struct_ptr inputs, struct_ptr outputs",
			fn:               func(r *DummyStruct) *DummyStruct { return nil },
			expectedInTypes:  []reflect.Type{reflect.TypeOf(&DummyStruct{})},
			expectedOutTypes: []reflect.Type{reflect.TypeOf(&DummyStruct{})},
		},
		{
			name:             "interface inputs, interface outputs",
			fn:               func(r DummyInterface) DummyInterface { return &DummyStruct{} },
			expectedInTypes:  []reflect.Type{reflect.TypeOf((*DummyInterface)(nil)).Elem()},
			expectedOutTypes: []reflect.Type{reflect.TypeOf((*DummyInterface)(nil)).Elem()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inTypes, outTypes := ExtractFuncType(tc.fn)

			if !reflect.DeepEqual(inTypes, tc.expectedInTypes) {
				t.Errorf("Expected input types %v, but got %v", tc.expectedInTypes, inTypes)
			}

			if !reflect.DeepEqual(outTypes, tc.expectedOutTypes) {
				t.Errorf("Expected output types %v, but got %v", tc.expectedOutTypes, outTypes)
			}
		})
	}
}

func TestConvertReflectValuesToAnySlice(t *testing.T) {
	intPtr := new(int)

	testCases := []struct {
		name          string
		reflectValues []reflect.Value
		expectedAny   []any
	}{
		{
			name:          "[int, int] => [int, int]",
			reflectValues: []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)},
			expectedAny:   []any{1, 2},
		},
		{
			name:          "[string, string] => [string, string]",
			reflectValues: []reflect.Value{reflect.ValueOf("foo"), reflect.ValueOf("bar")},
			expectedAny:   []any{"foo", "bar"},
		},
		{
			name:          "[] => []",
			reflectValues: []reflect.Value{},
			expectedAny:   []any{},
		},
		{
			name:          "[struct] => [struct]",
			reflectValues: []reflect.Value{reflect.ValueOf(DummyStruct{})},
			expectedAny:   []any{DummyStruct{}},
		},
		{
			name:          "[[1,2,3]] => [[1,2,3]]",
			reflectValues: []reflect.Value{reflect.ValueOf([]int{1, 2, 3})},
			expectedAny:   []any{[]int{1, 2, 3}},
		},
		{
			name:          "[[[1,2,3].[4,5,6]]] => [[[1,2,3],[4,5,6]]]",
			reflectValues: []reflect.Value{reflect.ValueOf([][]int{{1, 2, 3}, {4, 5, 6}})},
			expectedAny:   []any{[][]int{{1, 2, 3}, {4, 5, 6}}},
		},
		{
			name:          "[*int] => [*int]",
			reflectValues: []reflect.Value{reflect.ValueOf(intPtr)},
			expectedAny:   []any{intPtr},
		},
		{
			name:          "[nil] => [nil]",
			reflectValues: []reflect.Value{reflect.ValueOf(nil)},
			expectedAny:   []any{nil},
		},
		{
			name:          "[any] => [any]",
			reflectValues: []reflect.Value{reflect.ValueOf(any(10))},
			expectedAny:   []any{any(10)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertReflectValuesToAnySlice(tc.reflectValues)

			if !reflect.DeepEqual(result, tc.expectedAny) {
				fmt.Printf("Test case '%s' failed\n", tc.name)
				fmt.Printf("Expected: %+v\n", tc.expectedAny)
				fmt.Printf("Actual:   %+v\n", result)
				t.FailNow()
			}
		})
	}
}

func TestConvertAnySliceToReflectValues(t *testing.T) {
	testCases := []struct {
		name       string
		anySlice   []any
		expectedRv []reflect.Value
	}{
		{
			name:     "[int, int] => [int, int]",
			anySlice: []any{1, 2},
			expectedRv: []reflect.Value{
				reflect.ValueOf(1),
				reflect.ValueOf(2),
			},
		},
		{
			name:     "[string, string] => [string, string]",
			anySlice: []any{"foo", "bar"},
			expectedRv: []reflect.Value{
				reflect.ValueOf("foo"),
				reflect.ValueOf("bar"),
			},
		},
		{
			name:       "[] => []",
			anySlice:   []any{},
			expectedRv: []reflect.Value{},
		},
		{
			name:       "[struct] => [struct]",
			anySlice:   []any{DummyStruct{}},
			expectedRv: []reflect.Value{reflect.ValueOf(DummyStruct{})},
		},
		{
			name:       "[[1,2,3]] => [[1,2,3]]",
			anySlice:   []any{[]int{1, 2, 3}},
			expectedRv: []reflect.Value{reflect.ValueOf([]int{1, 2, 3})},
		},
		{
			name:       "[[[1,2,3].[4,5,6]]] => [[[1,2,3],[4,5,6]]]",
			anySlice:   []any{[][]int{{1, 2, 3}, {4, 5, 6}}},
			expectedRv: []reflect.Value{reflect.ValueOf([][]int{{1, 2, 3}, {4, 5, 6}})},
		},
		{
			name:       "[*int] => [*int]",
			anySlice:   []any{new(int)},
			expectedRv: []reflect.Value{reflect.ValueOf(new(int)).Elem()},
		},
		{
			name:       "[any] => [any]",
			anySlice:   []any{any(10)},
			expectedRv: []reflect.Value{reflect.ValueOf(any(10))},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertAnySliceToReflectValues(tc.anySlice)

			if !reflect.DeepEqual(result, tc.expectedRv) {
				fmt.Printf("Test case '%s' failed\n", tc.name)
				fmt.Printf("Expected: %+v\n", tc.expectedRv)
				fmt.Printf("Actual:   %+v\n", result)
				t.FailNow()
			}
		})
	}
}

func TestConvertReflectValuesToReflectTypes(t *testing.T) {
	testCases := []struct {
		name            string
		reflectValues   []reflect.Value
		expectedRtSlice []reflect.Type
	}{
		{
			name:          "Test case 1: Convert int values to reflect.Type slice",
			reflectValues: []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)},
			expectedRtSlice: []reflect.Type{
				reflect.TypeOf(1),
				reflect.TypeOf(2),
			},
		},
		{
			name:          "Test case 2: Convert string values to reflect.Type slice",
			reflectValues: []reflect.Value{reflect.ValueOf("foo"), reflect.ValueOf("bar")},
			expectedRtSlice: []reflect.Type{
				reflect.TypeOf("foo"),
				reflect.TypeOf("bar"),
			},
		},
		{
			name:            "Test case 3: Convert empty values to empty reflect.Type slice",
			reflectValues:   []reflect.Value{},
			expectedRtSlice: []reflect.Type{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertReflectValuesToReflectTypes(tc.reflectValues)

			if !reflect.DeepEqual(result, tc.expectedRtSlice) {
				t.Errorf("Expected %v, but got %v", tc.expectedRtSlice, result)
			}
		})
	}
}
