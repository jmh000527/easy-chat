package interceptor

import (
	"context"
	"easy-chat/pkg/xerr"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Idempotent interface {
	Identify(ctx context.Context, method string) string                             // 获取请求标识
	IsIdempotentMethod(method string) bool                                          // 是否支持幂等性
	TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquired bool)  // 幂等性的验证
	SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error // 执行后结果的保存
}

var (
	TKey = "easy-chat-idempotence-task-id"      // 请求任务标识
	DKey = "easy-chat-idempotence-dispatch-key" // 设置rpc调度中rpc请求的标识
)

// ContextWithVal 函数将一个新的唯一UUID添加到传入的context中，并返回新的context。
// ctx是传入的context，TKey是context中存储UUID的键。
// 返回值是新的context，包含新生成的UUID。
func ContextWithVal(ctx context.Context) context.Context {
	// 创建一个新的上下文，将请求的id设置为键值对中的值
	// 设置请求的id
	return context.WithValue(ctx, TKey, utils.NewUuid())
}

// NewIdempotenceClient 创建一个gRPC客户端的拦截器，用于实现幂等性控制
func NewIdempotenceClient(idempotent Idempotent) grpc.UnaryClientInterceptor {
	// 返回一个实现了grpc.UnaryClientInterceptor接口的函数
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 调用Idempotent接口的Identify方法，传入当前上下文和请求方法，获取唯一的幂等性key
		identify := idempotent.Identify(ctx, method)

		// 使用metadata.NewOutgoingContext创建一个新的上下文对象，将幂等性key添加到请求头部信息中
		// DKey是预设的键名，用于在gRPC请求头部中存储幂等性key
		ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
			DKey: {identify},
		})

		// 调用grpc.UnaryInvoker接口的invoker方法，发起实际的gRPC请求
		// 将新的上下文、请求方法、请求体、响应体、客户端连接和调用选项传递给invoker
		// 返回请求执行的结果，该结果表示请求是否成功执行
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewIdempotenceServer 返回一个grpc服务端拦截器，用于实现gRPC请求的幂等性处理
func NewIdempotenceServer(idempotent Idempotent) grpc.UnaryServerInterceptor {
	// 返回一个grpc.UnaryServerInterceptor类型的函数，作为gRPC服务端拦截器
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 从请求的上下文中获取请求的id
		identify := metadata.ValueFromIncomingContext(ctx, DKey)
		// 如果请求的id为空，或者请求的方法不需要进行幂等性处理
		if len(identify) == 0 || !idempotent.IsIdempotentMethod(info.FullMethod) {
			// 直接调用handler处理请求，并返回结果
			// 不进行幂等性处理
			return handler(ctx, req)
		}

		// 打印日志，表示请求进入幂等性处理，并显示请求的id
		fmt.Println("----", "请求进入 幂等性处理 ", identify)

		// 尝试获取锁，并判断是否成功获取到锁
		r, isAcquire := idempotent.TryAcquire(ctx, identify[0])
		if isAcquire {
			// 如果成功获取到锁，则调用handler处理请求，并保存响应结果和错误信息
			resp, err = handler(ctx, req)
			fmt.Println("---- 执行任务", identify)

			// 保存响应结果到幂等性处理模块中
			if err := idempotent.SaveResp(ctx, identify[0], resp, err); err != nil {
				// 如果保存失败，则返回错误信息
				return resp, err
			}

			// 返回处理结果和错误信息
			return resp, err
		}

		// 如果未能获取到锁，表示任务已经在执行中
		fmt.Println("----- 任务在执行", identify)

		// 如果任务已经执行完成，并且保存了响应结果
		if r != nil {
			// 打印日志，表示任务已经执行完成，并显示请求的id
			fmt.Println("--- 任务已经执行完了 ", identify)
			// 直接返回之前保存的响应结果，而不重新执行任务
			return r, nil
		}

		// 如果任务可能还在执行中，则返回错误信息，表示存在其他任务在执行相同的id
		// 可能还在执行
		return nil, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), fmt.Sprintf("存在其他任务在执行 id %v", identify[0])))
	}
}

