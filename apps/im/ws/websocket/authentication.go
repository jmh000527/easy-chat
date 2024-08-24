package websocket

import (
	"fmt"
	"net/http"
	"time"
)

// Authentication 定义了一个认证接口，包含认证方法和获取用户ID的方法。
type Authentication interface {
	// Authenticate 尝试对当前请求进行认证，返回认证结果。
	// 参数:
	//   w: http.ResponseWriter，用于向客户端发送响应。
	//   r: *http.Request，当前的HTTP请求。
	// 返回值:
	//   bool: 表示认证是否成功。
	Authenticate(w http.ResponseWriter, r *http.Request) bool

	// UserId 从请求中提取并返回用户ID。
	// 参数:
	//   r: *http.Request，当前的HTTP请求。
	// 返回值:
	//   string: 用户ID，如果无法获取，则返回默认值。
	UserId(r *http.Request) string
}

// webSocketAuthentication 是Authentication接口的一个实现，用于WebSocket的认证。
type webSocketAuthentication struct{}

// Authenticate 实现了Authentication接口的Authenticate方法。
// 对于WebSocket认证，假设总是成功的，因此返回true。
func (a *webSocketAuthentication) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}

// UserId 实现了Authentication接口的UserId方法。
// 从请求的查询参数中获取用户ID，如果不存在，则生成一个基于当前时间戳的唯一ID。
//
// 参数:
//
//	r: *http.Request，当前的HTTP请求。
//
// 返回值:
//
//	string: 用户ID。
func (a *webSocketAuthentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}

	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
