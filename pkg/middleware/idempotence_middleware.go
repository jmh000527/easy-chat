package middleware

import (
	"easy-chat/pkg/interceptor"
	"net/http"
)

// IdempotenceMiddleware 结构体用于实现幂等性中间件。
// 幂等性中间件旨在确保相同的请求多次执行结果相同，且对系统没有副作用。
type IdempotenceMiddleware struct{}

// NewIdempotenceMiddleware 创建并返回一个新的 IdempotenceMiddleware 实例。
// 该函数是中间件的构造函数，负责初始化中间件实例。
func NewIdempotenceMiddleware() *IdempotenceMiddleware {
	return &IdempotenceMiddleware{}
}

// Handler 方法用于装饰接下来的 HTTP 处理函数，以实现幂等性。
// 它接收一个 http.HandlerFunc 类型的参数 next，表示接下来的处理函数。
// 返回一个新的 http.HandlerFunc，该函数将在请求中注入一些上下文信息后，调用 next 进行处理。
func (m *IdempotenceMiddleware) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 通过 WithContext 方法为请求注入一些上下文信息。
		// 这里的 interceptor.ContextWithVal 函数可能是用于设置请求的某些特殊上下文值，以实现幂等性。
		r = r.WithContext(interceptor.ContextWithVal(r.Context()))
		// 调用接下来的处理函数，传递修改后的请求和响应写入器。
		next(w, r)
	}
}
