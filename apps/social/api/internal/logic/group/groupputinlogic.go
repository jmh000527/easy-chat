package group

import (
	"context"
	"easy-chat/apps/im/rpc/imclient"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/ctxdata"

	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutIn 处理用户加入群组的请求（被邀请直接加入，自行申请则需要由 GroupPutInHandle 处理）
//
// 功能描述:
//   - 该方法首先调用外部服务将用户加入群组。
//   - 如果加入成功，则创建一个新的群组会话。
//   - 如果操作失败，则返回错误信息。
//
// 参数:
//   - req: `*types.GroupPutInRep` 类型，包含用户请求加入的群组ID、请求消息、请求时间和加入来源等信息。
//
// 返回值:
//   - `*types.GroupPutInResp`: 处理结果的响应。
//   - `error`: 处理过程中发生的错误。
func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInRep) (resp *types.GroupPutInResp, err error) {
	// 获取当前用户的ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用外部服务，将用户请求加入群组
	res, err := l.svcCtx.Social.GroupPutin(l.ctx, &socialclient.GroupPutinReq{
		GroupId:    req.GroupId,           // 群组ID
		ReqId:      uid,                   // 用户ID
		ReqMsg:     req.ReqMsg,            // 请求消息
		ReqTime:    req.ReqTime,           // 请求时间
		JoinSource: int32(req.JoinSource), // 加入来源
	})
	if err != nil {
		// 如果加入群组操作失败，返回错误信息
		return nil, err
	}

	// 检查返回的群组ID是否有效
	if res.GroupId == "" {
		// 如果群组ID无效，则返回错误
		return nil, err
	}

	// 创建与群组的会话
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,                            // 发送者ID（用户ID）
		RecvId:   res.GroupId,                    // 接收者ID（群组ID）
		ChatType: int32(constants.GroupChatType), // 聊天类型，群组聊天
	})
	if err != nil {
		// 如果创建会话操作失败，返回错误信息
		return nil, err
	}

	// 返回成功响应（此处未返回具体数据）
	return nil, nil
}
