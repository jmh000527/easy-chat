package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"easy-chat/pkg/ctxdata"
	"easy-chat/pkg/encrypt"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"time"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegistered = xerr.New(xerr.ServerCommonError, "手机号码没有注册")
	ErrUserPwdError       = xerr.New(xerr.ServerCommonError, "密码不正确")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 处理用户登录请求。
//
// 功能描述:
//   - 验证用户是否已经注册。
//   - 验证用户提供的密码是否正确。
//   - 如果验证通过，为用户生成 JWT 令牌并返回。
//
// 参数:
//   - in: 包含用户登录信息的请求结构体。
//
// 返回值:
//   - *user.LoginResp: 包含用户ID、JWT令牌、过期时间和用户信息的响应结构体。
//   - error: 如果登录过程中出现错误，则返回相应的错误信息。
func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 验证用户是否注册过
	// 调用 FindOneByPhoneNumber 方法通过手机号查找用户
	userEntity, err := l.svcCtx.UsersModel.FindOneByPhoneNumber(l.ctx, in.Phone)
	if err != nil {
		// 如果用户不存在，返回 ErrPhoneNotRegistered 错误
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegistered)
		}
		// 返回数据库错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err: %v, req %v", err, in.Phone)
	}

	// 调用 ValidatePasswordHash 方法验证输入密码是否正确
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// 生成 token
	// 获取当前时间的 Unix 时间戳
	now := time.Now().Unix()
	// 调用 GetJwtToken 方法生成 JWT 令牌
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err: %v", err)
	}

	// 返回登录响应，其中包含用户ID、生成的 token、过期时间和用户信息
	return &user.LoginResp{
		Id:     userEntity.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
		User: &user.UserEntity{
			Id:       userEntity.Id,
			Avatar:   userEntity.Avatar,
			Nickname: userEntity.Nickname,
			Phone:    in.Phone,
			Status:   int32(userEntity.Status.Int64),
			Sex:      int32(userEntity.Sex.Int64),
		},
	}, nil
}
