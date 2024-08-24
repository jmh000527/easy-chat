package group

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/constants"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupUserOnline 查询群组中所有用户的在线状态
//
// 功能描述:
//   - 获取指定群组的所有成员
//   - 检查这些成员是否在线，依据缓存中的在线用户信息
//   - 返回每个成员的在线状态
//
// 参数:
//   - req: `*types.GroupUserOnlineReq` 类型，包含群组ID，用于指定要查询的群组
//
// 返回值:
//   - `*types.GroupUserOnlineResp`: 包含在线用户的状态映射
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// 获取当前群组的所有成员信息
	groupUsers, err := l.svcCtx.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId, // 群组ID
	})
	if err != nil {
		// 如果获取群组成员信息失败，则返回空响应和错误
		return &types.GroupUserOnlineResp{}, err
	}

	// 如果群组没有成员，则返回空的在线用户状态响应
	if len(groupUsers.List) == 0 {
		return &types.GroupUserOnlineResp{}, nil
	}

	// 提取群组成员的UID列表
	uids := make([]string, 0, len(groupUsers.List))
	for _, groupUser := range groupUsers.List {
		uids = append(uids, groupUser.UserId)
	}

	// 查询缓存中所有在线用户的状态
	onlines, err := l.svcCtx.Redis.Hgetall(constants.RedisOnlineUser)
	if err != nil {
		// 如果查询缓存失败，则返回空响应和错误
		return nil, err
	}

	// 创建一个映射，用于存储每个用户的在线状态
	resOnLineList := make(map[string]bool, len(uids))
	for _, s := range uids {
		// 如果用户ID在缓存中，则表示该用户在线，否则为离线
		if _, ok := onlines[s]; ok {
			resOnLineList[s] = true
		} else {
			resOnLineList[s] = false
		}
	}

	// 返回群组用户在线状态的响应
	return &types.GroupUserOnlineResp{
		OnlineList: resOnLineList, // 在线用户状态映射
	}, nil
}
