package group

import (
	"context"
	"easy-chat/apps/im/rpc/imclient"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/ctxdata"
	"errors"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateGroup 创建一个新的群组并建立会话
//
// 功能描述:
//   - 从上下文中获取当前用户ID（群组创建者）
//   - 调用服务层接口创建群组，并获取群组ID
//   - 如果群组创建成功，则建立群组会话
//
// 参数:
//   - req: `*types.GroupCreateReq` 类型，包含群组创建所需的信息（群组名称、图标等）
//
// 返回值:
//   - `*types.GroupCreateResp`: 响应对象，当前未使用
//   - `error`: 如果在创建群组或建立会话过程中发生错误，则返回相应的错误信息
func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	// 从上下文中获取群组创建者ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用服务层接口创建群组
	res, err := l.svcCtx.Social.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name, // 群组名称
		Icon:       req.Icon, // 群组图标
		CreatorUid: uid,      // 群组创建者ID
	})
	if err != nil {
		// 如果创建群组失败，返回错误信息
		return nil, err
	}

	// 检查群组ID是否为空，确保群组创建成功
	if res.Id == "" {
		return nil, errors.New("failed to create group: empty group ID")
	}

	// 创建群组会话
	_, err = l.svcCtx.Im.CreateGroupConversation(l.ctx, &imclient.CreateGroupConversationReq{
		GroupId:  res.Id, // 新创建的群组ID
		CreateId: uid,    // 创建会话的用户ID（群组创建者）
	})
	if err != nil {
		// 如果建立群组会话失败，返回错误信息
		return nil, err
	}

	// 返回成功的响应，当前未使用
	return &types.GroupCreateResp{}, nil
}
