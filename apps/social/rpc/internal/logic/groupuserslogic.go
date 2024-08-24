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

type GroupUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUsersLogic {
	return &GroupUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupUsers 处理查询群组成员的请求
//
// 功能描述:
//   - 该方法根据提供的群组ID查询群组成员列表。
//   - 如果查询成功，将成员列表返回给调用者。
//   - 如果查询过程中发生错误，返回错误信息。
//
// 参数:
//   - in: `*social.GroupUsersReq` 类型，包含查询群组成员所需的群组ID。
//
// 返回值:
//   - `*social.GroupUsersResp`: 包含群组成员列表的响应。
//   - `error`: 处理过程中发生的错误。
func (l *GroupUsersLogic) GroupUsers(in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	// 通过群组ID查询群组成员列表
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByGroupId(l.ctx, in.GroupId)
	if err != nil {
		// 如果查询出错，返回带有错误信息的DB错误，并附上请求的GroupId
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err: %v req: %v", err, in.GroupId)
	}

	// 创建一个空的切片，用于存储群组成员信息的副本
	var respList []*social.GroupMembers
	// 使用copier库将查询结果复制到响应切片中
	err = copier.Copy(&respList, &groupMembers)
	if err != nil {
		// 如果复制过程中出错，返回错误
		return nil, err
	}

	// 返回包含群组成员列表的响应，以及nil错误（表示成功）
	return &social.GroupUsersResp{
		List: respList,
	}, nil
}
