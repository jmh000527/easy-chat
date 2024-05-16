package logic

import (
	"context"
	"database/sql"
	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"time"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// 查询申请人是否与目标是好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err: %v req: %v ", err, in)
	}
	// 是好友关系
	if friends != nil {
		return &social.FriendPutInResp{}, nil
	}

	// 是否已经有过申请，申请是不成功，没有完成
	friendRequests, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by rid and uid err: %v req: %v ", err, in)
	}
	// 已发送好友请求
	if friendRequests != nil {
		return &social.FriendPutInResp{}, err
	}

	// 创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v ", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
