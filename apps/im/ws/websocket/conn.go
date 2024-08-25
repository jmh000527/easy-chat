package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// Conn 表示 WebSocket 连接。
//
// 该结构体定义了一个WebSocket连接的主要属性和状态，包括用户ID、WebSocket连接实例、
// 连接的管理服务器以及用于消息处理和连接状态维护的各种信息。此结构体用于管理和维护
// WebSocket连接的生命周期和消息传递。
//
// 字段:
//   - idleMu: 连接空闲状态的互斥锁，用于保护空闲时间的读写操作。
//   - Uid: 用户标识符，用于标识与该连接关联的用户。
//   - wsConn: WebSocket连接实例，表示与客户端的实际WebSocket连接。
//   - s: 连接所属的WebSocket服务器，用于访问服务器相关的功能和状态。
//   - idle: 连接的空闲时间，用于检测连接的活动状态。
//   - maxConnectionIdle: 允许的最大空闲时间，超过该时间连接将被认为是超时。
//   - messageMu: 消息队列的互斥锁，用于保护消息队列的读写操作。
//   - readMessage: 读消息队列，存储尚未处理的消息。
//   - readMessageSeq: 读消息队列的序列化映射，用于按序号存储消息。
//   - message: 消息通道，用于接收和发送消息。
//   - done: 关闭连接时的信号通道，用于通知连接的结束。
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

// appendMsgMq 将消息添加到消息队列中。
//
// 该方法用于将传入的消息添加到读消息队列中，并维护消息序列化映射。
// 如果消息已经存在于队列中且其确认序号小于等于之前存储的消息，则忽略该消息。
// 如果消息类型是确认消息（FrameAck），则不处理。
// 否则，将消息添加到队列并更新消息序列化映射。
//
// 参数:
//   - msg: 要添加到消息队列中的消息结构体。
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
//
// 该方法关闭WebSocket连接并通知所有相关操作停止。
// 如果连接已经关闭，则不会重复关闭。关闭操作通过信号通道通知。
// 关闭操作返回WebSocket连接关闭的错误（如果有的话）。
//
// 返回:
//   - error: 关闭WebSocket连接时发生的错误，如果成功关闭则返回nil。
func (c *Conn) Close() error {
	select {
	case <-c.done:
		// 如果 c.done 通道已经关闭，执行此分支，什么都不做。
	default:
		// 如果 c.done 通道未关闭，执行此分支，关闭该通道。
		close(c.done)
	}

	// 关闭 websocket 连接
	return c.wsConn.Close()
}

// ReadMessage 从 WebSocket 连接中读取消息。
//
// 该方法从WebSocket连接中读取消息，并重置连接的空闲时间。
// 如果读取消息时发生错误，则返回该错误。
// 空闲时间用于管理连接的活跃状态。
//
// 返回:
//   - messageType: 消息类型，指示消息的格式（文本或二进制）。
//   - p: 读取到的消息内容。
//   - err: 读取消息时发生的错误，如果没有错误则返回nil。
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.wsConn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 置零，表示当前连接不再空闲
	c.idle = time.Time{}
	return
}

// WriteMessage 向 WebSocket 连接中写入消息。
//
// 该方法将消息写入WebSocket连接，并更新连接的空闲时间。
// 如果写入消息时发生错误，则返回该错误。
// 空闲时间用于管理连接的活跃状态。
//
// 参数:
//   - messageType: 消息类型，指示消息的格式（文本或二进制）。
//   - data: 要写入的消息内容。
//
// 返回:
//   - error: 写入消息时发生的错误，如果成功写入则返回nil。
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	err := c.wsConn.WriteMessage(messageType, data)
	// 更新空闲时间，表示当前连接空闲
	c.idle = time.Now()
	return err
}

// keepalive 定期检查连接的空闲状态，确保连接在超过最大空闲时间后被优雅地关闭。
// 该方法会启动一个定时器，根据连接的最大空闲时间进行检查。
// 如果连接空闲时间超过了最大空闲时间，连接将被关闭。
// 如果连接未超过空闲时间，定时器将重置以继续监控连接状态。
// 方法会监听连接的关闭事件，以便在连接关闭时终止检查。
func (c *Conn) keepalive() {
	// 创建一个新的计时器，定时检查连接的空闲状态
	idleTimer := time.NewTimer(c.maxConnectionIdle)

	// 确保在方法结束时停止计时器
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		// 当计时器超时时，执行空闲状态检查
		case <-idleTimer.C:
			c.idleMu.Lock() // 锁定连接的空闲状态
			// 获取连接的空闲时间
			idle := c.idle
			if idle.IsZero() { // 连接非空闲状态
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle) // 重置计时器
				continue
			}
			// 计算剩余的最大空闲时间
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// 连接已空闲超过最大空闲时间，优雅地关闭连接
				c.s.Close(c)
				return
			}
			// 重置计时器以继续监控
			idleTimer.Reset(val)
		// 监听连接关闭事件
		case <-c.done:
			return
		}
	}
}

// NewConn 创建一个新的 WebSocket 连接。
//
// 该函数用于创建一个新的 WebSocket 连接，并返回一个包含连接信息的 `Conn` 对象。
// 它使用服务器的升级器将 HTTP 请求升级为 WebSocket 连接，并根据请求头中的 "Sec-WebSocket-Protocol" 头设置响应头。
// 连接成功后，初始化连接对象的相关字段，并启动一个后台协程用于保持连接的活动状态。
//
// 参数:
//   - s: 服务器实例，提供用于 WebSocket 升级的上下文和配置。
//   - w: HTTP 响应写入器，用于将 WebSocket 升级响应写入客户端。
//   - r: HTTP 请求对象，包含 WebSocket 升级请求的信息。
//
// 返回:
//   - *Conn: 新创建的 WebSocket 连接对象。如果连接创建失败，则返回 nil。
//
// 错误处理:
//   - 如果升级过程中发生错误，日志中记录错误信息，并返回 nil。
func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	// 初始化响应头，用于设置 WebSocket 协议升级的响应。
	var responseHeader http.Header

	// 如果请求指定了 Sec-WebSocket-Protocol 头，则在响应中进行相应设置。
	if protocol := r.Header.Get("Sec-WebSocket-Protocol"); protocol != "" {
		responseHeader = http.Header{
			"Sec-WebSocket-Protocol": []string{protocol},
		}
	}

	// 使用服务器的升级器将 HTTP 连接升级为 WebSocket 连接，并应用响应头。
	// 如果升级失败，记录错误并返回 nil。
	c, err := s.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		s.Errorf("Error upgrading connection: %s", err)
		return nil
	}

	// 初始化 Conn 实例，并设置相关属性。
	// 包括 WebSocket 连接、服务器实例、连接空闲时间、最大空闲时间、消息队列等。
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

	// 启动后台协程执行心跳检测，以保持连接的活跃状态。
	go conn.keepalive()
	return conn
}
