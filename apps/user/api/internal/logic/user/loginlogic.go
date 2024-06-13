package user

import (
	"context"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/constants"
	"github.com/jinzhu/copier"

	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 处理用户登录请求
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 调用 svcCtx 的 User.Login 方法进行用户登录
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// 将 user.LoginResp 转换为 types.LoginResp
	var res types.LoginResp
	if err := copier.Copy(&res, loginResp); err != nil {
		return nil, err
	}

	// 处理登录后的业务，将用户标记为在线用户
	err = l.svcCtx.Redis.HsetCtx(l.ctx, constants.RedisOnlineUser, loginResp.Id, "1")
	if err != nil {
		return nil, err
	}

	return &res, nil
}
