package msgtransfer

import (
	"context"
	"easy-chat/apps/task/mq/internal/svc"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	svc *svc.ServiceContext
	logx.Logger
}

func (m MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("key:", key, "value:", value)
	return nil
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}
