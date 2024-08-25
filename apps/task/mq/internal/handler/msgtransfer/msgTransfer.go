package msgtransfer

import (
	"context"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/logx"
)

type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

// Transfer 处理消息的转发，根据聊天类型决定消息的处理方式。
//
// 该方法根据消息的聊天类型调用相应的处理函数：
// - 对于单聊消息，调用 single 方法。
// - 对于群聊消息，调用 group 方法。
//
// 参数:
//   - ctx: 上下文对象，用于传递请求范围的数据。
//   - data: 包含要推送的数据的 Push 结构体。
//
// 返回值:
//   - error: 如果处理过程中出现错误，返回相应的错误；否则返回 nil。
func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.SingleChatType:
		err = m.single(ctx, data)
	case constants.GroupChatType:
		err = m.group(ctx, data)
	}
	return err
}

// single 处理单聊消息的转发。
//
// 该方法通过 WebSocket 客户端将单聊消息推送给指定的用户。
//
// 参数:
//   - ctx: 上下文对象，用于传递请求范围的数据。
//   - data: 包含要推送的数据的 Push 结构体。
//
// 返回值:
//   - error: 如果推送过程中出现错误，返回相应的错误；否则返回 nil。
func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	// 推送消息
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SystemRootUid,
		Data:      data,
	})
}

// group 处理群聊消息的转发。
//
// 该方法首先查询群成员，然后将群聊消息推送给所有群成员，
// 除了消息的发送者外。
//
// 参数:
//   - ctx: 上下文对象，用于传递请求范围的数据。
//   - data: 包含要推送的数据的 Push 结构体。
//
// 返回值:
//   - error: 如果查询群成员或推送消息过程中出现错误，返回相应的错误；否则返回 nil.
func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 查询群用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}

	// 获取待发送的群用户ID
	data.RecvIds = make([]string, 0, len(users.List))
	for _, user := range users.List {
		// 不包含发送者自己
		if user.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, user.UserId)
	}

	// 向用户发送消息
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SystemRootUid,
		Data:      data,
	})
}
