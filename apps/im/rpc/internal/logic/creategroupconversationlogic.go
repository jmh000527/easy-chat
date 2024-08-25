package logic

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateGroupConversation 创建群聊会话
//
// 该方法用于创建群聊会话。如果指定的群ID已存在，则直接返回；
// 如果群ID不存在，则新建一个群聊会话，并将创建者的用户会话列表进行更新。
//
// 参数:
//   - in: 包含群聊会话创建请求的结构体，包括群ID和创建者用户ID。
//
// 返回:
//   - *im.CreateGroupConversationResp: 创建群聊会话的响应结构体，包含相关的响应数据。
//   - error: 发生的错误（如果有的话），返回nil表示操作成功。
func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	resp := &im.CreateGroupConversationResp{}

	// 查询是否已存在指定的群聊会话
	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	// 如果会话存在，则直接返回成功响应
	if err == nil {
		return resp, nil
	}
	// 如果发生其他错误（非未找到错误），返回数据库错误
	if !errors.Is(err, immodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "cannot find group conversation with group id %s, err: %v", in.GroupId, err)
	}

	// 如果会话不存在，则创建新的群聊会话
	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "cannot create group conversation with group id %s, err: %v", in.GroupId, err)
	}

	// 更新创建者用户的会话列表
	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	// 返回响应及可能发生的错误
	return resp, nil
}
