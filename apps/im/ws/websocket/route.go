package websocket

// HandlerFunc 定义了处理 WebSocket 消息的函数类型。
//
// 参数:
//   - srv: *Server
//     表示当前的 WebSocket 服务器实例，提供了与服务器交互的方法。
//   - conn: *Conn
//     表示当前的 WebSocket 连接，代表一个与客户端的连接实例，提供了与客户端通信的方法。
//   - msg: *Message
//     表示收到的消息，包含消息的详细信息，包括消息的类型、内容等。
type HandlerFunc func(srv *Server, conn *Conn, msg *Message)

// Route 表示一个路由条目，用于将特定的请求方法映射到对应的处理函数。
//
// 字段:
//   - Method: 字符串类型，表示请求的方法（例如 "GET"、"POST"、"PUT" 等）。
//   - Handler: 函数类型，表示与该请求方法关联的处理函数。该处理函数负责处理匹配此路由的请求。
type Route struct {
	Method  string      // 请求方法
	Handler HandlerFunc // 与请求方法关联的处理函数
}
