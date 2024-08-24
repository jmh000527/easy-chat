package xerr

// codeText 是一个错误码与错误信息的映射，用于根据错误码快速查找对应的错误信息。
var codeText = map[int]string{
	ServerCommonError: "服务器异常，稍后再尝试",
	RequestParamError: "请求参数有误",
	DbError:           "数据库繁忙，稍后再尝试",
}

// ErrMsg 根据错误码返回对应的错误信息。
//
// 功能描述:
//   - 根据传入的错误码，从预定义的错误码与错误信息的映射中查找并返回对应的错误信息。
//   - 如果错误码不存在于映射中，则返回一个通用的服务器错误信息。
//
// 参数:
//   - errCode: int
//     错误码，用于查找错误信息。
//
// 返回值:
//   - string: 对应错误码的错误信息字符串。如果错误码不存在，则返回通用的服务器错误信息。
func ErrMsg(errCode int) string {
	// 尝试从 codeText 中查找错误码对应的错误信息
	if msg, ok := codeText[errCode]; ok {
		return msg
	}
	// 如果错误码不存在，返回一个通用的服务器错误信息
	return codeText[ServerCommonError]
}
