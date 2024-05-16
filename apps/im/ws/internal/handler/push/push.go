package push

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 解析消息
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			err := srv.Send(websocket.NewErrMessage(err))
			if err != nil {
				srv.Errorf("push err: %v", err)
			}
			return
		}
		// 发送的目标
		switch data.ChatType {
		case constants.SingleChatType:
			err := single(srv, &data, data.RecvId)
			if err != nil {
				srv.Errorf("push err: %v", err)
				return
			}
		case constants.GroupChatType:
			group(srv, &data)
		}
	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 获取发送的目标
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// todo: 目标离线
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
