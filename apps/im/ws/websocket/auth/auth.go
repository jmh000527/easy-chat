package auth

import "net/http"

// Authentication 定义了一个认证接口，包含认证方法和获取用户ID的方法。
type Authentication interface {
	// Authenticate 尝试对当前请求进行认证，返回认证结果。
	Authenticate(w http.ResponseWriter, r *http.Request) bool

	// UserId 从请求中提取并返回用户ID。
	UserId(r *http.Request) string
}
