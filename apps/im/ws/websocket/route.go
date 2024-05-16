package websocket

import "github.com/gorilla/websocket"

type HandlerFunc func(srv *Server, conn *websocket.Conn, msg *Message)

type Route struct {
	Method  string
	Handler HandlerFunc
}
