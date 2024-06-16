package xerr

import "github.com/zeromicro/x/errors"

// New 创建一个带有错误代码和错误消息的新错误。
// code: 错误代码，用于标识错误的类型。
// msg: 错误消息，对错误进行详细描述。
// 返回一个封装了错误代码和错误消息的新错误。
func New(code int, msg string) error {
	return errors.New(code, msg)
}

// NewMsg 创建一个带有预定义错误代码和自定义错误消息的新错误。
// msg: 自定义错误消息，对错误进行详细描述。
// 返回一个封装了预定义错误代码和自定义错误消息的新错误。
func NewMsg(msg string) error {
	return errors.New(ServerCommonError, msg)
}

// NewDBErr 创建一个表示数据库错误的新错误。
// 该方法通过调用ErrMsg获取数据库错误的详细消息，然后封装成错误返回。
// 返回一个封装了数据库错误代码和详细消息的新错误。
func NewDBErr() error {
	return errors.New(DbError, ErrMsg(DbError))
}

// NewInternalErr 创建一个表示内部服务器错误的新错误。
// 该方法通过调用ErrMsg获取内部服务器错误的详细消息，然后封装成错误返回。
// 返回一个封装了内部服务器错误代码和详细消息的新错误。
func NewInternalErr() error {
	return errors.New(ServerCommonError, ErrMsg(ServerCommonError))
}
