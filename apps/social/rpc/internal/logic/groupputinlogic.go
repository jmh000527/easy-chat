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

// GroupPutin 处理群组加入请求
//
// 功能描述:
//   - 根据请求的类型（用户申请、群成员邀请、管理员/创建者邀请）处理群组加入请求。
//   - 如果群组不需要验证，直接将用户加入群组；否则，根据请求的来源和群组角色处理请求。
//
// 参数:
//   - in: `*social.GroupPutinReq` 类型，包含用户ID、群组ID、请求消息、请求时间、加入来源等信息。
//
// 返回值:
//   - `*social.GroupPutinResp`: 处理结果的响应，可能包含群组ID。
//   - `error`: 处理过程中发生的错误。
func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	// 定义局部变量，用于存储查询结果和错误信息
	var (
		inviteGroupMember *socialmodels.GroupMembers // 邀请者的群组成员信息
		userGroupMember   *socialmodels.GroupMembers // 用户的群组成员信息
		groupInfo         *socialmodels.Groups       // 群组信息
		err               error                      // 错误信息
	)

	// 查询用户是否已经是群组成员
	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && err != socialmodels.ErrNotFound {
		// 如果查询出错且错误不是“未找到”，则返回错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by group id and req id err %v, req %v, %v", err, in.GroupId, in.ReqId)
	}
	if userGroupMember != nil {
		// 如果用户已经是群组成员，则直接返回成功响应
		return &social.GroupPutinResp{}, nil
	}

	// 查询用户是否有加入群组的请求
	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != socialmodels.ErrNotFound {
		// 如果查询出错且错误不是“未找到”，则返回错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req by group id and req id err %v, req %v, %v", err, in.GroupId, in.ReqId)
	}
	if groupReq != nil {
		// 如果用户已经有加入群组的请求，则直接返回成功响应
		return &social.GroupPutinResp{}, nil
	}

	// 初始化群组请求数据
	groupReq = &socialmodels.GroupRequests{
		ReqId:      in.ReqId,
		GroupId:    in.GroupId,
		ReqMsg:     sql.NullString{String: in.ReqMsg, Valid: true},
		ReqTime:    sql.NullTime{Time: time.Unix(in.ReqTime, 0), Valid: true},
		JoinSource: sql.NullInt64{Int64: int64(in.JoinSource), Valid: true},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}

	// 定义闭包函数，用于创建群组成员
	createGroupMember := func() {
		if err != nil {
			return
		}
		err = l.createGroupMember(in)
	}

	// 查询群组信息
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		// 如果查询出错，则返回错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by group id err %v, req %v", err, in.GroupId)
	}

	// 验证是否要验证加入请求
	if !groupInfo.IsVerify {
		// 如果不需要验证，则直接创建群组成员，并设置请求处理结果为通过
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}

		// 创建并返回群组请求成功响应
		return l.createGroupReq(groupReq, true)
	}

	// 验证进群方式
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		// 如果是用户主动申请加入，则创建群组请求，并返回成功响应（待审核状态）
		return l.createGroupReq(groupReq, false)
	}

	// 查询邀请者是否是群组的管理者或创建者
	inviteGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
	if err != nil {
		// 如果查询出错，则返回错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by group id and user id err %v, req %v", in.InviterUid, in.GroupId)
	}

	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.CreatorGroupRoleLevel || constants.
		GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		// 如果邀请者是群组的管理者或创建者，则执行以下逻辑：
		// 1. 使用defer确保在函数返回前执行createGroupMember，即创建群组成员
		// 2. 设置请求处理结果为通过，并指定处理请求的用户为邀请者
		// 3. 调用createGroupReq创建群组请求并返回成功响应（传入的第二个参数为true）
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.createGroupReq(groupReq, true)
	}

	// 如果邀请者既不是管理者也不是创建者，则调用createGroupReq创建群组请求并返回成功响应（传入的第二个参数为false）
	return l.createGroupReq(groupReq, false)
}

// createGroupReq 插入群组请求并返回响应
//
// 功能描述:
//   - 将群组请求插入数据库，并根据请求是否通过返回相应的响应。
//
// 参数:
//   - groupReq: `*socialmodels.GroupRequests` 类型，包含群组请求的信息。
//   - isPass: `bool` 类型，指示请求是否通过（true表示通过，false表示待审核）
//
// 返回值:
//   - `*social.GroupPutinResp`: 处理结果的响应，可能包含群组ID。
//   - `error`: 处理过程中发生的错误。
func (l *GroupPutinLogic) createGroupReq(groupReq *socialmodels.GroupRequests, isPass bool) (*social.GroupPutinResp, error) {
	// 将群组请求插入数据库
	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		// 如果插入失败，则返回错误，并包装为数据库错误类型，同时附带错误信息和请求详情
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group req err: %v req: %v", err, groupReq)
	}

	// 如果isPass为true，表示请求通过，直接返回带有GroupId的响应
	if isPass {
		return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
	}

	// 否则，返回空的GroupPutinResp响应和nil错误（表示没有错误发生）
	return &social.GroupPutinResp{}, nil
}

// createGroupMember 创建群组成员
//
// 功能描述:
//   - 将用户作为群组成员插入数据库。
//
// 参数:
//   - in: `*social.GroupPutinReq` 类型，包含用户ID和群组ID等信息。
//
// 返回值:
//   - `error`: 处理过程中发生的错误。
func (l *GroupPutinLogic) createGroupMember(in *social.GroupPutinReq) error {
	// 构建群组成员结构体，并设置相关属性
	groupMember := &socialmodels.GroupMembers{
		GroupId:     in.GroupId,
		UserId:      in.ReqId,
		RoleLevel:   int(constants.AtLargeGroupRoleLevel),
		OperatorUid: in.InviterUid,
	}
	// 将群组成员插入数据库（注意这里传入了nil作为第二个参数，可能需要根据实际情况调整）
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, nil, groupMember)
	if err != nil {
		// 如果插入失败，则返回错误，并包装为数据库错误类型，同时附带错误信息和群组成员详情
		return errors.Wrapf(xerr.NewDBErr(), "insert group member err: %v req: %v", err, groupMember)
	}

	// 如果没有错误发生，则返回nil表示操作成功
	return nil
}
