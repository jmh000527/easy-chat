package job

import "time"

// 定义重试机制的默认参数常量
const (
	// DefaultRetryJetLag 默认的重试延迟时间，表示两次重试之间的时间间隔
	DefaultRetryJetLag = time.Second
	// DefaultRetryTimeout 默认的重试超时时间，超过这个时间不再重试
	DefaultRetryTimeout = 2 * time.Second
	// DefaultRetryNums 默认的重试次数，表示最多尝试多少次
	DefaultRetryNums = 5
)
