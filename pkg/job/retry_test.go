package job

import (
	"context"
	"errors"
	"testing"
	"time"
)

// 测试WithRetry函数
func TestWithRetry(t *testing.T) {
	// 定义一个测试错误
	var (
		ErrTest = errors.New("测试异常")
		// 模拟的handler函数，返回ErrTest错误
		handler = func(ctx context.Context) error {
			t.Log("执行handler")
			return ErrTest
		}
	)

	// 定义测试用例的参数类型
	type args struct {
		ctx     context.Context
		handler func(context.Context) error
		opts    []RetryOptions
	}

	// 定义测试用例
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			// 测试用例1：没有设置重试选项，期望超时错误
			"1", args{
				ctx:     context.Background(),
				handler: handler,
				opts:    []RetryOptions{},
			}, ErrJobTimeout,
		},
		{
			// 测试用例2：设置重试超时时间和重试间隔，期望最后返回测试错误
			"2", args{
				ctx:     context.Background(),
				handler: handler,
				opts: []RetryOptions{
					WithRetryTimeout(6 * time.Second), // 设置重试超时时间为3秒
					WithRetryJetLagFunc(func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
						return 1 * time.Second // 设置重试间隔为1秒
					}),
				},
			}, ErrTest,
		},
		{
			// 测试用例3：设置重试条件函数为false，期望立即返回测试错误
			"3", args{
				ctx:     context.Background(),
				handler: handler,
				opts: []RetryOptions{
					WithIsRetryFunc(func(ctx context.Context, retryCount int, err error) bool {
						return false // 不进行重试
					}),
				},
			}, ErrTest,
		},
	}

	// 遍历测试用例并运行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WithRetry(tt.args.ctx, tt.args.handler, tt.args.opts...); !errors.Is(err, tt.wantErr) {
				t.Errorf("WithRetry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
