package user

import (
	"context"
	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/apps/user/api/internal/types"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/ctxdata"
	"github.com/jinzhu/copier"

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

// Detail 通过用户ID获取用户的详细信息。
//
// 功能描述:
//   - 从上下文中获取用户ID。
//   - 使用该用户ID调用 svcCtx 的 User.GetUserInfo 方法获取用户详细信息。
//   - 将获取到的用户信息转换为 types.User 类型。
//   - 构建并返回包含用户详细信息的响应对象。
//
// 参数:
//   - req: *types.UserInfoReq
//     请求参数，包含需要获取用户详细信息的用户ID（从上下文中获取）.
//
// 返回值:
//   - *types.UserInfoResp: 包含用户详细信息的响应对象。
//   - error: 如果获取用户信息或处理过程中发生错误，则返回相应的错误信息。
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

	// 构建并返回用户详细信息响应
	return &types.UserInfoResp{Info: res}, nil
}
