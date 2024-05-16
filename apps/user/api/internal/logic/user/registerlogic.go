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

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerResp, err := l.svcCtx.Register(l.ctx, &user.RegisterReq{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	if err != nil {
		return nil, err
	}

	var res types.RegisterResp
	copier.Copy(&res, registerResp)

	return &res, nil
}
