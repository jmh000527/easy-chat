package user

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
)

// OnLine 处理 WebSocket 消息，向客户端发送在线用户列表。
//
// 该函数返回一个 websocket.HandlerFunc 处理函数，用于接收并处理请求在线用户列表的消息。
// 它从 WebSocket 服务器中获取所有在线用户的列表，并将该列表发送到请求的客户端。
// 如果消息发送过程中出现错误，将记录错误信息。
//
// 参数:
//   - svc: 包含服务上下文的 *svc.ServiceContext，用于访问服务相关功能。
//
// 返回:
//   - websocket.HandlerFunc: 处理 WebSocket 消息的处理函数。
func OnLine(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 获取所有在线用户ID
		uids := srv.GetUsers()

		// 获取请求者的用户ID
		users := srv.GetUsers(conn)

		// 发送在线用户列表到请求的客户端
		err := srv.Send(websocket.NewMessage(users[0], uids), conn)

		// 记录错误信息（如果有的话）
		srv.Info("err: ", err)
	}
}
