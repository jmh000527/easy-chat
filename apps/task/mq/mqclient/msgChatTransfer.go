package mqclient

import (
	"easy-chat/apps/task/mq/mq"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
)

// MsgChatTransferClient 提供发送聊天消息的方法。
//
// 该接口定义了发送聊天消息的操作。
type MsgChatTransferClient interface {
	// Push 发送聊天消息。
	//
	// 参数:
	//   - msg: 包含聊天消息的结构体，该消息将被发送到消息队列中。
	//
	// 返回值:
	//   - error: 如果发送过程中出现错误，则返回相应的错误信息；否则返回 nil。
	Push(msg *mq.MsgChatTransfer) error
}

// msgChatTransferClient 实现了 MsgChatTransferClient 接口，用于将聊天消息推送到消息队列中。
type msgChatTransferClient struct {
	pusher *kq.Pusher
}

// NewMsgChatTransferClient 创建一个新的 MsgChatTransferClient 实例。
//
// 该函数用于初始化并返回一个新的消息推送客户端。
//
// 参数:
//   - addr: 消息队列的地址列表。
//   - topic: 消息主题。
//   - opts: 其他可选的推送选项。
//
// 返回值:
//   - MsgChatTransferClient: 初始化好的消息推送客户端实例。
func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

// Push 将聊天消息推送到消息队列中。
//
// 该方法将聊天消息序列化为 JSON 格式，并通过 pusher 推送到消息队列中。
//
// 参数:
//   - msg: 包含聊天消息的结构体，该消息将被发送到消息队列中。
//
// 返回值:
//   - error: 如果发送过程中出现错误，则返回相应的错误信息；否则返回 nil。
func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.pusher.Push(string(body))
}
