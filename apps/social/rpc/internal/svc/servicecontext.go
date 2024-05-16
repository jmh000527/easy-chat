package svc

import (
	"easy-chat/apps/social/rpc/internal/config"
	"easy-chat/apps/social/socialmodels"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	socialmodels.FriendsModel
	socialmodels.FriendRequestsModel
	socialmodels.GroupsModel
	socialmodels.GroupRequestsModel
	socialmodels.GroupMembersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	//// 设置 SQL 模式
	//_, err := sqlConn.ExecCtx(context.Background(), "SET @@sql_mode='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'")
	//if err != nil {
	//	panic(err)
	//}

	return &ServiceContext{
		Config: c,

		FriendsModel:        socialmodels.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: socialmodels.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:         socialmodels.NewGroupsModel(sqlConn, c.Cache),
		GroupRequestsModel:  socialmodels.NewGroupRequestsModel(sqlConn, c.Cache),
		GroupMembersModel:   socialmodels.NewGroupMembersModel(sqlConn, c.Cache),
	}
}
