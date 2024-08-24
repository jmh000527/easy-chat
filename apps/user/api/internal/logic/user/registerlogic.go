package user

import (
	"context"
	"easy-chat/apps/user/rpc/user"
	"github.com/jinzhu/copier"

	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 实现用户注册逻辑。
//
// 功能描述:
//   - 接收一个注册请求类型 req，并返回一个注册响应类型 resp 和可能的错误。
//   - 该方法主要负责将注册请求转发给服务层处理，并将处理结果转换为统一的响应格式。
//
// 参数:
//   - req: *types.RegisterReq
//     用户注册请求的输入参数，包括手机号、昵称、密码、头像和性别。
//
// 返回值:
//   - *types.RegisterResp: 包含注册成功后的响应数据。
//   - error: 如果在服务层处理过程中出现错误，则返回相应的错误信息。
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 调用服务层的 Register 方法，传入注册请求信息
	registerResp, err := l.svcCtx.Register(l.ctx, &user.RegisterReq{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	// 如果服务层注册过程中出现错误，直接返回错误
	if err != nil {
		return nil, err
	}

	// 使用 copier 将服务层的注册响应拷贝到业务层的注册响应结构体中
	// 这一步是为了将底层实现的响应格式转换为上层统一的响应格式
	var res types.RegisterResp
	err = copier.Copy(&res, registerResp)
	// 如果拷贝过程中出现错误，返回错误
	if err != nil {
		return nil, err
	}

	// 返回拷贝后的注册响应
	return &res, nil
}
