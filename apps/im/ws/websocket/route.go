package websocket

type HandlerFunc func(srv *Server, conn *Conn, msg *Message)

type Route struct {
	Method  string
	Handler HandlerFunc
}
