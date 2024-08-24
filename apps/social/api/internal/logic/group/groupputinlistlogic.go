package group

import (
	"context"
	"easy-chat/apps/social/rpc/socialclient"
	"github.com/jinzhu/copier"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutInList 查询未处理的群组加入请求列表
//
// 功能描述:
//   - 根据提供的群组ID从服务层获取所有未处理的群组加入请求
//   - 将获取到的请求转换为响应格式并返回
//
// 参数:
//   - req: `*types.GroupPutInListRep` 类型，包含群组ID，用于查询相关的群组请求
//
// 返回值:
//   - `*types.GroupPutInListResp`: 包含未处理的群组请求列表
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListRep) (resp *types.GroupPutInListResp, err error) {
	// 调用服务层方法获取未处理的群组加入请求
	list, err := l.svcCtx.Social.GroupPutinList(l.ctx, &socialclient.GroupPutinListReq{
		GroupId: req.GroupId, // 提供的群组ID
	})
	if err != nil {
		// 如果在获取群组请求过程中发生错误，返回错误信息
		return nil, err
	}

	// 创建一个用于存储响应的切片
	var respList []*types.GroupRequests

	// 将获取到的群组请求列表复制到响应列表中
	err = copier.Copy(&respList, list.List)
	if err != nil {
		// 如果复制过程中发生错误，返回错误
		return nil, err
	}

	// 返回包含未处理群组请求列表的响应结构体
	return &types.GroupPutInListResp{
		List: respList, // 响应中的群组请求列表
	}, nil
}
