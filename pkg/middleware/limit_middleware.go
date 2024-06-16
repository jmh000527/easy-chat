package middleware

import (
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

// LimitMiddleware 使用Redis实现的限流中间件结构体
// 它基于Token Bucket算法，通过Redis来存储和管理令牌桶的状态。
type LimitMiddleware struct {
	redisCfg            redis.RedisConf // Redis配置信息
	*limit.TokenLimiter                 // TokenLimiter实例，用于限流
}

// NewLimitMiddleware 创建一个新的LimitMiddleware实例
// 参数redisCfg: Redis的配置信息，用于连接Redis服务器。
// 返回值: 新创建的LimitMiddleware实例。
func NewLimitMiddleware(redisCfg redis.RedisConf) *LimitMiddleware {
	return &LimitMiddleware{
		redisCfg: redisCfg,
	}
}

// TokenLimitHandler 创建一个限流处理函数，并返回一个rest.Middleware
// 参数rate: 令牌桶的填充速率，即每秒允许的请求数。
// 参数burst: 令牌桶的容量，即允许瞬间爆发的请求数。
// 返回值: 一个rest.Middleware，用于限制HTTP请求的速率。
// 该方法通过TokenLimiter来实现限流，如果请求速率超过设定的阈值，则拒绝请求。
func (l *LimitMiddleware) TokenLimitHandler(rate, burst int) rest.Middleware {
	// 初始化TokenLimiter，使用Redis作为存储后端
	l.TokenLimiter = limit.NewTokenLimiter(rate, burst, redis.MustNewRedis(l.redisCfg), "REDIS_TOKEN_LIMIT_KEY")
	return func(next http.HandlerFunc) http.HandlerFunc {
		// 返回一个HTTP处理函数，用于实际处理请求
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查是否允许通过，即是否有足够的令牌
			if l.TokenLimiter.AllowCtx(r.Context()) {
				// 如果允许通过，则调用下一个处理函数
				next(w, r)
				return
			}
			// 如果不允许通过，则返回HTTP状态码429，表示请求过多
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}
}
