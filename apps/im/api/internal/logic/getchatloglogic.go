package logic

import (
	"context"
	"easy-chat/apps/im/rpc/imclient"
	"github.com/jinzhu/copier"

	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetChatLog 获取聊天记录 chatlog。
//
// 该方法调用服务上下文中的 GetChatLog 方法来从数据源中获取聊天记录。
// 根据请求中的参数，方法会查询特定会话的聊天记录，并将结果返回给调用方。
//
// 参数:
//   - req: 请求对象，包含查询聊天记录所需的所有信息。
//
// 返回值:
//   - *types.ChatLogResp: 查询结果的响应对象，包含聊天记录的列表。
//   - error: 如果在查询过程中发生错误，则返回具体的错误信息。成功时返回 nil。
func (l *GetChatLogLogic) GetChatLog(req *types.ChatLogReq) (resp *types.ChatLogResp, err error) {
	// 调用服务上下文中的 GetChatLog 方法获取聊天记录
	data, err := l.svcCtx.GetChatLog(l.ctx, &imclient.GetChatLogReq{
		ConversationId: req.ConversationId,
		StartSendTime:  req.StartSendTime,
		EndSendTime:    req.EndSendTime,
		Count:          req.Count,
	})
	if err != nil {
		// 如果获取聊天记录时发生错误，返回 nil 和错误信息
		return nil, err
	}

	var res types.ChatLogResp
	// 将获取到的数据复制到响应对象中
	copier.Copy(&res, &data)

	// 返回包含聊天记录的响应对象和 nil 错误
	return &res, nil
}
