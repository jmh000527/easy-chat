package websocket

import "net/http"

// DialOptions 是用于配置拨号选项的函数类型。
type DialOptions func(option *dialOption)

// dialOption 结构体保存了 WebSocket 连接的配置选项。
type dialOption struct {
	header  http.Header // HTTP 头部
	pattern string      // 连接路径模式
}

// newDialOptions 创建一个具有默认值的新的 dialOption 结构体，并根据传入的选项进行配置。
func newDialOptions(opts ...DialOptions) dialOption {
	// 默认值
	o := dialOption{
		header:  nil,
		pattern: "/ws",
	}
	// 应用传入的选项
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithClientPattern 返回一个设置连接路径模式的 DialOptions 函数。
func WithClientPattern(pattern string) DialOptions {
	return func(opt *dialOption) {
		opt.pattern = pattern
	}
}

// WithClientHeader 返回一个设置 HTTP 头部的 DialOptions 函数。
func WithClientHeader(header http.Header) DialOptions {
	return func(opt *dialOption) {
		opt.header = header
	}
}
