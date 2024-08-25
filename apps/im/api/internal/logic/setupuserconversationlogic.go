package logic

import (
	"context"
	"easy-chat/apps/im/rpc/im"

	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SetUpUserConversation 设置用户之间的会话信息。
//
// 该方法用于初始化用户之间的会话信息，
// 包括发送者和接收者的用户 ID 以及聊天类型。
//
// 参数:
//   - req: 请求对象，包含需要设置的会话信息。
//
// 返回值:
//   - *types.SetUpUserConversationResp: 响应对象，表示设置操作的结果。
//   - error: 如果在设置过程中发生错误，返回具体的错误信息；成功时返回 nil。
func (l *SetUpUserConversationLogic) SetUpUserConversation(req *types.SetUpUserConversationReq) (resp *types.SetUpUserConversationResp, err error) {
	// 调用服务上下文中的方法进行会话设置
	_, err = l.svcCtx.SetUpUserConversation(l.ctx, &im.SetUpUserConversationReq{
		SendId:   req.SendId,
		RecvId:   req.RecvId,
		ChatType: req.ChatType,
	})

	// 返回设置结果
	return
}
