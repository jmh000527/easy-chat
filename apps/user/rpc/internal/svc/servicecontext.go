package svc

import (
	"easy-chat/apps/user/models"
	"easy-chat/apps/user/rpc/internal/config"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.Datasource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

func (svc *ServiceContext) SetRootToken() error {
	// 生成jwt
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 999999999, constants.SystemRootUid)
	if err != nil {
		return err
	}
	// 写入到redis
	return svc.Redis.Set(constants.RedisSystemRootToken, systemToken)
}
