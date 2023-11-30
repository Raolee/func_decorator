package v2

import (
	"reflect"
	"testing"
)

type Foo struct {
	A bool
	B int
}

type Bar interface {
	dummy()
}

func TestGetGenericType(t *testing.T) {

	intType := GetGenericType[int]()
	if intType.String() != "int" {
		t.Errorf("GetGenericType() = %s, want int", intType.String())
	}
	stringType := GetGenericType[string]()
	if stringType.String() != "string" {
		t.Errorf("GetGenericType() = %s, want string", stringType.String())
	}
	intPtrType := GetGenericType[*int]()
	if intPtrType.String() != "*int" {
		t.Errorf("GetGenericType() = %s, want *int", intPtrType.String())
	}
	intSliceType := GetGenericType[[]int]()
	if intSliceType.String() != "[]int" {
		t.Errorf("GetGenericType() = %s, want []int", intSliceType.String())
	}
	mapIntIntType := GetGenericType[map[int]int]()
	if mapIntIntType.String() != "map[int]int" {
		t.Errorf("GetGenericType() = %s, want map[int]int", mapIntIntType.String())
	}
	fooType := GetGenericType[Foo]()
	if fooType.String() != "v2.Foo" {
		t.Errorf("GetGenericType() = %s, want v2.Foo", fooType.String())
	}
	fooPtrType := GetGenericType[*Foo]()
	if fooPtrType.String() != "*v2.Foo" {
		t.Errorf("GetGenericType() = %s, want *v2.Foo", fooPtrType.String())
	}
	barType := GetGenericType[Bar]()
	if barType.String() != "v2.Bar" {
		t.Errorf("GetGenericType() = %s, want v2.Bar", barType.String())
	}
}

// TestEqualType - EqualType 함수를 테스트하기 위한 함수
func TestEqualType(t *testing.T) {
	tests := []struct {
		name   string
		t      reflect.Type
		v      any
		expect bool
	}{
		{
			name:   "IntType",
			t:      reflect.TypeOf(0),
			v:      42,
			expect: true,
		},
		{
			name:   "IntAndStringType",
			t:      reflect.TypeOf(0),
			v:      "hello",
			expect: false,
		},
		{
			name:   "SliceType",
			t:      reflect.TypeOf([]int{}),
			v:      []int{1, 2, 3},
			expect: true,
		},
		{
			name:   "DifferentSliceType",
			t:      reflect.TypeOf([]int{}),
			v:      []string{"a", "b", "c"},
			expect: false,
		},
		{
			name:   "PointerType",
			t:      reflect.TypeOf(&struct{}{}),
			v:      &struct{}{},
			expect: true,
		},
		{
			name:   "DifferentPointerType",
			t:      reflect.TypeOf(&struct{}{}),
			v:      &[]int{},
			expect: false,
		},
		{
			name:   "VoidType",
			t:      reflect.TypeOf(VoidType{}),
			v:      VoidType{},
			expect: true,
		},
		// 다른 타입에 대한 추가 테스트 케이스들...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualType(tt.t, tt.v); got != tt.expect {
				t.Errorf("EqualType() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func Test_zeroValue(t *testing.T) {

	zeroInt := zeroValue[int]()
	if zeroInt != 0 {
		t.Errorf("zeroValue[int]() want 0")
	}
	zeroString := zeroValue[string]()
	if zeroString != "" {
		t.Errorf("zeroValue[string]() want string.empty")
	}
	zeroIntPtr := zeroValue[*int]()
	if zeroIntPtr != nil {
		t.Errorf("zeroValue[*int]() want nil")
	}
	zeroIntSlice := zeroValue[[]int]()
	if zeroIntSlice != nil {
		t.Errorf("zeroValue[[]int]() want nil")
	}
	zeroMap := zeroValue[map[int]int]()
	if zeroMap != nil {
		t.Errorf("zeroValue[map[int]int]() want nil")
	}
	zeroFoo := zeroValue[Foo]()
	if !reflect.DeepEqual(zeroFoo, Foo{}) {
		t.Errorf("zeroValue[Foo]() = %v, want %v", zeroFoo, Foo{})
	}
	zeroFooPtr := zeroValue[*Foo]()
	if zeroFooPtr != nil {
		t.Errorf("zeroValue[*Foo]() want nil")
	}
	zeroBar := zeroValue[Bar]()
	if zeroBar != nil {
		t.Errorf("zeroValue[Bar]() want nil")
	}
}
