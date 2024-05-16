package websocket

import "time"

type ServerOptions func(opt *websocketOption)

type websocketOption struct {
	Authentication
	patten string

	ack          AckType
	ackTimeout   time.Duration
	sendErrCount int

	maxConnectionIdle time.Duration
}

func newWebsocketServerOption(opts ...ServerOptions) websocketOption {
	o := websocketOption{
		Authentication:    new(webSocketAuthentication),
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		sendErrCount:      defaultSendErrCount,
		patten:            "/ws",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithWebsocketAuthentication(auth Authentication) ServerOptions {
	return func(opt *websocketOption) {
		opt.Authentication = auth
	}
}

func WithWebsocketPatten(patten string) ServerOptions {
	return func(opt *websocketOption) {
		opt.patten = patten
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *websocketOption) {
		opt.ack = ack
	}
}

func WithServerSendErrCount(sendErrCount int) ServerOptions {
	return func(opt *websocketOption) {
		opt.sendErrCount = sendErrCount
	}
}

func WithWebsocketMaxConnectionIdle(duration time.Duration) ServerOptions {
	return func(opt *websocketOption) {
		if opt.maxConnectionIdle > 0 {
			opt.maxConnectionIdle = duration
		}
	}
}
