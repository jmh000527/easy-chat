package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Database string
	Redisx   redis.RedisConf
	UserRpc  zrpc.RpcClientConf

	JwtAuth struct {
		AccessSecret string
		//AccessExpire int64
	}
}
