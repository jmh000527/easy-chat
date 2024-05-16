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

func (l *ConversationLogic) SingleChat(data *ws.Chat, userId string) error {
	// 查看是否存在会话ID，否则新建一个会话ID
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}
	//time.Sleep(time.Minute)
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
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)

	return err
}
