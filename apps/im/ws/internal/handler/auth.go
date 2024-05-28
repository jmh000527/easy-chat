package handler

import (
	"context"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/pkg/ctxdata"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
)

type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser
	logx.Logger
}

func (j *JwtAuth) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	if t := r.Header.Get("sec-websocket-protocol"); t != "" {
		r.Header.Set("Authorization", t)
	}
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("parse token err: %v", err)
		return false
	}

	if !tok.Valid {
		j.Errorf("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		j.Errorf("invalid token")
		return false
	}

	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))
	return true
}

func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUId(r.Context())
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}
