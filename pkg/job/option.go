package job

import "time"

// RetryOptions 是一个函数类型，用于设置 retryOptions
type RetryOptions func(opts *retryOptions)

// retryOptions 结构体定义了重试的配置选项
type retryOptions struct {
	timeout     time.Duration   // 超时时间
	retryNums   int             // 重试次数
	isRetryFunc IsRetryFunc     // 判断是否需要重试的函数
	retryJetLag RetryJetLagFunc // 重试间隔时间的函数
}

// newOptions 返回一个带有默认配置的 retryOptions 实例，并应用可选的配置函数
func newOptions(opts ...RetryOptions) *retryOptions {
	// 设置默认配置
	opt := &retryOptions{
		timeout:     DefaultRetryTimeout,
		retryNums:   DefaultRetryNums,
		isRetryFunc: RetryAlways,
		retryJetLag: RetryJetLagAlways,
	}

	// 应用所有传入的可选配置函数
	for _, options := range opts {
		options(opt)
	}
	return opt
}

// WithRetryTimeout 返回一个设置重试超时时间的 RetryOptions 函数
func WithRetryTimeout(timeout time.Duration) RetryOptions {
	return func(opts *retryOptions) {
		if timeout > 0 {
			opts.timeout = timeout
		}
	}
}

// WithRetryNums 返回一个设置重试次数的 RetryOptions 函数
func WithRetryNums(nums int) RetryOptions {
	return func(opts *retryOptions) {
		// 确保重试次数至少为1次
		opts.retryNums = 1

		if nums > 1 {
			opts.retryNums = nums
		}
	}
}

// WithIsRetryFunc 返回一个设置判断是否需要重试函数的 RetryOptions 函数
func WithIsRetryFunc(retryFunc IsRetryFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryFunc != nil {
			opts.isRetryFunc = retryFunc
		}
	}
}

// WithRetryJetLagFunc 返回一个设置重试间隔时间函数的 RetryOptions 函数
func WithRetryJetLagFunc(retryJetLagFunc RetryJetLagFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryJetLagFunc != nil {
			opts.retryJetLag = retryJetLagFunc
		}
	}
}
