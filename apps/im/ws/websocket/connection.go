package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// Conn 表示 WebSocket 连接。
type Conn struct {
	idleMu            sync.Mutex      // 用于保护空闲时间的互斥锁。
	Uid               string          // 连接的唯一标识符。
	wsConn            *websocket.Conn // 底层的 WebSocket 连接。
	s                 *Server         // 与该连接关联的 WebSocket 服务器。
	idle              time.Time       // 上次活动时间。
	maxConnectionIdle time.Duration   // 最大连接空闲时间。
	done              chan struct{}   // 关闭信号通道。
}

// Close 关闭 WebSocket 连接。
func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.wsConn.Close()
}

// ReadMessage 从 WebSocket 连接中读取消息。
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.wsConn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 重置空闲时间
	c.idle = time.Time{}
	return
}

// WriteMessage 向 WebSocket 连接中写入消息。
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	err := c.wsConn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}

// keepalive 保持连接的活动状态。
func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // 连接非空闲状态。
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// 连接已空闲超过 keepalive.MaxConnectionIdle 指定的时间。
				// 优雅地关闭连接。
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}

// NewConn 创建一个新的 WebSocket 连接。
func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("Error upgrading connection: %s", err)
		return nil
	}

	conn := &Conn{
		wsConn:            c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		done:              make(chan struct{}),
	}

	//go conn.keepalive()

	return conn
}
