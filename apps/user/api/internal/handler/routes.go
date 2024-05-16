// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	user "easy-chat/apps/user/api/internal/handler/user"
	"easy-chat/apps/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user",
				Handler: user.DetailHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/v1/user"),
	)
}
