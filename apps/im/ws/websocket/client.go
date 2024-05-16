package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/url"
)

// Client 表示 WebSocket 客户端。
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
func NewClient(host string, opts ...DialOptions) *client {
	opt := newDialOptions(opts...)
	c := &client{
		Conn: nil,
		host: host,
		opt:  opt,
	}
	conn, err := c.dial()
	if err != nil {
		panic(err)
	}
	c.Conn = conn
	return c
}

// dial 与 WebSocket 服务器建立连接。
func (c *client) dial() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

// Close 关闭 WebSocket 连接。
func (c *client) Close() error {
	if c.Conn == nil {
		return errors.New("connection is nil")
	}
	return c.Conn.Close()
}

// Send 序列化并发送消息到 WebSocket。
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
