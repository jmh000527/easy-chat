package svc

import (
	"easy-chat/apps/im/api/internal/config"
	"easy-chat/apps/im/rpc/imclient"
)

type ServiceContext struct {
	Config config.Config

	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
