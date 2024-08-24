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

// Login 处理用户登录请求。
//
// 功能描述:
//   - 调用 svcCtx 的 User.Login 方法进行用户登录。
//   - 将 user.LoginResp 转换为 types.LoginResp。
//   - 将用户ID和在线状态"1"存储到Redis的hash中，标记用户为在线。
//   - 使用Redis来管理在线用户，因为Redis的高并发读写性能和键值对存储特性适合此类场景。
//
// 参数:
//   - req: *types.LoginReq
//     登录请求的输入参数，包含手机号和密码。
//
// 返回值:
//   - *types.LoginResp: 包含登录成功后的用户信息和生成的token。
//   - error: 如果登录验证、数据转换或Redis操作中出现错误，则返回相应的错误信息。
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 调用 svcCtx 的 User.Login 方法进行用户登录
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		// 如果登录验证失败，返回错误。
		return nil, err
	}

	// 将 user.LoginResp 转换为 types.LoginResp
	var res types.LoginResp
	if err := copier.Copy(&res, loginResp); err != nil {
		// 如果复制过程中发生错误，返回错误。
		return nil, err
	}

	// 将用户ID和在线状态"1"存储到Redis的hash中，标记用户为在线。
	// 这里使用Redis来管理在线用户，是因为Redis的高并发读写性能和键值对存储特性适合此类场景。
	err = l.svcCtx.Redis.HsetCtx(l.ctx, constants.RedisOnlineUser, loginResp.Id, "1")
	if err != nil {
		// 如果设置Redis中用户在线状态失败，返回错误。
		return nil, err
	}

	// 登录成功，返回复制后的登录响应。
	return &res, nil
}
