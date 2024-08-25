package logic

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/wuid"
	"time"
)

type ConversationLogic struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *ConversationLogic {
	return &ConversationLogic{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

// SingleChat 处理单聊消息
//
// 该方法负责处理单聊消息，首先检查是否已存在会话ID，如果不存在则生成一个新的会话ID。
// 然后记录聊天日志到数据库中。
//
// 参数:
// - data: 包含聊天消息数据的结构体，包括会话ID、接收者ID、聊天类型、消息类型和内容。
// - userId: 发送者的用户ID。
//
// 返回:
// - error: 发生的错误（如果有的话），返回nil表示操作成功。
func (l *ConversationLogic) SingleChat(data *ws.Chat, userId string) error {
	// 查看是否存在会话ID，否则新建一个会话ID
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	// 记录消息
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}

	// 将聊天日志插入数据库
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)

	// 返回操作结果
	return err
}
