package logic

import (
	"context"
	"database/sql"
	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var (
	ErrGroupReqBeforePass   = xerr.NewMsg("请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsg("请求已拒绝")
)

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutInHandle 处理群组加入请求
//
// 功能描述:
//   - 通过群组请求ID查找群组请求记录
//   - 根据请求的处理结果更新群组请求状态
//   - 如果请求被批准，将申请者添加到群组成员列表中
//   - 使用事务确保操作的原子性
//
// 参数:
//   - in: `*social.GroupPutInHandleReq` 类型，包含群组请求处理的信息，包括请求ID、处理结果及操作用户ID
//
// 返回值:
//   - `*social.GroupPutInHandleResp`: 包含群组ID的响应对象（如果处理结果为批准）
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// 通过群组请求ID查找对应的群组请求记录
	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, int64(in.GroupReqId))
	if err != nil {
		// 若查找群组请求记录失败，返回错误，并附带错误信息和请求ID
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req err %v req %v", err, in.GroupReqId)
	}

	// 根据处理结果判断请求状态
	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		// 如果请求已经被批准，返回错误，表示之前已处理通过
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		// 如果请求已经被拒绝，返回错误，表示之前已处理拒绝
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}

	// 更新处理结果
	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}

	// 开启事务处理群组请求更新和群组成员插入操作
	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新群组请求记录
		if err := l.svcCtx.GroupRequestsModel.Update(l.ctx, session, groupReq); err != nil {
			// 若更新请求记录失败，返回错误，并附带错误信息和群组请求
			return errors.Wrapf(xerr.NewDBErr(), "update group req err: %v req: %v", err, groupReq)
		}

		// 如果处理结果不是通过，直接返回
		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}

		// 构建群组成员信息
		groupMember := &socialmodels.GroupMembers{
			GroupId:     groupReq.GroupId,                     // 群组ID
			UserId:      groupReq.ReqId,                       // 申请者用户ID
			RoleLevel:   int(constants.AtLargeGroupRoleLevel), // 群组角色等级
			OperatorUid: in.HandleUid,                         // 操作用户ID
		}

		// 将群组成员信息插入到群组成员表中
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			// 若插入失败，返回错误，并附带错误信息和群组成员信息
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err: %v req: %v", err, groupMember)
		}

		return nil
	})

	// 如果处理结果不是通过，直接返回错误
	if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
		return &social.GroupPutInHandleResp{}, err
	}

	// 返回群组加入处理响应（包含群组ID）和错误
	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId, // 处理的群组ID
	}, err
}
