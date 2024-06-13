package svc

import (
	"easy-chat/apps/im/rpc/imclient"
	"easy-chat/apps/social/api/internal/config"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/userclient"
	"easy-chat/pkg/interceptor"
	"easy-chat/pkg/middleware"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// ServiceContext 包含所有服务客户端和配置
type ServiceContext struct {
	Config config.Config // 全局配置

	IdempotenceMiddleware rest.Middleware
	socialclient.Social   // 社交服务客户端
	userclient.User       // 用户服务客户端
	imclient.Im           // 即时通讯服务客户端
	*redis.Redis          // Redis 客户端
}

// retryPolicy 定义了 gRPC 客户端的重试策略
var retryPolicy = `{
	"methodConfig": [{
		"name": [{
			"service": "social.social"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

// NewServiceContext 创建一个新的 ServiceContext 实例
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient)),
		),
		User: userclient.NewUser(
			zrpc.MustNewClient(c.UserRpc),
		),
		Im: imclient.NewIm(
			zrpc.MustNewClient(c.ImRpc),
		),
		Redis:                 redis.MustNewRedis(c.Redisx),
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
	}
}
