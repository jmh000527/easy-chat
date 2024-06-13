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
	ErrPhoneNotRegistered = xerr.New(xerr.ServerCommonError, "手机号没有注册")
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

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 验证用户是否注册过
	// 调用FindOneByPhoneNumber方法通过手机号查找用户
	userEntity, err := l.svcCtx.UsersModel.FindOneByPhoneNumber(l.ctx, in.Phone)
	if err != nil {
		// 如果用户不存在，返回ErrPhoneNotRegistered错误
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegistered)
		}
		// 其他错误，返回数据库错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err: %v, req %v", err, in.Phone)
	}

	// 密码验证
	// 调用ValidatePasswordHash方法验证输入密码是否正确
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// 生成token
	// 获取当前时间的Unix时间戳
	now := time.Now().Unix()
	// 调用GetJwtToken方法生成JWT令牌
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err: %v", err)
	}

	// 返回登录响应，其中包含用户ID、生成的token、过期时间和用户信息
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
