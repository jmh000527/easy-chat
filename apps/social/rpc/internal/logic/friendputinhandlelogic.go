package logic

import (
	"context"
	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsg("好友申请并已经通过")
	ErrFriendReqBeforeRefuse = xerr.NewMsg("好友申请已经被拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 获取好友申请记录
	friendRequest, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, int64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqid err: %v req: %v ", err,
			in.FriendReqId)
	}

	// 验证是否有处理
	switch constants.HandlerResult(friendRequest.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}
	friendRequest.HandleResult.Int64 = int64(in.HandleResult)

	// 修改申请结果 -> 通过【建立两条好友关系记录】 -> 事务
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, friendRequest)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err: %v, req: %v", err, friendRequest)
		}

		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		friends := []*socialmodels.Friends{
			{
				UserId:    friendRequest.UserId,
				FriendUid: friendRequest.ReqUid,
			}, {
				UserId:    friendRequest.ReqUid,
				FriendUid: friendRequest.UserId,
			},
		}

		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "friends inserts err %v, req %v", err, friends)
		}

		return nil
	})

	return &social.FriendPutInHandleResp{}, nil
}
