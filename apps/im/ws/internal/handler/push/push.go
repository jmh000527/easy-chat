package push

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
	"github.com/mitchellh/mapstructure"
)

// Push 处理 WebSocket 消息，转发推送消息，由 kafka 消息队列远程调用。
//
// 该函数返回一个 websocket.HandlerFunc 处理函数，用于接收并处理推送消息。
// 它将 WebSocket 消息解码为 ws.Push 结构体，并根据聊天类型将消息推送到目标用户。
// 如果消息解码失败，或推送过程中出现错误，将通过 WebSocket 向客户端发送错误信息。
//
// 参数:
//   - svc: 包含服务上下文的 *svc.ServiceContext，用于访问服务相关功能。
//
// 返回:
//   - websocket.HandlerFunc: 处理 WebSocket 消息的处理函数。
func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 解析 WebSocket 消息数据为 ws.Push 结构体
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			// 如果解码失败，发送错误信息到客户端
			err := srv.Send(websocket.NewErrMessage(err))
			if err != nil {
				srv.Errorf("push err: %v", err)
			}
			return
		}

		// 根据聊天类型进行不同的推送处理
		switch data.ChatType {
		case constants.SingleChatType:
			// 处理单聊消息推送
			err := single(srv, &data, data.RecvId)
			if err != nil {
				srv.Errorf("push err: %v", err)
				return
			}
		case constants.GroupChatType:
			// 处理群聊消息推送
			group(srv, &data)
		}
	}
}

// single 处理单聊消息的推送。
//
// 该函数根据接收者ID从服务器获取连接，并将消息推送给接收者。
// 如果目标用户离线，当前实现没有处理离线用户的逻辑。
// 如果推送过程中出现错误，记录错误日志。
//
// 参数:
//   - srv: WebSocket 服务器实例。
//   - data: 包含推送消息的数据结构体。
//   - recvId: 接收者用户ID。
//
// 返回:
//   - error: 发生的错误（如果有的话），返回nil表示推送成功。
func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 获取发送的目标用户连接
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// 目标用户离线，当前实现未处理离线用户的逻辑
		return nil
	}
	// 发送消息
	srv.Infof("push msg: %v", data)
	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MsgId:       data.MsgId,
			ReadRecords: data.ReadRecords,
			MType:       data.MType,
			Content:     data.Content,
		},
	}), rconn)
}

// group 处理群聊消息的推送。
//
// 该函数将消息推送到所有指定的群聊成员。
// 对于每个成员，使用单聊消息推送的方式发送消息。
// 错误日志会在单聊消息推送过程中记录。
//
// 参数:
//   - srv: WebSocket 服务器实例。
//   - data: 包含推送消息的数据结构体。
func group(srv *websocket.Server, data *ws.Push) {
	for _, id := range data.RecvIds {
		func(recvId string) {
			srv.Schedule(func() {
				err := single(srv, data, recvId)
				if err != nil {
					srv.Errorf("push err: %v", err)
					return
				}
			})
		}(id)
	}
	return
}
