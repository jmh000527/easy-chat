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

// NewServiceContext 初始化一个新的 ServiceContext 对象。
//
// 参数:
//   - c: 配置对象，包含数据库、Redis 等配置信息。
//
// 返回值:
//   - *ServiceContext: 返回一个包含配置信息、Redis 客户端和用户模型的 ServiceContext 实例。
func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.Datasource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

// SetRootToken 为系统设置根令牌。
//
// 功能描述:
//   - 生成一个系统根令牌，并将其存储到 Redis 中。
//   - 该根令牌具有最高权限，通常用于系统初始化或关键操作。
//
// 参数:
//   - svc: ServiceContext 对象，包含系统的配置信息和 Redis 客户端。
//
// 返回值:
//   - error: 错误信息。如果在生成或存储令牌的过程中出现错误，将返回相应的错误信息。
func (svc *ServiceContext) SetRootToken() error {
	// 生成 JWT 令牌。
	// 使用服务配置中的访问密钥，结合当前时间戳和一个非常长的过期时间（999999999 秒），
	// 以及系统根用户的 UID，来生成一个具有最高权限的令牌。
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 999999999, constants.SystemRootUid)
	if err != nil {
		// 如果生成令牌过程中出现错误，返回错误信息。
		return err
	}

	// 将生成的根令牌存储到 Redis 中。
	return svc.Redis.Set(constants.RedisSystemRootToken, systemToken)
}
