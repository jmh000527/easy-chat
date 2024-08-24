package friend

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/userclient"
	"easy-chat/pkg/ctxdata"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendList 根据请求参数获取用户的好友列表
//
// 功能描述:
//   - 该方法根据当前上下文中的用户ID获取该用户的好友列表，并查询每个好友的详细信息，最终返回完整的好友列表。
//
// 参数:
//   - req: `*types.FriendListReq` 类型，包含请求参数的信息。
//   - 此处没有使用 req 直接从上下文获取用户ID，假设这是为了简化代码示例。
//
// 返回值:
//   - `*types.FriendListResp`: 返回包含好友列表的响应对象。
//   - `error`: 如果在过程中发生错误，返回相应的错误信息。
func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// 从上下文中获取当前用户ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用服务层接口获取用户的好友列表
	friends, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		// 如果获取好友列表失败，返回错误
		return nil, err
	}

	// 如果好友列表为空，返回一个空的好友列表响应
	if len(friends.List) == 0 {
		return &types.FriendListResp{}, nil
	}

	// 获取好友ID列表
	uids := make([]string, 0, len(friends.List))
	for _, i := range friends.List {
		uids = append(uids, i.FriendUid)
	}

	// 调用用户服务接口根据好友ID列表获取用户信息列表
	users, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		// 如果获取用户信息失败，返回一个空的好友列表响应
		return &types.FriendListResp{}, nil
	}

	// 将获取的用户信息存储在一个字典中，以便快速查找
	userRecords := make(map[string]*userclient.UserEntity, len(users.User))
	for i := range users.User {
		userRecords[users.User[i].Id] = users.User[i]
	}

	// 构造返回的好友列表
	respList := make([]*types.Friends, 0, len(friends.List))
	for _, v := range friends.List {
		friend := &types.Friends{
			Id:        v.Id,
			FriendUid: v.FriendUid,
		}

		// 如果找到好友的详细信息，则填充好友的昵称和头像
		if u, ok := userRecords[v.FriendUid]; ok {
			friend.Nickname = u.Nickname
			friend.Avatar = u.Avatar
		}
		respList = append(respList, friend)
	}

	// 返回包含好友列表的响应对象
	return &types.FriendListResp{
		List: respList,
	}, nil
}
