package websocket

import (
	"easy-chat/apps/im/ws/websocket/auth"
	"time"
)

// ServerOptions 定义 WebSocket 服务器的选项配置函数。
type ServerOptions func(opt *websocketOption)

type websocketOption struct {
	auth.Authentication        // WebSocket 服务器的身份认证设置
	patten              string // WebSocket 路由模式

	ack          AckType       // 消息确认类型
	ackTimeout   time.Duration // 消息确认超时时间
	sendErrCount int           // 发送错误次数限制

	maxConnectionIdle time.Duration // 最大连接空闲时间

	concurrency int // 群消息并发处理量级
}

// newWebsocketServerOption 创建一个新的 websocketOption 实例。
//
// 该函数初始化 WebSocket 服务器的选项，并应用提供的 ServerOptions 函数。
//
// 参数:
//   - opts: 可选的 ServerOptions 函数，用于配置 WebSocket 服务器选项。
//
// 返回:
//   - websocketOption: 配置好的 WebSocket 服务器选项。
func newWebsocketServerOption(opts ...ServerOptions) websocketOption {
	o := websocketOption{
		Authentication:    new(auth.WebSocketAuth),
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		sendErrCount:      defaultSendErrCount,
		patten:            "/ws",
		concurrency:       defaultConcurrency,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithWebsocketAuthentication 配置 WebSocket 服务器的身份认证。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的身份认证。
//
// 参数:
//   - auth: 用于 WebSocket 服务器的身份认证实例。
//
// 返回:
//   - ServerOptions: 配置身份认证的函数。
func WithWebsocketAuthentication(auth auth.Authentication) ServerOptions {
	return func(opt *websocketOption) {
		opt.Authentication = auth
	}
}

// WithWebsocketPatten 配置 WebSocket 的路由模式。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的路由模式。
//
// 参数:
//   - patten: WebSocket 的路由模式路径。
//
// 返回:
//   - ServerOptions: 配置路由模式的函数。
func WithWebsocketPatten(patten string) ServerOptions {
	return func(opt *websocketOption) {
		opt.patten = patten
	}
}

// WithServerAck 配置消息确认类型。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的消息确认类型。
//
// 参数:
//   - ack: 消息确认类型。
//
// 返回:
//   - ServerOptions: 配置消息确认类型的函数。
func WithServerAck(ack AckType) ServerOptions {
	return func(opt *websocketOption) {
		opt.ack = ack
	}
}

// WithServerSendErrCount 配置发送错误次数限制。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的发送错误次数限制。
//
// 参数:
//   - sendErrCount: 发送错误次数限制。
//
// 返回:
//   - ServerOptions: 配置发送错误次数限制的函数。
func WithServerSendErrCount(sendErrCount int) ServerOptions {
	return func(opt *websocketOption) {
		opt.sendErrCount = sendErrCount
	}
}

// WithWebsocketMaxConnectionIdle 配置最大连接空闲时间。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的最大连接空闲时间。
// 如果当前设置的 maxConnectionIdle 大于0，则更新为提供的时间。
//
// 参数:
//   - duration: 最大连接空闲时间。
//
// 返回:
//   - ServerOptions: 配置最大连接空闲时间的函数。
func WithWebsocketMaxConnectionIdle(duration time.Duration) ServerOptions {
	return func(opt *websocketOption) {
		if opt.maxConnectionIdle > 0 {
			opt.maxConnectionIdle = duration
		}
	}
}
