package logic

import (
	"context"
	"database/sql"
	"easy-chat/apps/user/models"
	"easy-chat/pkg/ctxdata"
	"easy-chat/pkg/encrypt"
	"easy-chat/pkg/wuid"
	"errors"
	"time"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// ErrPhoneIsRegistered 表示尝试注册一个已经注册过的手机号的错误。
	ErrPhoneIsRegistered = errors.New("手机号已经注册过")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 处理用户注册请求。
//
// 功能描述:
//   - 验证用户是否已经注册。
//   - 如果用户未注册，创建一个新的用户并将其信息存入数据库。
//   - 为用户生成 JWT 令牌并返回。
//
// 参数:
//   - in: 包含用户注册信息的请求结构体。
//
// 返回值:
//   - *user.RegisterResp: 包含生成的 JWT 令牌和过期时间的响应结构体。
//   - error: 如果注册过程中出现错误，则返回相应的错误信息。
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 验证用户是否注册过
	// 调用 FindOneByPhoneNumber 方法通过手机号查找用户，
	// 如果发生错误且错误不是 models.ErrNotFound，返回错误
	userEntity, err := l.svcCtx.UsersModel.FindOneByPhoneNumber(l.ctx, in.Phone)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	// 如果用户已经存在，返回 ErrPhoneIsRegistered 错误
	if userEntity != nil {
		return nil, ErrPhoneIsRegistered
	}

	// 创建一个新的用户实体，并填充其数据
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.Datasource), // 生成唯一用户ID
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	// 如果输入的密码不为空
	if len(in.Password) > 0 {
		// 调用 GenPasswordHash 方法生成密码哈希
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		// 将密码哈希赋值给用户实体的 Password 字段
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	// 将用户实体插入数据库
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// 生成 token
	// 获取当前时间的 Unix 时间戳
	now := time.Now().Unix()
	// 调用 GetJwtToken 方法生成 JWT 令牌
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	// 返回注册响应，其中包含生成的 token 和过期时间
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
