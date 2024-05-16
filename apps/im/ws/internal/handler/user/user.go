package user

import (
	"easy-chat/apps/im/ws/internal/svc"
	websocketx "easy-chat/apps/im/ws/websocket"
	"github.com/gorilla/websocket"
)

func OnLine(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocket.Conn, msg *websocketx.Message) {
		uids := srv.GetUsers()

		userId := srv.GetUsers(conn)
		err := srv.Send(websocketx.NewMessage(userId[0], uids), conn)
		srv.Info("err: ", err)
	}
}
