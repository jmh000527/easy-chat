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

// SetRootToken 为系统设置根令牌。
// 这个方法用于生成一个根令牌，并存储到Redis中，以供系统后续使用。
// 根令牌具有最高的权限，通常用于系统初始化或关键操作。
//
// svc: 服务上下文，包含配置信息和Redis客户端等。
// 返回值: 错误信息，如果设置令牌过程中出现错误，则返回相应的错误。
func (svc *ServiceContext) SetRootToken() error {
	// 生成jwt令牌。
	// 使用服务配置中的访问密钥，结合当前时间戳和一个非常长的过期时间（999999999秒），
	// 以及系统根用户的UID，来生成一个具有最高权限的令牌。
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 999999999, constants.SystemRootUid)
	if err != nil {
		// 如果生成令牌过程中出现错误，返回错误信息。
		return err
	}

	// 将生成的根令牌存储到Redis中。
	return svc.Redis.Set(constants.RedisSystemRootToken, systemToken)
}
