package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// Conn 表示 WebSocket 连接。
type Conn struct {
	idleMu sync.Mutex
	Uid    string
	wsConn *websocket.Conn
	s      *Server

	idle              time.Time
	maxConnectionIdle time.Duration

	messageMu      sync.Mutex
	readMessage    []*Message          // 读消息队列
	readMessageSeq map[string]*Message // 读消息队列序列化

	message chan *Message
	done    chan struct{}
}

func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()
	// 读队列中，判断之前是否在队列中存过消息
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 该消息已经有Ack的确认过程
		if len(c.readMessage) == 0 {
			// 队列中没有该消息
			return
		}
		// 要求 msg.AckSeq > m.AckSeq
		if msg.AckSeq <= m.AckSeq {
			// 没有进行ack的确认, 重复
			return
		}
		// 更新最新的消息
		c.readMessageSeq[msg.Id] = msg
		return
	}
	// 还没有进行Ack确认，避免客户端重复发送多余的Ack消息
	if msg.FrameType == FrameAck {
		return
	}
	// 记录消息
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg
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
	var responseHeader http.Header
	if protocol := r.Header.Get("Sec-WebSocket-Protocol"); protocol != "" {
		responseHeader = http.Header{
			"Sec-WebSocket-Protocol": []string{protocol},
		}
	}

	c, err := s.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		s.Errorf("Error upgrading connection: %s", err)
		return nil
	}
	conn := &Conn{
		wsConn:            c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		message:           make(chan *Message, 1),
		done:              make(chan struct{}),
	}
	go conn.keepalive()
	return conn
}
