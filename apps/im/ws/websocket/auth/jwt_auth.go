package auth

import (
	"context"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/pkg/ctxdata"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
)

// JwtAuth 用于处理基于 JWT 的身份认证。
//
// 该结构体包含了服务上下文、令牌解析器和日志记录器，用于验证 WebSocket 请求的 JWT 令牌，
// 并从中提取用户标识符。
type JwtAuth struct {
	svc         *svc.ServiceContext // 服务上下文，提供服务配置和依赖。
	parser      *token.TokenParser  // JWT 令牌解析器，用于解析和验证 JWT 令牌。
	logx.Logger                     // 日志记录器，用于记录日志信息。
}

// NewJwtAuth 创建一个新的 JwtAuth 实例。
//
// 该方法用于初始化 JwtAuth 结构体，并返回一个新的实例。
//
// 参数:
//   - svc: 服务上下文，提供服务配置和依赖。
//
// 返回值:
//   - *JwtAuth: 返回新创建的 JwtAuth 实例。
func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

// Authenticate 验证请求的 JWT 令牌。
//
// 该方法从请求头中提取 JWT 令牌，并使用解析器进行验证。如果令牌有效，
// 将用户标识符注入到请求的上下文中。
//
// 参数:
//   - w: HTTP 响应写入器，用于写入响应。
//   - r: HTTP 请求对象，其中包含了请求的所有信息。
//
// 返回值:
//   - bool: 如果验证成功返回 true，否则返回 false。
func (j *JwtAuth) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	// 从 WebSocket 协议头中提取 JWT 令牌，并将其设置到 Authorization 头中。
	if t := r.Header.Get("sec-websocket-protocol"); t != "" {
		r.Header.Set("Authorization", t)
	}

	// 解析并验证 JWT 令牌。
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("parse token err: %v", err)
		return false
	}

	// 检查令牌的有效性。
	if !tok.Valid {
		j.Errorf("invalid token")
		return false
	}

	// 提取令牌中的声明 (claims) 并验证其类型。
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		j.Errorf("invalid token")
		return false
	}

	// 将用户标识符注入到请求上下文中。
	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))
	return true
}

// UserId 从请求的上下文中获取用户标识符。
//
// 该方法从请求上下文中提取之前注入的用户标识符。
//
// 参数:
//   - r: HTTP 请求对象，其中包含了请求的所有信息。
//
// 返回值:
//   - string: 返回提取到的用户标识符。
func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUId(r.Context())
}
