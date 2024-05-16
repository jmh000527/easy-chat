package constants

type MType int
type ChatType int
type ContentType int

const (
	TextMType MType = iota
)

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

const (
	ContentChatMsg ContentType = iota
	ContentMakeRead
)
