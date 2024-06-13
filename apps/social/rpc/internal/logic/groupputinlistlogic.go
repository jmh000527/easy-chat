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

func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	// 通过GroupId获取未处理的群请求列表
	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		// 若发生错误，返回错误信息，并包装为数据库错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group req err: %v req: %v", err, in.GroupId)
	}

	// 创建一个用于存储响应的切片
	var respList []*social.GroupRequests
	// 将获取到的群请求列表复制到响应列表中
	err = copier.Copy(&respList, groupReqs)
	if err != nil {
		return nil, err
	}

	// 返回包含响应列表的GroupPutinListResp结构体和nil错误
	return &social.GroupPutinListResp{
		List: respList,
	}, nil
}
