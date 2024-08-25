package logic

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/wuid"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetUpUserConversation 建立会话: 群聊, 私聊
//
// 该方法负责创建或设置用户会话，支持单聊和群聊两种类型。
// 根据聊天类型，创建或更新相应的会话记录。
//
// 参数:
// - in: 包含建立会话所需的请求参数，主要包括发送者ID、接收者ID和聊天类型。
//
// 返回:
// - *im.SetUpUserConversationResp: 返回会话设置的响应结构体。
// - error: 发生的错误（如果有的话），返回nil表示操作成功。
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	var resp im.SetUpUserConversationResp

	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		// 生成会话ID
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		// 查询是否有会话记录
		conversation, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			// 没有建立过会话，建立会话
			if errors.Is(err, immodels.ErrNotFound) {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})
				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "insert conversation err: %v", err)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation err: %v", err)
			}
		} else if conversation != nil {
			// 会话已经建立过，不需要重复建立
			return &resp, nil
		}
		// 建立发送者和接收者的会话
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}
		// 接收者的会话设置为不展示
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}
	case constants.GroupChatType:
		// 对于群聊，会话ID就是群组ID
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.GroupChatType, true)
		if err != nil {
			return nil, err
		}
	}

	return &im.SetUpUserConversationResp{}, nil
}

// setUpUserConversation 设置用户会话
//
// 该方法管理用户的会话列表，添加或更新会话记录。
// 根据会话ID和用户ID更新会话列表，并确保会话记录存在。
//
// 参数:
// - conversationId: 会话ID，用于标识不同的会话。
// - userId: 用户ID，需要设置会话的用户ID。
// - recvId: 接收者ID，用于标识接收者的会话记录。
// - chatType: 聊天类型，表示是单聊还是群聊。
// - isShow: 是否在用户会话列表中显示该会话。
//
// 返回:
// - error: 发生的错误（如果有的话），返回nil表示操作成功。
func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string, chatType constants.ChatType, isShow bool) error {
	// 获取用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if errors.Is(err, immodels.ErrNotFound) {
			// 如果用户没有会话列表，则创建新的会话列表
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}
		} else {
			return errors.Wrapf(xerr.NewDBErr(), "find by user id err: %v", err)
		}
	}
	// 检查会话是否已存在
	// 如果存在，则直接返回
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}
	// 添加新的会话记录
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       chatType,
		IsShow:         isShow,
	}
	// 更新会话列表
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "insert conversation err: %v", err)
	}
	return nil
}
