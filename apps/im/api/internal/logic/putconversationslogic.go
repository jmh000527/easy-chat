package logic

import (
	"context"
	"easy-chat/apps/im/rpc/imclient"
	"easy-chat/pkg/ctxdata"
	"github.com/jinzhu/copier"

	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PutConversations 更新用户的会话信息。
//
// 该方法将请求中的会话信息更新到数据库中，
// 其中包括会话的显示状态、已读记录和其他相关信息。
//
// 参数:
//   - req: 请求对象，包含用户 ID 和需要更新的会话信息。
//
// 返回值:
//   - *types.PutConversationsResp: 响应对象，表示更新操作的结果。
//   - error: 如果在更新过程中发生错误，返回具体的错误信息；成功时返回 nil。
func (l *PutConversationsLogic) PutConversations(req *types.PutConversationsReq) (resp *types.PutConversationsResp, err error) {
	// 从上下文中获取用户 ID
	uid := ctxdata.GetUId(l.ctx)

	// 将请求中的会话信息复制到一个新的映射中
	var conversationList map[string]*imclient.Conversation
	copier.Copy(&conversationList, req.ConversationList)

	// 调用服务上下文中的方法进行会话更新
	_, err = l.svcCtx.PutConversations(l.ctx, &imclient.PutConversationsReq{
		UserId:           uid,
		ConversationList: conversationList,
	})

	// 返回更新结果
	return
}
