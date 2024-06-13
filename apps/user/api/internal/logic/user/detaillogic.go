package user

import (
	"context"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/ctxdata"
	"github.com/jinzhu/copier"

	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Detail 获取用户详细信息
func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// 从上下文中获取用户ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用 svcCtx 的 User.GetUserInfo 方法获取用户信息
	userInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}

	// 将 user.UserInfoResp.User 转换为 types.User
	var res types.User
	if err := copier.Copy(&res, userInfoResp.User); err != nil {
		return nil, err
	}

	// 返回用户详细信息
	return &types.UserInfoResp{Info: res}, nil
}
