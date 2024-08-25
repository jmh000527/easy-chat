package logic

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"

	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PutConversations 更新会话信息。
//
// 该方法更新指定用户的会话列表，将新的会话信息保存到数据库中。
// 如果用户原本有会话数据，将会合并新数据；如果没有，将会创建新的会话记录。
//
// 参数:
//   - in: 请求对象，包含需要更新的会话信息。
//
// 返回值:
//   - *im.PutConversationsResp: 响应对象，表示更新操作的结果。
//   - error: 如果在更新过程中发生错误，返回具体的错误信息；成功时返回 nil。
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// 查询用户的会话列表
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		// 查询会话列表失败，返回 nil 和错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversations by user id failed, uid: %s, err: %v", in.UserId, err)
	}

	// 如果会话列表为空，则初始化一个空的会话列表
	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}

	for s, conversation := range in.ConversationList {
		// 获取用户原本读取的会话消息量
		var oldTotal int
		if data.ConversationList[s] != nil {
			oldTotal = data.ConversationList[s].Total
		}

		// 更新会话信息
		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			Total:          int(conversation.Read) + oldTotal, // 新的已读记录量 + 原本读取的会话消息量
			Seq:            conversation.Seq,
		}
	}

	// 将更新后的会话列表保存到数据库
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		// 更新会话列表失败，返回 nil 和错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "update conversations failed, uid: %s, err: %v", in.UserId, err)
	}

	// 更新成功，返回空响应和 nil 错误
	return &im.PutConversationsResp{}, nil
}
