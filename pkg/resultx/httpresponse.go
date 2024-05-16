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

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		// 先设置默认错误码和错误描述
		errCode := xerr.ServerCommonError
		errMsg := xerr.ErrMsg(errCode)

		causeErr := errors.Cause(err)
		// 自定义的错误类型
		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
			errCode = e.Code
			errMsg = e.Msg
		} else {
			// gRPC的错误
			if gstatus, ok := status.FromError(causeErr); ok {
				errCode = int(gstatus.Code())
				errMsg = gstatus.Message()
			}
		}

		// 日志记录
		logx.WithContext(ctx).Errorf("【%s】 err: %v", name, err)

		return http.StatusBadRequest, Fail(errCode, errMsg)
	}
}
