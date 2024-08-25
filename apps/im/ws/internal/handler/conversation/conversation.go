package conversation

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/wuid"
	"github.com/mitchellh/mapstructure"
	"time"
)

// Chat 处理 WebSocket 消息，进行聊天消息的转发。
//
// 该函数返回一个 websocket.HandlerFunc 处理函数，用于接收并处理聊天消息。
// 它将 WebSocket 消息解码为 ws.Chat 结构体，若消息未指定会话ID，则根据聊天类型生成会话ID。
// 处理完成后，将聊天消息推送到消息聊天传输客户端进行处理。
// 如果解码或消息处理失败，将通过 WebSocket 向客户端发送错误信息。
//
// 参数:
//   - svc: 包含服务上下文的 *svc.ServiceContext，用于访问消息聊天传输客户端。
//
// 返回:
//   - websocket.HandlerFunc: 处理 WebSocket 消息的处理函数。
func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Chat
		// 解码 WebSocket 消息数据为 ws.Chat 结构体
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			// 如果解码失败，发送错误信息到客户端
			err := srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("error message send error: %v", err)
			}
			return
		}

		// 如果消息未指定会话ID，根据聊天类型生成会话ID
		if data.ConversationId == "" {
			switch data.ChatType {
			// 单聊
			case constants.SingleChatType:
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			// 群聊
			case constants.GroupChatType:
				data.ConversationId = data.RecvId
			}
		}

		// 将聊天消息推送 kafka 消息队列进行处理
		err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixMilli(),
			MType:          data.Msg.MType,
			Content:        data.Msg.Content,
			MsgId:          msg.Id,
		})
		if err != nil {
			// 如果消息推送失败，发送错误信息到客户端
			err := srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("error message send error: %v", err)
			}
			return
		}
	}
}

// MarkRead 处理 WebSocket 消息，标记消息为已读。
//
// 该函数返回一个 websocket.HandlerFunc 处理函数，用于接收并处理标记消息为已读的请求。
// 它将 WebSocket 消息解码为 ws.MarkRead 结构体，并将其传递给消息读取传输客户端进行处理。
// 如果解码或消息处理失败，将通过 WebSocket 向客户端发送错误信息。
//
// 参数:
//   - svc: 包含服务上下文的 *svc.ServiceContext，用于访问消息读取传输客户端。
//
// 返回:
//   - websocket.HandlerFunc: 处理 WebSocket 消息的处理函数。
func MarkRead(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 已读未读处理
		var data ws.MarkRead
		// 解码 WebSocket 消息数据为 ws.MarkRead 结构体
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			// 如果解码失败，发送错误信息到客户端
			err := srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("error message send error: %v", err)
			}
			return
		}

		// 将标记已读的请求发送到消息读取传输客户端
		err := svc.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ChatType:       data.ChatType,
			ConversationId: data.ConversationId,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			MsgIds:         data.MsgIds,
		})
		if err != nil {
			// 如果消息处理失败，发送错误信息到客户端
			err := srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("error message send error: %v", err)
			}
			return
		}
	}
}
