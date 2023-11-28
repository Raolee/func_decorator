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
