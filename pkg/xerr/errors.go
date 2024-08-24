package xerr

import "github.com/zeromicro/x/errors"

// New 创建一个带有错误代码和错误消息的新错误。
//
// 功能描述:
//   - 根据指定的错误代码和错误消息创建一个新的错误对象。
//   - 错误代码用于标识错误的类型，错误消息提供对错误的详细描述。
//
// 参数:
//   - code: int
//     错误代码，用于标识错误的类型。
//   - msg: string
//     错误消息，对错误进行详细描述。
//
// 返回值:
//   - error: 返回一个封装了错误代码和错误消息的新错误。
func New(code int, msg string) error {
	return errors.New(code, msg)
}

// NewMsg 创建一个带有预定义错误代码和自定义错误消息的新错误。
//
// 功能描述:
//   - 创建一个新的错误对象，使用预定义的错误代码和用户提供的自定义错误消息。
//   - 预定义错误代码用于标识错误类型，自定义消息提供对错误的具体描述。
//
// 参数:
//   - msg: string
//     自定义错误消息，对错误进行详细描述。
//
// 返回值:
//   - error: 返回一个封装了预定义错误代码和自定义错误消息的新错误。
func NewMsg(msg string) error {
	return errors.New(ServerCommonError, msg)
}

// NewDBErr 创建一个表示数据库错误的新错误。
//
// 功能描述:
//   - 创建一个新的错误对象，用于表示数据库相关的错误。
//   - 使用预定义的数据库错误代码，并通过 ErrMsg 获取详细的错误消息。
//
// 返回值:
//   - error: 返回一个封装了数据库错误代码和详细消息的新错误。
func NewDBErr() error {
	return errors.New(DbError, ErrMsg(DbError))
}

// NewInternalErr 创建一个表示内部服务器错误的新错误。
//
// 功能描述:
//   - 创建一个新的错误对象，用于表示内部服务器错误。
//   - 使用预定义的内部服务器错误代码，并通过 ErrMsg 获取详细的错误消息。
//
// 返回值:
//   - error: 返回一个封装了内部服务器错误代码和详细消息的新错误。
func NewInternalErr() error {
	return errors.New(ServerCommonError, ErrMsg(ServerCommonError))
}
