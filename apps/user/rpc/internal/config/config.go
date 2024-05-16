package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		Datasource string
	}

	Redisx redis.RedisConf

	Cache cache.CacheConf

	Jwt struct {
		AccessSecret string
		AccessExpire int64
	}
}
