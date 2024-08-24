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

// FriendPutIn 处理添加好友请求
//
// 功能描述:
//   - 检查申请人和目标用户是否已是好友。
//   - 检查是否已有未处理的好友请求。
//   - 如果未找到好友关系或已有未处理的好友请求，则创建新的好友请求记录。
//
// 参数:
//   - in: `social.FriendPutInReq` 类型，包含用户ID、请求用户ID、请求消息和请求时间等信息。
//
// 返回值:
//   - `*social.FriendPutInResp`: 包含处理结果的响应对象。
//   - `error`: 如果发生错误，则返回相应的错误信息。
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// 查询申请人是否与目标用户是好友关系。
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	// 如果查询过程中出现错误且错误不是未找到，则返回数据库错误。
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err: %v req: %v ", err, in)
	}
	// 如果存在好友关系，直接返回空的响应。
	if friends != nil {
		return &social.FriendPutInResp{}, nil
	}

	// 检查是否已经有过好友申请，且申请未成功。
	friendRequests, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	// 如果查询过程中出现错误且错误不是未找到，则返回数据库错误。
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by rid and uid err: %v req: %v ", err, in)
	}
	// 如果存在未处理的好友请求，返回相应的错误。
	if friendRequests != nil {
		return &social.FriendPutInResp{}, err
	}

	// 创建新的好友请求记录。
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
		// 如果插入过程中出现错误，则返回数据库错误。
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v ", err, in)
	}

	// 返回空的响应，表示请求已成功添加。
	return &social.FriendPutInResp{}, nil
}
