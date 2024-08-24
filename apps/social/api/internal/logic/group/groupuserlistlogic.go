package group

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/userclient"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupUserList 处理查询群组成员及其用户信息的请求
//
// 功能描述:
//   - 该方法查询指定群组的所有成员，并获取这些成员的用户信息。
//   - 将群组成员及其对应的用户信息构造为响应返回。
//
// 参数:
//   - req: `*types.GroupUserListReq` 类型，包含查询所需的群组ID。
//
// 返回值:
//   - `*types.GroupUserListResp`: 包含群组成员及其用户信息的响应。
//   - `error`: 处理过程中发生的错误。
func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {
	// 获取群组成员列表
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		// 查询群组成员失败，返回错误
		return nil, err
	}

	// 提取成员用户ID以便批量查询用户信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}

	// 获取用户信息
	userList, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		// 查询用户信息失败，返回错误
		return nil, err
	}

	// 构造用户信息映射
	userRecords := make(map[string]*userclient.UserEntity, len(userList.User))
	for i := range userList.User {
		userRecords[userList.User[i].Id] = userList.User[i]
	}

	// 构造返回结果列表
	respList := make([]*types.GroupMembers, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		// 创建群组成员对象
		member := &types.GroupMembers{
			Id:        int64(v.Id),
			GroupId:   v.GroupId,
			UserId:    v.UserId,
			RoleLevel: int(v.RoleLevel),
		}
		// 赋值用户信息
		if u, ok := userRecords[v.UserId]; ok {
			member.Nickname = u.Nickname
			member.UserAvatarUrl = u.Avatar
		}
		respList = append(respList, member)
	}

	// 返回包含群组成员及其用户信息的响应
	return &types.GroupUserListResp{List: respList}, nil
}
