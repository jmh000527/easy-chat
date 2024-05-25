package friend

import (
	"context"
	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"
	"easy-chat/apps/social/rpc/social"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendsOnlineResp, err error) {
	// 获取当前用户ID
	uid := ctxdata.GetUId(l.ctx)
	// 获取当前用户所有好友列表
	friendList, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return &types.FriendsOnlineResp{}, err
	}
	if len(friendList.List) == 0 {
		return &types.FriendsOnlineResp{}, nil
	}
	// 查询缓存中在线的用户
	uids := make([]string, 0, len(friendList.List))
	for _, friend := range friendList.List {
		uids = append(uids, friend.UserId)
	}
	onlines, err := l.svcCtx.Redis.Hgetall(constants.RedisOnlineUser)
	if err != nil {
		return nil, err
	}
	resOnlineList := make(map[string]bool, len(uids))
	for _, uid := range uids {
		if _, ok := onlines[uid]; !ok {
			resOnlineList[uid] = true
		} else {
			resOnlineList[uid] = false
		}
	}

	return &types.FriendsOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
