package svc

import (
	"easy-chat/apps/im/rpc/imclient"
	"easy-chat/apps/social/api/internal/config"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	socialclient.Social
	userclient.User
	imclient.Im
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Redis:  redis.MustNewRedis(c.Redisx),
	}
}
