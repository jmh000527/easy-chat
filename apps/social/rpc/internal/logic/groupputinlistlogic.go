package logic

import (
	"context"
	"easy-chat/pkg/xerr"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutinList 查询未处理的群组加入请求列表
//
// 功能描述:
//   - 根据群组ID获取所有未处理的群组加入请求
//   - 将这些请求转换为响应格式并返回
//
// 参数:
//   - in: `*social.GroupPutinListReq` 类型，包含群组ID，用于查询相关的群组请求
//
// 返回值:
//   - `*social.GroupPutinListResp`: 包含未处理的群组请求列表
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	// 使用群组ID获取所有未处理的群组请求
	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		// 如果获取未处理的请求列表失败，返回错误信息，并包装为数据库错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group req err: %v req: %v", err, in.GroupId)
	}

	// 创建一个用于存储响应的切片
	var respList []*social.GroupRequests

	// 将从数据库获取的请求列表复制到响应列表中
	err = copier.Copy(&respList, groupReqs)
	if err != nil {
		// 如果复制过程中发生错误，返回错误
		return nil, err
	}

	// 返回包含未处理群组请求列表的响应结构体
	return &social.GroupPutinListResp{
		List: respList, // 响应中的群组请求列表
	}, nil
}
