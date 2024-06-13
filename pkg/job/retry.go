package job

import (
	"context"
	"errors"
	"time"
)

// ErrJobTimeout 定义任务超时的错误
var ErrJobTimeout = errors.New("任务超时")

// RetryJetLagFunc 定义重试的时间策略函数类型
type RetryJetLagFunc func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration

// RetryJetLagAlways 返回默认的重试间隔时间
func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

// IsRetryFunc 定义是否进行重试的函数类型
type IsRetryFunc func(ctx context.Context, retryCount int, err error) bool

// RetryAlways 始终返回true，表示总是重试
func RetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

// WithRetry 执行一个带有重试逻辑的处理函数
func WithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {
	// 使用传入的选项来初始化retryOptions
	opt := newOptions(opts...)

	// 检查传入的上下文是否设置了超时时间
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		// 如果没有设置超时时间，则使用默认的超时时间
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		herr        error                 // 用于存储handler的错误
		retryJetLag time.Duration         // 重试间隔时间
		ch          = make(chan error, 1) // 用于接收handler的执行结果
	)

	// 执行重试逻辑
	for i := 0; i < opt.retryNums; i++ {
		go func() {
			ch <- handler(ctx)
		}()

		select {
		case herr = <-ch: // 处理函数执行完毕
			if herr == nil {
				return nil // 如果没有错误，直接返回
			}

			if !opt.isRetryFunc(ctx, i, herr) {
				return herr // 错误不为空，如果不需要重试，返回错误
			}

			// 计算下次重试的间隔时间并等待
			retryJetLag = opt.retryJetLag(ctx, i, retryJetLag)
			time.Sleep(retryJetLag)
		case <-ctx.Done(): // 上下文超时或被取消
			return ErrJobTimeout
		}
	}

	return herr // 返回最后一次的错误
}
