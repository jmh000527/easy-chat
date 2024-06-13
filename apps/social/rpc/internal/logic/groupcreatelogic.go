package logic

import (
	"context"
	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/wuid"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// 创建一个Groups结构体实例，并设置相关属性
	groups := &socialmodels.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource), // 生成唯一ID作为群组ID
		Name:       in.Name,                                       // 设置群组名称
		Icon:       in.Icon,                                       // 设置群组图标
		CreatorUid: in.CreatorUid,                                 // 设置群组创建者ID
		// IsVerify设置为false，表示群组不需要验证（可选，默认值为false）
		//IsVerify:   true,
		IsVerify: false,
	}

	// 调用GroupsModel的Trans方法，开启一个事务，并执行相关数据库操作
	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务中插入群组信息到数据库
		_, err := l.svcCtx.GroupsModel.Insert(l.ctx, session, groups)
		if err != nil {
			// 插入群组信息失败，返回错误信息
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v req %v", err, in)
		}

		// 插入群组成员信息到数据库，成员包括创建者，并设置其角色等级为CreatorGroupRoleLevel
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,                            // 设置群组成员所属群组ID
			UserId:    in.CreatorUid,                        // 设置群组成员用户ID
			RoleLevel: int(constants.CreatorGroupRoleLevel), // 设置群组成员角色等级
		})
		if err != nil {
			// 插入群组成员信息失败，返回错误信息
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
		}
		return nil // 事务执行成功，返回nil
	})

	// 返回群组创建响应，包括群组ID和错误信息（如果有）
	return &social.GroupCreateResp{
		Id: groups.Id, // 设置响应中的群组ID
	}, err
}
