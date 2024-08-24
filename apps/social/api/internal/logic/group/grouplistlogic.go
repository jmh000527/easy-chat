package group

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/ctxdata"
	"github.com/jinzhu/copier"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupList 获取用户所在的群组列表
//
// 功能描述:
//   - 从上下文中获取当前用户ID
//   - 调用服务层接口获取用户所在的群组列表
//   - 将获取的群组信息转换为响应格式并返回
//
// 参数:
//   - req: `*types.GroupListReq` 类型，包含请求获取群组列表的信息（当前未使用）
//
// 返回值:
//   - `*types.GroupListResp`: 包含用户所在群组列表的响应对象
//   - `error`: 如果获取群组列表过程中发生错误，则返回相应的错误信息
func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	// 从上下文中获取当前用户ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用服务层接口获取用户所在的群组列表
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid, // 当前用户ID
	})
	if err != nil {
		// 如果获取群组列表失败，返回错误信息
		return nil, err
	}

	// 将获取到的群组信息转换为响应格式
	var respList []*types.Groups
	copier.Copy(&respList, list.List)

	// 返回群组列表响应，包括群组列表和nil错误信息
	return &types.GroupListResp{
		List: respList, // 用户所在的群组列表
	}, nil
}
