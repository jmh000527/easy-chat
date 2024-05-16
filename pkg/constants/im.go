package constants

type MType int
type ChatType int

const (
	TextMType MType = iota
)

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)
