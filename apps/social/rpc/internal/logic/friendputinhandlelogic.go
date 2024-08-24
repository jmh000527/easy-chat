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

// FriendPutInHandle 处理好友申请
//
// 功能描述:
//   - 根据好友申请ID获取申请记录。
//   - 验证申请是否已经被处理。
//   - 根据处理结果更新申请记录，并在通过时建立两条好友关系记录。
//   - 所有操作在数据库事务中执行，以确保数据一致性。
//
// 参数:
//   - in: `social.FriendPutInHandleReq` 类型，包含好友申请ID和处理结果。
//
// 返回值:
//   - `*social.FriendPutInHandleResp`: 处理结果的响应对象。
//   - `error`: 如果处理过程中发生错误，则返回相应的错误信息。
func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 获取指定好友申请记录
	friendRequest, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, int64(in.FriendReqId))
	if err != nil {
		// 如果查询过程中发生错误，返回数据库错误，并包装详细的错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqid err: %v req: %v ", err, in.FriendReqId)
	}

	// 验证申请是否已经被处理
	switch constants.HandlerResult(friendRequest.HandleResult.Int64) {
	case constants.PassHandlerResult:
		// 申请已经通过，返回相应的错误
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		// 申请已经拒绝，返回相应的错误
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}

	// 请求未被处理
	// 处理申请，填写客户端发送来的期望结果
	friendRequest.HandleResult.Int64 = int64(in.HandleResult)

	// 在事务中执行以下操作：更新申请记录和建立好友关系
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新好友申请记录
		err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, friendRequest)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err: %v, req: %v", err, friendRequest)
		}

		// 如果处理结果为拒绝，直接返回
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		// 如果处理结果为通过，建立两条好友关系记录
		friends := []*socialmodels.Friends{
			{
				UserId:    friendRequest.UserId,
				FriendUid: friendRequest.ReqUid,
			}, {
				UserId:    friendRequest.ReqUid,
				FriendUid: friendRequest.UserId,
			},
		}

		// 插入好友关系记录
		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "friends inserts err %v, req %v", err, friends)
		}

		return nil
	})

	// 返回处理成功的响应对象
	return &social.FriendPutInHandleResp{}, nil
}
