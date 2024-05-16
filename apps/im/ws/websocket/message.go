package websocket

type Message struct {
	Method string      `json:"method"`
	FormId string      `json:"formId"`
	Data   interface{} `json:"data"`
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FormId: formId,
		Data:   data,
	}
}
