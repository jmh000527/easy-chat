package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogInterceptor 是 gRPC 的拦截器，用于记录 gRPC 服务端处理请求的日志，并将错误信息包装为 gRPC 的状态错误。
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 调用下一个处理器处理请求
	resp, err := handler(ctx, req)
	if err == nil {
		// 如果没有错误，直接返回响应
		return resp, nil
	}

	// 如果有错误，记录错误日志
	logx.WithContext(ctx).Errorf("【RPC SRV ERR】 %v", err)

	// 将错误信息包装为 gRPC 的状态错误
	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*zerr.CodeMsg); ok {
		err = status.Error(codes.Code(e.Code), e.Msg)
	}

	return resp, err
}
