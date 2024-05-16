package logic

import (
	"context"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// todo: add your logic here and delete this line

	return &social.GroupCreateResp{}, nil
}
