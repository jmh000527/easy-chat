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

// GroupList 是GroupListLogic结构体的方法，用于获取用户所在的群组列表
func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	// 根据用户ID查询用户所在的群组成员列表
	userGroup, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		// 如果查询群组成员列表失败，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err: %v req: %v", err, in.UserId)
	}
	// 如果用户没有加入任何群组，则返回一个空的群组列表响应
	if len(userGroup) == 0 {
		return &social.GroupListResp{}, nil
	}

	// 提取用户所在的所有群组的ID
	ids := make([]string, 0, len(userGroup))
	for _, v := range userGroup {
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
