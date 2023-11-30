package v2

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

// TestAnyFunction - AnyFunction 인터페이스와 그 구현체에 대한 테스트
func TestAnyFunction(t *testing.T) {
	// 테스트를 위한 더미 함수들
	reqEchoFunc := func(ctx context.Context, req any) (res any, err error) {
		return req, nil
	}

	stringEchoFunc := func(ctx context.Context, req any) (res any, err error) {
		return req.(string), nil
	}

	errorFunc := func(ctx context.Context, req any) (res any, err error) {
		return nil, fmt.Errorf("error")
	}
	voidFunc := func(ctx context.Context, _ any) (_ any, err error) {
		return VoidType{}, nil
	}

	tests := []struct {
		name       string
		reqType    reflect.Type
		resType    reflect.Type
		fn         func(ctx context.Context, req any) (res any, err error)
		signature  Signature
		request    any
		wantResult any
		wantErr    bool
	}{
		{
			name:       "ExistsRequestExistsResponseAny",
			reqType:    reflect.TypeOf([]int{0, 1}),
			resType:    reflect.TypeOf([]int{0, 1}),
			fn:         reqEchoFunc,
			signature:  ExistsRequestExistsResponse,
			request:    []int{0, 1},
			wantResult: []int{0, 1},
			wantErr:    false,
		},
		{
			name:       "ExistsRequestExistsResponseString",
			reqType:    reflect.TypeOf(""),
			resType:    reflect.TypeOf(""),
			fn:         stringEchoFunc,
			signature:  ExistsRequestExistsResponse,
			request:    "test",
			wantResult: "test",
			wantErr:    false,
		},
		{
			name:       "ReqTypeNil",
			reqType:    nil,
			resType:    reflect.TypeOf(0),
			fn:         reqEchoFunc,
			signature:  ExistsRequestExistsResponse,
			request:    123,
			wantResult: nil,
			wantErr:    true,
		},
		{
			name:       "ResTypeNil",
			reqType:    reflect.TypeOf(0),
			resType:    nil,
			fn:         reqEchoFunc,
			signature:  ExistsRequestExistsResponse,
			request:    123,
			wantResult: nil,
			wantErr:    true,
		},
		{
			name:       "FunctionNil",
			reqType:    reflect.TypeOf(0),
			resType:    reflect.TypeOf(0),
			fn:         nil,
			signature:  ExistsRequestExistsResponse,
			request:    123,
			wantResult: nil,
			wantErr:    true,
		},
		{
			name:       "FunctionReturnsError",
			reqType:    reflect.TypeOf(0),
			resType:    nil,
			fn:         errorFunc,
			signature:  ExistsRequestNoResponse,
			request:    123,
			wantResult: nil,
			wantErr:    true,
		},
		{
			name:       "VoidFunction",
			reqType:    reflect.TypeOf(VoidType{}),
			resType:    reflect.TypeOf(VoidType{}),
			fn:         voidFunc,
			signature:  NoRequestNoResponse,
			request:    VoidType{},
			wantResult: VoidType{},
			wantErr:    false,
		},
		// ... 필요에 따라 더 많은 케이스들 추가 ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			af, err := NewAnyFunction(tc.reqType, tc.resType, tc.fn)

			if (err != nil) == tc.wantErr {
				return // 생성자 에러 에상이 된 것은 성공
			}

			gotResult, err := af.Call(context.Background(), tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tc.wantResult) {
				t.Errorf("Call() gotResult = %v, want %v", gotResult, tc.wantResult)
			}
			if af.GetSignature() != tc.signature {
				t.Errorf("GetSignature() = %v, want %v", af.GetSignature(), tc.signature)
			}
			if af.GetRequestType() != tc.reqType {
				t.Errorf("GetRequestType() = %v, want %v", af.GetRequestType(), tc.reqType)
			}
			if af.GetResponseType() != tc.resType {
				t.Errorf("GetResponseType() = %v, want %v", af.GetResponseType(), tc.resType)
			}
		})
	}
}
