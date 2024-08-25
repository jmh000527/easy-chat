package svc

import (
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/internal/config"
	"easy-chat/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	mqclient.MsgChatTransferClient
	mqclient.MsgReadTransferClient
}

// NewServiceContext 创建并返回一个新的 ServiceContext 实例。
//
// 参数:
//   - c: 配置结构体，包含服务所需的各种配置信息。
//
// 返回值:
//   - *ServiceContext: 一个指向新创建的 ServiceContext 的指针。
//
// 说明:
//
//	该函数通过使用传入的配置结构体 `c`，初始化并返回一个新的 ServiceContext 实例。它会根据配置创建消息传输客户端和聊天日志模型实例。
//	- `MsgChatTransferClient`: 初始化消息聊天传输客户端，用于处理聊天消息的传输。
//	- `MsgReadTransferClient`: 初始化消息已读传输客户端，用于处理消息已读状态的传输。
//	- `ChatLogModel`: 初始化聊天日志模型，用于与 MongoDB 交互，存储和检索聊天日志。
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
