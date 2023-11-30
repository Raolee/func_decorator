package v2

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// TestDecoratedFunctionCallBasicFunctionality - Function의 기본 기능 테스트
func TestDecoratedFunctionCallBasicFunctionality(t *testing.T) {
	fn := func(ctx context.Context, req string) (string, error) {
		return req + "/processed", nil
	}

	f := DecoratedFunction[string, string]{
		fn: fn,
	}

	res, err := f.Call(context.Background(), "test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != "test/processed" {
		t.Errorf("Expected 'test/processed', got '%s'", res)
	}

	af, err := f.Any()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if af.GetRequestType().String() != "string" {
		t.Errorf("request ypte is 'string', got '%s'", af.GetRequestType().String())
	}
	if af.GetResponseType().String() != "string" {
		t.Errorf("response ypte is 'string', got '%s'", af.GetResponseType().String())
	}
	anyRes, err := af.Call(context.Background(), "test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if anyRes != "test/processed" {
		t.Errorf("Expected 'test/processed', got '%s'", res)
	}
}

// TestDecoratedFunctionCallPanicHandling - 패닉 핸들링 테스트
func TestDecoratedFunctionCallPanicHandling(t *testing.T) {
	fn := func(ctx context.Context, req string) (string, error) {
		panic("Test panic")
	}

	f := DecoratedFunction[string, string]{
		fn:            fn,
		panicHandling: true,
	}

	_, err := f.Call(context.Background(), "Test Request")
	if err == nil || err.Error() != "Test panic" {
		t.Errorf("Expected panic to be caught and returned as error")
	}
}

// TestDecoratedFunctionCallWithRequestMiddleware - 요청 미들웨어 테스트
func TestDecoratedFunctionCallWithRequestMiddleware(t *testing.T) {
	reqDecorator := func(ctx context.Context, req string) (string, error) {
		return req + "/modified", nil
	}

	fn := func(ctx context.Context, req string) (string, error) {
		return req + "/processed", nil
	}

	f := DecoratedFunction[string, string]{
		requestDecorators: []func(ctx context.Context, req string) (string, error){reqDecorator},
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

// TestDecoratedFunctionCallWithResponseMiddleware - 응답 미들웨어 테스트
func TestDecoratedFunctionCallWithResponseMiddleware(t *testing.T) {
	resDecorator := func(ctx context.Context, res string) (string, error) {
		return res + "/modified", nil
	}

	fn := func(ctx context.Context, req string) (string, error) {
		return req + "/processed", nil
	}

	f := DecoratedFunction[string, string]{
		fn:                 fn,
		responseDecorators: []func(ctx context.Context, res string) (string, error){resDecorator},
	}

	res, err := f.Call(context.Background(), "test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != "test/processed/modified" {
		t.Errorf("Expected 'test/processed/modified', got '%s'", res)
	}
}

// TestDecoratedFunctionCallWithErrorMiddleware - 에러 미들웨어 테스트
func TestDecoratedFunctionCallWithErrorMiddleware(t *testing.T) {
	exDecorator := func(ctx context.Context, req string, err error) error {
		return fmt.Errorf("modified error: %w", err)
	}

	fn := func(ctx context.Context, req string) (string, error) {
		return "", errors.New("original error")
	}

	f := DecoratedFunction[string, string]{
		fn:                  fn,
		exceptionDecorators: []func(ctx context.Context, req string, err error) error{exDecorator},
	}

	_, err := f.Call(context.Background(), "Test Request")
	if err == nil || err.Error() != "Modified error: original error" {
		t.Errorf("Expected 'Modified error: original error', got '%v'", err)
	}
}

func TestDecoratedFunctionCallWithAllCaseUsingStruct(t *testing.T) {
	type TestRequest struct {
		A int
		B string
	}
	type TestResponse struct {
		A int
		B string
	}
	reqDecorator := func(ctx context.Context, req *TestRequest) (*TestRequest, error) {
		req.A++
		req.B = req.B + "/modified"
		return req, nil
	}
	resDecorator := func(ctx context.Context, res *TestResponse) (*TestResponse, error) {
		res.A++
		res.B = res.B + "/modified"
		return res, nil
	}
	f := DecoratedFunction[*TestRequest, *TestResponse]{
		fn: func(ctx context.Context, req *TestRequest) (*TestResponse, error) {
			return &TestResponse{A: req.A + 1, B: req.B + "/processed"}, nil
		},
		requestDecorators:  []func(ctx context.Context, req *TestRequest) (*TestRequest, error){reqDecorator},
		responseDecorators: []func(ctx context.Context, res *TestResponse) (*TestResponse, error){resDecorator},
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
