package main

import (
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/interceptor"
	"easy-chat/pkg/interceptor/rpcserver"
	"flag"
	"fmt"
	"log"
	"sync"

	"easy-chat/apps/social/rpc/internal/config"
	"easy-chat/apps/social/rpc/internal/server"
	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "social/conf/social-rpc.yaml", "the config file")
var grpcServer *grpc.Server
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "social",
		Configs:        "social-rpc.yaml",
		ConfigFilePath: "/social/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		err := configserver.LoadFromJsonBytes(bytes, &c)
		if err != nil {
			log.Println("load config err:", err)
			return err
		}
		log.Println("load config success, config info:", c)
		// 停止接受请求
		grpcServer.GracefulStop()
		// 另外启动一个服务
		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()

			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	wg.Wait()

}

func Run(c config.Config) {
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {

		grpcServer = srv

		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 为gRPC服务添加一元拦截器，以增强服务的功能和性能。
	// 这里分别添加了日志拦截器、幂等性拦截器和同步限流拦截器。
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	s.AddUnaryInterceptors(interceptor.NewIdempotenceServer(interceptor.NewDefaultIdempotent(c.Cache[0].RedisConf)))
	s.AddUnaryInterceptors(rpcserver.SyncXLimitInterceptor(100))

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
