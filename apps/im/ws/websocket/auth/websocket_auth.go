package auth

import (
	"fmt"
	"net/http"
	"time"
)

// WebSocketAuth 是Authentication接口的一个实现，用于WebSocket的认证。
type WebSocketAuth struct{}

// Authenticate 实现了Authentication接口的Authenticate方法。
// 对于WebSocket认证，假设总是成功的，因此返回true。
func (a *WebSocketAuth) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}

// UserId 根据请求从 URL 查询参数中提取用户 ID。
//
// 参数:
//   - r: *http.Request
//     HTTP 请求对象，包含客户端的请求信息。
//
// 返回值:
//   - string: 提取到的用户 ID。如果 URL 查询参数中存在 "userId" 参数，则返回其值；
//     如果 "userId" 参数不存在，则返回当前时间的毫秒级时间戳作为用户 ID。
//
// 功能说明:
//
//	`UserId` 函数从 HTTP 请求的 URL 查询参数中提取用户 ID。如果查询参数中包含 "userId" 参数，
//	则返回该参数的值；如果没有该参数，则返回当前时间的毫秒级时间戳，
//	以确保生成的用户 ID 是唯一的。
func (a *WebSocketAuth) UserId(r *http.Request) string {
	// 提取URL中的查询参数
	query := r.URL.Query()
	// 检查查询参数是否存在且包含"userId"
	if query != nil && query["userId"] != nil {
		// 如果存在"userId"参数，返回其值
		return fmt.Sprintf("%v", query["userId"])
	}

	// 如果不存在"userId"参数，生成并返回当前时间的Unix毫秒时间戳作为用户ID
	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
