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

func (l *GroupUsersLogic) GroupUsers(in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	// 通过GroupId查询群组成员，可能会返回错误
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByGroupId(l.ctx, in.GroupId)
	if err != nil {
		// 如果查询出错，则返回带有错误信息的DB错误，同时附上请求的GroupId
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err: %v req: %v", err, in.GroupId)
	}

	// 创建一个空的social.GroupMembers类型切片，用于存储复制的群组成员信息
	var respList []*social.GroupMembers
	// 使用copier库将groupMembers复制到respList中
	err = copier.Copy(&respList, &groupMembers)
	if err != nil {
		return nil, err
	}

	// 返回一个social.GroupUsersResp类型的指针，其中包含复制的群组成员列表，以及nil错误（表示无错误）
	return &social.GroupUsersResp{
		List: respList,
	}, nil
}
