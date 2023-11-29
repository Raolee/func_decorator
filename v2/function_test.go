package v2

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// TestFunctionCallBasicFunctionality - Function의 기본 기능 테스트
func TestFunctionCallBasicFunctionality(t *testing.T) {
	fn := func(ctx context.Context, req string) (string, error) {
		return req + "/processed", nil
	}

	f := Function[string, string]{
		fn: fn,
	}

	res, err := f.Call(context.Background(), "test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != "test/processed" {
		t.Errorf("Expected 'test/processed', got '%s'", res)
	}
}

// TestFunctionCallPanicHandling - 패닉 핸들링 테스트
func TestFunctionCallPanicHandling(t *testing.T) {
	fn := func(ctx context.Context, req string) (string, error) {
		panic("Test panic")
	}

	f := Function[string, string]{
		fn:            fn,
		panicHandling: true,
	}

	_, err := f.Call(context.Background(), "Test Request")
	if err == nil || err.Error() != "Test panic" {
		t.Errorf("Expected panic to be caught and returned as error")
	}
}

// TestFunctionCallWithRequestMiddleware - 요청 미들웨어 테스트
func TestFunctionCallWithRequestMiddleware(t *testing.T) {
	requestInterceptor := func(ctx context.Context, req string) (string, error) {
		return req + "/modified", nil
	}

	fn := func(ctx context.Context, req string) (string, error) {
		return req + "/processed", nil
	}

	f := Function[string, string]{
		requestMiddleware: []func(ctx context.Context, req string) (string, error){requestInterceptor},
		fn:                fn,
	}

	res, err := f.Call(context.Background(), "test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != "test/modified/processed" {
		t.Errorf("Expected 'test/modified/processed', got '%s'", res)
	}
}

// TestFunctionCallWithResponseMiddleware - 응답 미들웨어 테스트
func TestFunctionCallWithResponseMiddleware(t *testing.T) {
	responseInterceptor := func(ctx context.Context, res string) (string, error) {
		return res + "/modified", nil
	}

	fn := func(ctx context.Context, req string) (string, error) {
		return req + "/processed", nil
	}

	f := Function[string, string]{
		fn:                 fn,
		responseMiddleware: []func(ctx context.Context, res string) (string, error){responseInterceptor},
	}

	res, err := f.Call(context.Background(), "test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != "test/processed/modified" {
		t.Errorf("Expected 'test/processed/modified', got '%s'", res)
	}
}

// TestFunctionCallWithErrorMiddleware - 에러 미들웨어 테스트
func TestFunctionCallWithErrorMiddleware(t *testing.T) {
	errorInterceptor := func(ctx context.Context, req string, err error) error {
		return fmt.Errorf("Modified error: %w", err)
	}

	fn := func(ctx context.Context, req string) (string, error) {
		return "", errors.New("original error")
	}

	f := Function[string, string]{
		fn:                  fn,
		exceptionMiddleware: []func(ctx context.Context, req string, err error) error{errorInterceptor},
	}

	_, err := f.Call(context.Background(), "Test Request")
	if err == nil || err.Error() != "Modified error: original error" {
		t.Errorf("Expected 'Modified error: original error', got '%v'", err)
	}
}

func TestFunctionCallWithAllCaseUsingStruct(t *testing.T) {
	type TestRequest struct {
		A int
		B string
	}
	type TestResponse struct {
		A int
		B string
	}
	requestInterceptor := func(ctx context.Context, req *TestRequest) (*TestRequest, error) {
		req.A++
		req.B = req.B + "/modified"
		return req, nil
	}
	responseInterceptor := func(ctx context.Context, res *TestResponse) (*TestResponse, error) {
		res.A++
		res.B = res.B + "/modified"
		return res, nil
	}
	f := Function[*TestRequest, *TestResponse]{
		fn: func(ctx context.Context, req *TestRequest) (*TestResponse, error) {
			return &TestResponse{A: req.A + 1, B: req.B + "/processed"}, nil
		},
		requestMiddleware:  []func(ctx context.Context, req *TestRequest) (*TestRequest, error){requestInterceptor},
		responseMiddleware: []func(ctx context.Context, res *TestResponse) (*TestResponse, error){responseInterceptor},
	}

	res, err := f.Call(context.Background(), &TestRequest{A: 0, B: "test"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res.A != 3 {
		t.Errorf("Expected '3', got '%d'", res.A)
	}
	if res.B != "test/modified/processed/modified" {
		t.Errorf("Expected 'test/modified/processed/modified', got '%s'", res.B)
	}

}
