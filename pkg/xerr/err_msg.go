package xerr

var codeText = map[int]string{
	ServerCommonError: "服务器异常，稍后再尝试",
	RequestParamError: "请求参数有误",
	DbError:           "数据库繁忙，稍后再尝试",
}

// ErrMsg 根据错误码获取错误信息
func ErrMsg(errCode int) string {
	if msg, ok := codeText[errCode]; ok {
		return msg
	}
	return codeText[ServerCommonError]
}
