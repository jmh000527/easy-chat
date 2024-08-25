package websocket

import "time"

// FrameType 表示 WebSocket 消息的帧类型。
type FrameType uint8

const (
	FrameData  FrameType = 0x0 // 数据帧
	FramePing  FrameType = 0x1 // Ping 帧
	FrameAck   FrameType = 0x2 // Ack 帧
	FrameNoAck FrameType = 0x3 // 无 Ack 帧
	FrameErr   FrameType = 0x9 // 错误帧

	// 其他可能的帧类型（已注释）
	//FrameHeaders      FrameType = 0x1
	//FramePriority     FrameType = 0x2
	//FrameRSTStream    FrameType = 0x3
	//FrameSettings     FrameType = 0x4
	//FramePushPromise  FrameType = 0x5
	//FrameGoAway       FrameType = 0x7
	//FrameWindowUpdate FrameType = 0x8
	//FrameContinuation FrameType = 0x9
)

// Message 表示 WebSocket 消息。
type Message struct {
	FrameType `json:"frameType"` // 帧类型
	Id        string             `json:"id"`       // 消息 ID
	AckSeq    int                `json:"ackSeq"`   // Ack序列号
	AckTime   time.Time          `json:"ackTime"`  // 确认时间
	ErrCount  int                `json:"errCount"` // 错误计数
	Method    string             `json:"method"`   // 方法
	FormId    string             `json:"formId"`   // 来源 ID
	Data      interface{}        `json:"data"`     // 数据（使用空接口）
}

// NewMessage 创建一个新的数据消息。
//
// 该函数用于创建一个包含数据的消息对象。消息的类型被设置为 `FrameData`，
// `FormId` 是消息的发起者标识，`Data` 是消息的实际内容。
//
// 参数:
//   - formId: 消息发起者的标识符，用于标识发送消息的来源。
//   - data: 消息的实际内容，可以是任何类型的值。
//
// 返回值:
//   - *Message: 返回创建好的数据消息对象。
func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}

// NewErrMessage 创建一个新的错误消息。
//
// 该函数用于创建一个包含错误信息的消息对象。消息的类型被设置为 `FrameErr`，
// `Data` 字段包含了错误信息的字符串表示。
//
// 参数:
//   - err: 错误对象，将其错误信息转换为字符串用于消息内容。
//
// 返回值:
//   - *Message: 返回创建好的错误消息对象。
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