var (
	DefaultIdempotent       = new(defaultIdempotent)                  // 默认幂等性的处理
	DefaultIdempotentClient = NewIdempotenceClient(DefaultIdempotent) // 默认幂等性的拦截器客户端
)

type defaultIdempotent struct {
	// 获取和设置请求的id，用于唯一标识每个请求
	*redis.Redis

	// 注意存储，用于缓存处理过的请求信息，以实现幂等性
	*collection.Cache

	// 设置方法对幂等的支持，用于指定哪些方法需要进行幂等性处理
	// key为方法名，value为是否需要幂等处理的布尔值
	method map[string]bool
}

// NewDefaultIdempotent 函数用于创建一个默认的Idempotent实例
// 它接受一个redis.RedisConf类型的参数c，表示Redis的配置信息
// 函数返回一个Idempotent接口类型的实例
func NewDefaultIdempotent(c redis.RedisConf) Idempotent {
	// 创建一个缓存实例，缓存过期时间为60*60秒（即1小时）
	cache, err := collection.NewCache(60 * 60)
	if err != nil {
		// 如果创建缓存实例出错，则抛出panic异常
		panic(err)
	}

	// 返回一个defaultIdempotent实例的指针，该实例包含以下字段：
	// Redis：使用传入的Redis配置信息创建一个Redis实例
	// Cache：上面创建的缓存实例
	// method：一个map，用于存储需要幂等处理的HTTP方法及其对应的true值
	// 目前仅"/social.social/GroupCreate"方法被设置为需要幂等处理
	return &defaultIdempotent{
		Redis: redis.MustNewRedis(c),
		Cache: cache,
		method: map[string]bool{
			"/social.social/GroupCreate": true,
		},
	}
}

// Identify 根据请求的上下文和方法生成唯一的RPC标识。
func (d *defaultIdempotent) Identify(ctx context.Context, method string) string {
	// 从上下文中获取标识，可能是用户ID或其他唯一标识符
	id := ctx.Value(TKey)
	if id == nil {
		// 如果上下文中没有标识，则可能需要处理错误或生成一个默认的标识
		// 这里简单起见，返回空字符串，实际使用时应该更严谨地处理
		return ""
	}

	// 将标识和方法拼接成RPC标识字符串
	rpcId := fmt.Sprintf("%v.%s", id, method)
	return rpcId
}

// IsIdempotentMethod 判断给定的方法是否支持幂等性。
func (d *defaultIdempotent) IsIdempotentMethod(fullMethod string) bool {
	// 假设d.method是一个map，存储了支持幂等性的方法列表
	// 返回该方法是否在map中存在
	if d.method == nil {
		return false
	}
	return d.method[fullMethod]
}

// TryAcquire 尝试获取幂等性锁，如果成功则执行任务，否则返回缓存结果。
func (d *defaultIdempotent) TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool) {
	// 尝试在Redis中设置键值对，表示正在执行的任务
	// SetnxEx是Redis的SETNX和EXPIRE命令的结合，用于设置键值对并设置过期时间
	retry, err := d.SetnxEx(id, "1", 60*60) // 假设过期时间为1小时
	if err != nil {
		// 如果设置失败，返回错误和false
		return nil, false
	}

	if retry {
		// 如果设置成功（即key不存在时设置成功），表示获取到锁，返回nil和true
		return nil, true
	}

	// 如果设置失败（即key已存在），则从缓存中获取结果
	resp, _ = d.Cache.Get(id)

	// 返回缓存结果和false，表示没有获取到锁
	return resp, false
}

// SaveResp 保存执行后的结果到缓存中。
func (d *defaultIdempotent) SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error {
	// 如果存在响应错误，可能需要将错误也保存到缓存中，以便后续可以识别任务失败的原因
	// 这里简单起见，只保存响应结果，不保存错误
	d.Cache.Set(id, resp)
	return nil
}
