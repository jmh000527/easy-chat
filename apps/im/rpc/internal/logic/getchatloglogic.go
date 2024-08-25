package logic

import (
	"context"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"

	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetChatLog 获取会话记录。
//
// 该方法根据请求中的参数从数据库中获取聊天记录。根据是否提供了 msgId，
// 方法会选择不同的查询方式：如果 msgId 不为空，则直接查询该消息记录；
// 如果 msgId 为空，则根据时间段进行查询。查询的结果会按照时间排序，并返回符合条件的聊天记录。
//
// 参数:
//   - in: 请求对象，包含查询条件。
//
// 返回值:
//
//   - *im.GetChatLogResp: 查询结果的响应对象。
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// 如果请求中提供了 msgId，直接查询该消息记录
	if in.MsgId != "" {
		chatLog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			// 如果查询过程中发生错误，返回包装后的错误信息
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chatlog by msgId %s failed", in.MsgId)
		}
		// 构造并返回响应对象，包含查询到的单条聊天记录
		return &im.GetChatLogResp{
			List: []*im.ChatLog{
				{
					Id:             chatLog.ID.Hex(),
					ConversationId: chatLog.ConversationId,
					SendId:         chatLog.SendId,
					RecvId:         chatLog.RecvId,
					MsgType:        int32(chatLog.MsgType),
					MsgContent:     chatLog.MsgContent,
					ChatType:       int32(chatLog.ChatType),
					SendTime:       chatLog.SendTime,
					ReadRecords:    chatLog.ReadRecords,
				},
			},
		}, nil
	}

	// 如果没有提供 msgId，基于时间范围查询聊天记录
	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		// 如果查询过程中发生错误，返回包装后的错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find chatLog list by SendTime failed, err: %v req: %v", err.Error(), in)
	}

	// 构造查询结果列表
	res := make([]*im.ChatLog, 0, len(data))
	for _, v := range data {
		res = append(res, &im.ChatLog{
			Id:             v.ID.Hex(),
			ConversationId: v.ConversationId,
			SendId:         v.SendId,
			RecvId:         v.RecvId,
			MsgType:        int32(v.MsgType),
			MsgContent:     v.MsgContent,
			ChatType:       int32(v.ChatType),
			SendTime:       v.SendTime,
			ReadRecords:    v.ReadRecords,
		})
	}
	// 返回包含聊天记录列表的响应对象
	return &im.GetChatLogResp{
		List: res,
	}, nil
}
