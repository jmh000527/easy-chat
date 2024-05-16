package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	ListenOn string

	MsgChatTransfer kq.KqConf

	Redisx redis.RedisConf

	Mongo struct {
		Url string
		Db  string
	}

	SocialRpc zrpc.RpcClientConf

	Ws struct {
		Host string
	}
}
