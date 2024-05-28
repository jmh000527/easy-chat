package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	routes   map[string]HandlerFunc
	addr     string
	patten   string
	opt      *websocketOption
	upgrader websocket.Upgrader
	logx.Logger

	connToUser map[*Conn]string
	userToConn map[string]*Conn
	*threading.TaskRunner
	sync.RWMutex

	authentication Authentication
}

// 判断当前消息是否进行Ack确认
func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck
}

// ACK确认后任务的处理
func (s *Server) handleWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接关闭
			return
		case message := <-conn.message:
			// 根据请求消息分发路由执行
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{
					FrameType: FramePing,
				}, conn)
			case FrameData:
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				}
			}

			// 清除消息确认
			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}

}

// 读取ACK确认
func (s *Server) readAck(conn *Conn) {
	send := func(msg *Message, conn *Conn) error {
		err := s.Send(msg, conn)
		if err == nil {
			return nil
		}

		s.Errorf("message ack OnlyAck send err: %v message: %v", conn, msg)
		conn.messageMu.Lock()
		conn.readMessage[0].errCount++
		conn.messageMu.Unlock()

		tempDelay := time.Duration(200*conn.readMessage[0].errCount) * time.Microsecond
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}
		time.Sleep(tempDelay)
		return err
	}

	for {
		select {
		case <-conn.done:
			// 连接关闭
			s.Infof("close message ack uid: %v ", conn.Uid)
			return
		default:
		}

		// 从队列中读取新的消息
		conn.messageMu.Lock()
		// 当前队列中没有消息
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			// 增加睡眠，让任务更好地切换
			time.Sleep(100 * time.Microsecond)
			continue
		}

		// 读取第一条消息
		message := conn.readMessage[0]
		if message.errCount > s.opt.sendErrCount {
			s.Infof("conn send fail, message: %v, ackType: %v, maxSendErrCount: %v", message, message.errCount, s.opt.sendErrCount)
			conn.messageMu.Unlock()
			// 因为发送消息多次错误，放弃发送消息
			delete(conn.readMessageSeq, message.Id)
			conn.readMessage = conn.readMessage[1:]
			continue
		}

		// 判断Ack方式
		switch s.opt.ack {
		// 只需要一次确认
		case OnlyAck:
			// 直接给客户端回复
			if err := send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn); err != nil {
				continue
			}
			// 进行业务处理
			// 把消息从队列中移除
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()
			conn.message <- message
			s.Infof("message ack OnlyAck send success, mid: %v", message.Id)
		case RigorAck:
			if message.AckSeq == 0 {
				// 还未发送过确认消息
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()
				if err := send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn); err != nil {
					continue
				}

				conn.messageMu.Unlock()
				s.Infof("message ack RigorAck send mid: %v, seq: %v , time: %v", message.Id, message.AckSeq,
					message.ackTime)
				continue
			}

			// 验证
			// 1. 客户端返回结果，再一次确认
			// 获取之前记录的Ack信息，得到客户端的序号
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 客户端进行了确认
				// 删除消息
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid: %v", message.Id)
				continue
			}

			// 2. 客户端没有确认，考虑是否超过了ack的确认时间
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				// 2.1 超过结束确认
				s.Infof("message ack RigorAck timeout: %v ack time: %v", message.Id, message.ackTime)
				// 删除消息序号
				delete(conn.readMessageSeq, message.Id)
				// 删除消息
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}

			// 2.2 未超时，重新发送
			conn.messageMu.Unlock()
			if val > 0 && val > 300*time.Microsecond {
				if err := send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn); err != nil {
					continue
				}
			}
			// 睡眠一定的时间
			time.Sleep(300 * time.Microsecond)
		}
	}
}

func (s *Server) SendByUserIds(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
	conn.Close()
}

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登入过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// handlerConn 根据连接对象进行任务处理
func (s *Server) handlerConn(conn *Conn) {
	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	// 处理任务
	go s.handleWrite(conn)
	// 判断是否开启Ack机制
	if s.isAck(nil) {
		// 进行Ack确认
		go s.readAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket read error: %v", err)
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("websocket unmarshal error: %v", err)
			s.Close(conn)
			return
		}

		// todo: 给客户端回复一个ACK

		// 启动了Ack机制且消息类型需要进行Ack确认
		if s.isAck(&message) {
			// 进行Ack确认
			s.Infof("conn message read ack msg: %v", message)
			conn.appendMsgMq(&message)
		} else {
			// 直接传递消息，不进行Ack确认
			conn.message <- &message
		}
	}
}

func (s *Server) AddRoutes(rs []Route) {
	for _, route := range rs {
		s.routes[route.Method] = route.Handler
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err: %v", r)
		}
	}()
	// 获取一个websocket连接对象
	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	//鉴权
	if !s.authentication.Authenticate(w, r) {
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprint("不具备访问权限")}, conn)
		conn.Close()
		return
	}
	// 记录连接
	s.addConn(conn, r)
	// 根据连接对象获取请求信息
	go s.handlerConn(conn)
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("stop service")
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newWebsocketServerOption(opts...)

	return &Server{
		routes: make(map[string]HandlerFunc),
		addr:   addr,
		patten: opt.patten,
		opt:    &opt,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Logger:         logx.WithContext(context.Background()),
		connToUser:     make(map[*Conn]string),
		userToConn:     make(map[string]*Conn),
		authentication: opt.Authentication,
		TaskRunner:     threading.NewTaskRunner(opt.concurrency),
	}
}
