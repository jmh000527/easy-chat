package handler

import (
	"easy-chat/apps/task/mq/internal/handler/msgtransfer"
	"easy-chat/apps/task/mq/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Listen struct {
	svc *svc.ServiceContext
}

func (l *Listen) Services() []service.Service {
	return []service.Service{
		// todo: 此处可以加载多个消费者
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, msgtransfer.NewMsgChatTransfer(l.svc)),
	}
}

func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{
		svc: svc,
	}
}
