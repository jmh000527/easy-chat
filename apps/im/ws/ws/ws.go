package ws

import "easy-chat/pkg/constants"

// Msg 表示一个基础消息的结构体。
//
// 该结构体包含消息的唯一标识符、已读记录、消息类型和消息内容。
type Msg struct {
	MsgId           string                 `mapstructure:"msgId"`       // 消息的唯一标识符
	ReadRecords     map[string]string      `mapstructure:"readRecords"` // 消息的已读记录，键为用户ID，值为已读时间戳
	constants.MType `mapstructure:"mType"` // 消息的类型，定义在 constants 中
	Content         string                 `mapstructure:"content"` // 消息的实际内容
}

// Chat 表示一个聊天消息的结构体。
//
// 该结构体继承了 Msg 结构体，包含了会话ID、聊天类型、发送者和接收者ID、发送时间等信息。
type Chat struct {
	ConversationId     string                    `mapstructure:"conversationId"` // 聊天会话的唯一标识符
	constants.ChatType `mapstructure:"chatType"` // 聊天的类型，定义在 constants 中
	SendId             string                    `mapstructure:"sendId"`   // 发送者的唯一标识符
	RecvId             string                    `mapstructure:"recvId"`   // 接收者的唯一标识符
	SendTime           int64                     `mapstructure:"sendTime"` // 消息发送的时间戳
	Msg                `mapstructure:"msg"`      // 嵌入的消息结构体，包含消息的详细信息
}

// Push 表示一个推送消息的结构体。
//
// 该结构体包含了推送消息所需的信息，包括会话ID、发送者和接收者ID列表、发送时间、消息内容等。
type Push struct {
	ConversationId     string                    `mapstructure:"conversationId"` // 推送消息所属的会话ID
	constants.ChatType `mapstructure:"chatType"` // 推送消息的聊天类型
	SendId             string                    `mapstructure:"sendId"`   // 推送消息的发送者ID
	RecvId             string                    `mapstructure:"recvId"`   // 单一接收者的ID
	RecvIds            []string                  `mapstructure:"recvIds"`  // 多个接收者的ID列表
	SendTime           int64                     `mapstructure:"sendTime"` // 推送消息发送的时间戳

	MsgId       string                `mapstructure:"msgId"`       // 消息的唯一标识符
	ReadRecords map[string]string     `mapstructure:"readRecords"` // 消息的已读记录，键为用户ID，值为已读时间戳
	ContentType constants.ContentType `mapstructure:"contentType"` // 消息内容的类型，定义在 constants 中

	constants.MType `mapstructure:"mType"` // 消息的类型，定义在 constants 中
	Content         string                 `mapstructure:"content"` // 推送消息的实际内容
}

// MarkRead 表示一个标记消息已读的结构体。
//
// 该结构体用于处理标记消息已读的操作，包括会话ID、接收者ID和已读的消息ID列表。
type MarkRead struct {
	constants.ChatType `mapstructure:"chatType"` // 聊天的类型，定义在 constants 中
	RecvId             string                    `mapstructure:"recvId"`         // 已读结果的接收者ID
	ConversationId     string                    `mapstructure:"conversationId"` // 会话的唯一标识符
	MsgIds             []string                  `mapstructure:"msgIds"`         // 已读消息的ID列表
}
