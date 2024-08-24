package resultx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"

	"easy-chat/pkg/xerr"
)

// Response 结构体定义了一个标准的HTTP响应格式，包含状态码、消息和数据。
type Response struct {
	Code int         `json:"code"` // HTTP状态码
	Msg  string      `json:"msg"`  // 响应消息或错误描述
	Data interface{} `json:"data"` // 响应数据，类型为任意类型
}

// success 返回一个成功的响应，包含指定的数据。
//
// 参数:
//   - data: 响应数据，可以是任意类型。
//
// 返回值:
//   - *Response: 包含成功状态码(200)、空消息和传入数据的Response对象。
func success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

// fail 返回一个失败的响应，包含指定的错误码和错误描述。
//
// 参数:
//   - code: 错误码，通常为HTTP状态码。
//   - err: 错误描述信息。
//
// 返回值:
//   - *Response: 包含指定错误码、错误消息和空数据的Response对象。
func fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

// OkHandler 是一个处理成功情况的处理函数，直接返回成功的响应。
//
// 参数:
//   - ctx: 上下文对象（未使用，但为了保持一致性）。
//   - v: 任意类型的数据，将包含在成功响应中。
//
// 返回值:
//   - any: 一个包含成功数据的Response对象。
func OkHandler(_ context.Context, v interface{}) any {
	return success(v)
}

// ErrHandler 返回一个用于处理失败情况的处理函数。
//
// 参数:
//   - name: 错误发生的上下文名称，用于日志记录。
//
// 返回值:
//   - func(ctx context.Context, err error) (int, any): 一个处理错误的函数，该函数返回HTTP状态码和包含错误信息的Response对象。
func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		// 先设置默认错误码和错误描述
		errCode := xerr.ServerCommonError
		errMsg := xerr.ErrMsg(errCode)

		// 获取错误的根本原因
		causeErr := errors.Cause(err)
		// 处理自定义错误类型
		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
			errCode = e.Code
			errMsg = e.Msg
		} else {
			// 处理gRPC错误
			if gstatus, ok := status.FromError(causeErr); ok {
				errCode = int(gstatus.Code())
				errMsg = gstatus.Message()
			}
		}

		// 记录错误日志
		logx.WithContext(ctx).Errorf("【%s】 err: %v", name, err)

		// 返回HTTP状态码和错误响应
		return http.StatusBadRequest, fail(errCode, errMsg)
	}
}
