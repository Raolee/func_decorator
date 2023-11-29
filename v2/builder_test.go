package v2

import (
	"context"
	"errors"
	"testing"
)

// 테스트용 변수 선언
var fn = func(ctx context.Context, req string) (string, error) {
	return req + "/processed", nil
}
var requestInterceptor = func(ctx context.Context, req string) (string, error) {
	return req + "/req_processed", nil
}
var responseInterceptor = func(ctx context.Context, res string) (string, error) {
	return res + "/res_processed", nil
}
var exceptionInterceptor = func(ctx context.Context, req string, err error) error {
	return errors.Join(err, errors.New("custom raol error"))
}

// TestFunctionBuilder - FunctionBuilder의 기능을 테스트합니다.
func TestFunctionBuilder(t *testing.T) {
	builder := NewFunctionBuilder[string, string]()

	// 기본 함수 설정
	builder.Func(fn)
	// 요청 미들웨어 설정
	builder.RequestMiddleware(requestInterceptor, requestInterceptor)
	// 응답 미들웨어 설정
	builder.ResponseMiddleware(responseInterceptor, responseInterceptor)
	// 예외 미들웨어 설정
	builder.ExceptionMiddleware(exceptionInterceptor)
	// 패닉 핸들링 설정
	builder.PanicHandling(true)
	// 빌더 패턴을 사용하여 Function 객체 생성
	function := builder.Build()

	// 생성된 Function 객체를 검증
	if function.fn == nil {
		t.Errorf("Function fn was not set properly")
	}
	if len(function.requestMiddleware) != 2 {
		t.Errorf("Function requestMiddleware was not set properly")
	}
	if len(function.responseMiddleware) != 2 {
		t.Errorf("Function responseMiddleware was not set properly")
	}
	if len(function.exceptionMiddleware) != 1 {
		t.Errorf("Function exceptionMiddleware was not set properly")
	}
	if !function.panicHandling {
		t.Errorf("Function panicHandling was not set to true")
	}

	res, err := function.Call(context.Background(), "testbuilder")
	if err != nil {
		t.Errorf("Unexceped error: %v", err)
	}
	if res != "testbuilder/req_processed/req_processed/processed/res_processed/res_processed" {
		t.Errorf("Expected 'testbuilder/req_processed/req_processed/processed/res_processed/res_processed', got %s", res)
	}
}
