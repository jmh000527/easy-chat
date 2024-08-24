package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"errors"
	"github.com/jinzhu/copier"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = errors.New("此用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserInfo 获取用户信息。
//
// 功能描述:
//   - 根据给定的用户ID从数据库中获取用户信息。
//   - 如果用户不存在，返回相应的错误。
//   - 将获取到的用户信息填充到响应对象中并返回。
//
// 参数:
//   - in: *user.GetUserInfoReq
//     包含要获取信息的用户ID。
//
// 返回值:
//   - *user.GetUserInfoResp: 包含用户信息的响应对象。
//   - error: 如果发生错误，返回相应的错误信息。
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// 通过用户ID从数据库中查找用户信息
	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		// 处理用户不存在的情况
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 将数据库中的用户实体复制到响应对象中
	var resp user.UserEntity
	if err := copier.Copy(&resp, userEntity); err != nil {
		return nil, err
	}

	// 返回用户信息
	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
