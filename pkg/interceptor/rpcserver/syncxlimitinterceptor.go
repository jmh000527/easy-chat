package rpcserver

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// SyncXLimitInterceptor 创建一个 unary server interceptor，用于限制同时处理的请求数量。
// maxCount 参数指定最大并发请求量。
func SyncXLimitInterceptor(maxCount int) grpc.UnaryServerInterceptor {
	// 初始化一个并发控制对象，用于限制并发请求的数量。
	l := syncx.NewLimit(maxCount)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 尝试获取一个并发控制令牌，如果成功，则表示可以处理请求。
		if l.TryBorrow() {
			// 请求处理完成后，必须归还令牌。
			defer func() {
				// 如果归还令牌时发生错误，记录错误信息。
				if err := l.Return(); err != nil {
					logx.Errorf(err.Error())
				}
			}()
			// 处理请求。
			return handler(ctx, req)
		} else {
			// 如果无法获取并发控制令牌，则表示并发请求已达到上限，拒绝处理请求。
			logx.Errorf("concurrent connections exceeded %d, rejected with code %d", maxCount, http.StatusTooManyRequests)
			// 返回一个资源耗尽的错误，告知客户端并发请求已超过限制。
			return nil, status.Error(codes.ResourceExhausted, "concurrent connections exceeded limit")
		}
	}
}
