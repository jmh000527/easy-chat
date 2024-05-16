package friend

import (
	"context"
	"easy-chat/apps/social/rpc/social"
	"easy-chat/pkg/ctxdata"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	uid := ctxdata.GetUId(l.ctx)
	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &social.FriendPutInReq{
		UserId:  uid,
		ReqUid:  req.UserId,
		ReqMsg:  req.ReqMsg,
		ReqTime: req.ReqTime,
	})

	return
}
