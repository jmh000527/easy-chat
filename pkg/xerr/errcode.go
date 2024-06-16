package xerr

// 定义了一组常量，用于表示服务端可能遇到的错误码
const (
	// ServerCommonError 表示服务端通用错误，用于指示服务端发生了一般性错误
	ServerCommonError = 100001
	// RequestParamError 表示请求参数错误，用于指示客户端发送的请求参数有误
	RequestParamError = 100002
	// DbError 表示数据库操作错误，用于指示服务端在进行数据库操作时发生了错误
	DbError = 100003
)
