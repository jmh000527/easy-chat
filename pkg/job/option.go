package job

import "time"

// RetryOptions 是一个函数类型，用于设置 retryOptions
type RetryOptions func(opts *retryOptions)

// retryOptions 结构体定义了重试的配置选项
type retryOptions struct {
	timeout     time.Duration   // 超时时间限制。
	retryNums   int             // 允许的重试次数。
	isRetryFunc IsRetryFunc     // 判断是否应该重试的函数。
	retryJetLag RetryJetLagFunc // 计算重试间隔的函数。
}

// newOptions 返回一个带有默认配置的 retryOptions 实例，并应用可选的配置函数
func newOptions(opts ...RetryOptions) *retryOptions {
	// 初始化默认配置
	opt := &retryOptions{
		timeout:     DefaultRetryTimeout,
		retryNums:   DefaultRetryNums,
		isRetryFunc: RetryAlways,
		retryJetLag: RetryJetLagAlways,
	}

	// 应用传入的配置函数
	for _, options := range opts {
		options(opt)
	}
	return opt
}

// WithRetryTimeout 返回一个 RetryOptions，用于设置重试操作的超时时间。
func WithRetryTimeout(timeout time.Duration) RetryOptions {
	return func(opts *retryOptions) {
		// 如果指定的超时时间大于0，则更新超时时间。
		if timeout > 0 {
			opts.timeout = timeout
		}
	}
}

// WithRetryNums 返回一个 RetryOptions，用于设置重试的次数。
func WithRetryNums(nums int) RetryOptions {
	return func(opts *retryOptions) {
		// 确保重试次数至少为1次
		opts.retryNums = 1

		// 如果指定的重试次数大于1，则更新重试次数。
		if nums > 1 {
			opts.retryNums = nums
		}
	}
}

// WithIsRetryFunc 返回一个 RetryOptions，用于设置判断是否应该重试的函数。
func WithIsRetryFunc(retryFunc IsRetryFunc) RetryOptions {
	return func(opts *retryOptions) {
		// 如果提供了重试判断函数，则替换现有的函数。
		if retryFunc != nil {
			opts.isRetryFunc = retryFunc
		}
	}
}

// WithRetryJetLagFunc 返回一个 RetryOptions，用于设置计算重试间隔的函数。
func WithRetryJetLagFunc(retryJetLagFunc RetryJetLagFunc) RetryOptions {
	return func(opts *retryOptions) {
		// 如果提供了重试间隔计算函数，则替换现有的函数。
		if retryJetLagFunc != nil {
			opts.retryJetLag = retryJetLagFunc
		}
	}
}
