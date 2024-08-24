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

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutInHandle 处理群组加入请求
//
// 功能描述:
//   - 处理群组加入请求，根据请求的处理结果更新请求状态
//   - 如果请求被批准，则建立一个新的用户和群组的聊天会话
//
// 参数:
//   - req: `*types.GroupPutInHandleReq` 类型，包含处理群组请求所需的信息，包括群组请求ID、群组ID、处理结果等
//
// 返回值:
//   - `*types.GroupPutInHandleResp`: 空响应对象
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	// 获取当前用户ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用社交服务处理群组请求
	res, err := l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,   // 群组请求ID
		GroupId:      req.GroupId,      // 群组ID
		HandleUid:    uid,              // 处理请求的用户ID
		HandleResult: req.HandleResult, // 处理结果（通过或拒绝）
	})
	if err != nil {
		return nil, err
	}

	// 如果客户端期望的处理结果不是通过，直接返回响应和错误
	if constants.HandlerResult(req.HandleResult) != constants.PassHandlerResult {
		return
	}

	// 如果群组ID为空，则返回错误
	if res.GroupId == "" {
		return nil, err
	}

	// 建立用户与群组之间的聊天会话
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,                            // 发送者ID（当前用户ID）
		RecvId:   res.GroupId,                    // 接收者ID（群组ID）
		ChatType: int32(constants.GroupChatType), // 聊天类型（群组聊天）
	})

	// 返回空响应和错误
	return nil, err
}
