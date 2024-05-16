package logic

import (
	"context"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	// todo: add your logic here and delete this line

	return &social.GroupPutinResp{}, nil
}
