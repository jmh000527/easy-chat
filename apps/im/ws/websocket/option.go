package websocket

type ServerOptions func(opt *websocketOption)

type websocketOption struct {
	Authentication
	patten string
}

func newWebsocketServerOption(opts ...ServerOptions) websocketOption {
	o := websocketOption{
		Authentication: new(webSocketAuthentication),
		patten:         "/ws",
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
