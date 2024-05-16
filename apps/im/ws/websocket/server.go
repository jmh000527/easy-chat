package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

type Server struct {
	routes   map[string]HandlerFunc
	addr     string
	patten   string
	upgrader websocket.Upgrader
	logx.Logger

	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn
	sync.RWMutex

	authentication Authentication
}

func (s *Server) SendByUserIds(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
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

func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*websocket.Conn) []string {

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

func (s *Server) Close(conn *websocket.Conn) {
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

func (s *Server) addConn(conn *websocket.Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// handlerConn 根据连接对象进行任务处理
func (s *Server) handlerConn(conn *websocket.Conn) {
	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket read error: %v", err)
			s.Close(conn)
			return
		}

		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("websocket unmarshal error: %v", err)
			s.Close(conn)
			return
		}

		// 根据请求方法分发路由执行
		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, &message)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在执行的方法: %v", message.Method)))
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
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("server handler ws upgrade err: %v", err)
		return
	}

	//鉴权
	if !s.authentication.Authenticate(w, r) {
		conn.WriteMessage(websocket.CloseMessage, []byte(fmt.Sprint("不具备访问权限")))
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
	fmt.Println("停止服务")
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newWebsocketServerOption(opts...)

	return &Server{
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		patten:         opt.patten,
		upgrader:       websocket.Upgrader{},
		Logger:         logx.WithContext(context.Background()),
		connToUser:     make(map[*websocket.Conn]string),
		userToConn:     make(map[string]*websocket.Conn),
		authentication: opt.Authentication,
	}
}
