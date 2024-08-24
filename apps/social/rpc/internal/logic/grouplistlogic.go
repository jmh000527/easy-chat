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

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupList 获取用户所在的群组列表
//
// 功能描述:
//   - 根据用户ID查询用户所在的群组成员列表，然后获取这些群组的详细信息，并返回群组列表。
//
// 参数:
//   - in: `social.GroupListReq` 类型，包含用户ID，用于查询用户所在的群组。
//
// 返回值:
//   - `*social.GroupListResp`: 包含用户所在的群组列表的响应对象。
//   - `error`: 如果在获取群组列表过程中发生错误，则返回相应的错误信息。
func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	// 根据用户ID查询用户所在的群组成员列表
	userGroups, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		// 如果查询群组成员列表失败，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err: %v req: %v", err, in.UserId)
	}

	// 如果用户没有加入任何群组，则返回一个空的群组列表响应
	if len(userGroups) == 0 {
		return &social.GroupListResp{}, nil
	}

	// 提取用户所在的所有群组的ID
	ids := make([]string, 0, len(userGroups))
	for _, v := range userGroups {
		ids = append(ids, v.GroupId)
	}

	// 根据群组ID列表查询群组信息
	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, ids)
	if err != nil {
		// 如果查询群组信息失败，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group err: %v req: %v", err, ids)
	}

	// 将查询到的群组信息复制到响应列表中
	var respList []*social.Groups
	err = copier.Copy(&respList, &groups)
	if err != nil {
		return nil, err
	}

	// 返回群组列表响应，包括群组列表和nil错误信息
	return &social.GroupListResp{
		List: respList,
	}, nil
}
