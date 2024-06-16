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

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 验证用户是否注册过
	// 调用FindOneByPhoneNumber方法通过手机号查找用户，如果发生错误且错误不是models.ErrNotFound，返回错误
	userEntity, err := l.svcCtx.UsersModel.FindOneByPhoneNumber(l.ctx, in.Phone)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	// 如果用户已经存在，返回ErrPhoneIsRegistered错误
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

	// 如果输入的密码不为空，生成密码哈希并赋值给用户实体
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
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

	// 生成token
	// 获取当前时间的Unix时间戳
	now := time.Now().Unix()
	// 调用GetJwtToken方法生成JWT令牌
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	// 返回注册响应，其中包含生成的token和过期时间
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
