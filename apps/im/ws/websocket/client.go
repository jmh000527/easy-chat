package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/url"
)

// Client 表示 WebSocket 客户端，在kafka中消费。
//
// 该接口定义了 WebSocket 客户端应实现的方法，包括关闭连接、发送消息和读取消息。
type Client interface {
	Close() error     // 关闭 WebSocket 连接。
	Send(v any) error // 发送消息到 WebSocket。
	Read(v any) error // 从 WebSocket 读取消息。
}

type client struct {
	*websocket.Conn            // WebSocket 连接。
	host            string     // WebSocket 服务器的主机地址。
	opt             dialOption // WebSocket 连接的拨号选项。
}

// NewClient 创建一个新的 WebSocket 客户端。
//
// 该函数用于创建一个新的 WebSocket 客户端实例，初始化连接到指定的 WebSocket 服务器。
// 如果连接失败，会导致程序 panic。
//
// 参数:
//   - host: WebSocket 服务器的主机地址。
//   - opts: 可选的拨号选项，用于配置 WebSocket 连接。
//
// 返回:
//   - *client: 新创建的 WebSocket 客户端实例。
func NewClient(host string, opts ...DialOptions) *client {
	// 创建新的拨号选项。
	opt := newDialOptions(opts...)
	// 创建新的 WebSocket 客户端实例。
	c := &client{
		Conn: nil,
		host: host,
		opt:  opt,
	}
	// 拨号连接到 WebSocket 服务器。
	conn, err := c.dial()
	if err != nil {
		panic(err)
	}
	// 将连接赋值给客户端。
	c.Conn = conn
	return c
}

// dial 与 WebSocket 服务器建立连接。
//
// 该方法用于与 WebSocket 服务器建立连接，并返回一个 WebSocket 连接实例。
// 如果连接失败，则返回错误。
//
// 返回:
//   - *websocket.Conn: 成功建立的 WebSocket 连接。
//   - error: 连接过程中发生的错误（如果有的话）。
func (c *client) dial() (*websocket.Conn, error) {
	// 构造WebSocket连接的URL。
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}

	// 使用DefaultDialer进行WebSocket连接。
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	if err != nil {
		return nil, err // 如果连接失败，返回错误。
	}

	return conn, nil // 如果连接成功，返回连接对象。
}

// Close 关闭 WebSocket 连接。
//
// 该方法用于关闭 WebSocket 连接。如果连接已经是 nil，则返回一个错误。
//
// 返回:
//   - error: 关闭连接过程中发生的错误（如果有的话）。
func (c *client) Close() error {
	if c.Conn == nil {
		return errors.New("connection is nil")
	}
	return c.Conn.Close()
}

// Send 序列化并发送消息到 WebSocket。
//
// 该方法将消息对象序列化为 JSON 格式，并通过 WebSocket 连接发送。
// 如果发送失败，会尝试重新连接并重新发送消息。
//
// 参数:
//   - v: 要发送的消息对象，可以是任意类型。
//
// 返回:
//   - error: 发送消息过程中发生的错误（如果有的话）。
func (c *client) Send(v any) error {
	if c.Conn == nil {
		return errors.New("connection is nil")
	}
	// 序列化消息。
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	// 发送消息。
	err = c.Conn.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return nil
	}
	// 重新连接并重新发送消息。
	conn, err := c.dial()
	if err != nil {
		return err
	}
	c.Conn = conn
	err = c.Conn.WriteMessage(websocket.TextMessage, data)
	return err
}

// Read 从 WebSocket 读取消息并反序列化。
//
// 该方法从 WebSocket 连接中读取消息，并将其反序列化为指定的对象类型。
// 如果读取或反序列化过程中发生错误，则返回错误。
//
// 参数:
//   - v: 用于接收反序列化后的消息对象，可以是任意类型。
//
// 返回:
//   - error: 读取消息过程中发生的错误（如果有的话）。
func (c *client) Read(v any) error {
	if c.Conn == nil {
		return errors.New("connection is nil")
	}
	// 读取消息。
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}
	// 反序列化消息。
	return json.Unmarshal(msg, v)
}
