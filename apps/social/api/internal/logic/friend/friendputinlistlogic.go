package friend

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/ctxdata"
	"github.com/jinzhu/copier"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendPutInList 获取用户的好友申请列表
//
// 功能描述:
//   - 从上下文中获取当前用户ID
//   - 调用服务层接口获取当前用户的好友申请列表
//   - 将获取的好友申请列表转换为响应对象并返回
//
// 参数:
//   - req: `*types.FriendPutInListReq` 类型，包含获取好友申请列表的请求信息
//   - `UserId`: 当前用户ID，表示请求获取好友申请列表的用户
//
// 返回值:
//   - `*types.FriendPutInListResp`: 响应对象，包含当前用户的好友申请列表
//   - `List`: 好友申请列表，包含所有待处理的好友申请记录
//   - `error`: 如果在获取好友申请列表过程中发生错误，则返回相应的错误信息
func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	// 从上下文中获取当前用户ID
	userId := ctxdata.GetUId(l.ctx)

	// 调用服务层接口获取当前用户的好友申请列表
	list, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialclient.FriendPutInListReq{
		UserId: userId, // 当前用户ID
	})
	if err != nil {
		// 如果获取好友申请列表失败，返回错误信息
		return nil, err
	}

	// 定义响应列表，用于存储转换后的好友申请记录
	var respList []*types.FriendRequests

	// 将获取到的好友申请列表复制到响应列表中
	err = copier.Copy(&respList, list.List)
	if err != nil {
		// 如果复制过程中出现错误，返回错误信息
		return nil, err
	}

	// 返回好友申请列表响应对象
	return &types.FriendPutInListResp{
		List: respList, // 好友申请列表
	}, nil
}
