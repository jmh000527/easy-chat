package main

//
import (
	"easy-chat/apps/user/rpc/internal/server"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/interceptor/rpcserver"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/reflection"
	"log"
	"sync"

	"easy-chat/apps/user/rpc/internal/config"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

var wg sync.WaitGroup

var grpcServer *grpc.Server

func main() {
	flag.Parse()

	var c config.Config
	//conf.MustLoad(*configFile, &c)
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "user",
		Configs:        "user-rpc.yaml",
		ConfigFilePath: "/user/conf",
		//ConfigFilePath: "./etc/conf",
		LogLevel: "DEBUG",
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

	if err := ctx.SetRootToken(); err != nil {
		panic(err)
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {
		grpcServer = srv

		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

//var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")
//
//var grpcServer *grpc.Server
//var wg sync.WaitGroup
//
//func main() {
//	flag.Parse()
//
//	var c config.Config
//	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
//		ETCDEndpoints:  "192.168.199.138:3379",
//		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
//		Namespace:      "user",
//		Configs:        "user-rpc.yaml",
//		ConfigFilePath: "./etc/conf",
//		LogLevel:       "DEBUG",
//	})).MustLoad(&c, func(bytes []byte) error {
//		var c config.Config
//		configserver.LoadFromJsonBytes(bytes, &c)
//
//		grpcServer.GracefulStop()
//
//		wg.Add(1)
//		go func(c config.Config) {
//			defer wg.Done()
//
//			Run(c)
//		}(c)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	wg.Add(1)
//	go func(c config.Config) {
//		defer wg.Done()
//
//		Run(c)
//	}(c)
//
//	wg.Wait()
//}
//
//func Run(c config.Config) {
//	ctx := svc.NewServiceContext(c)
//
//	if err := ctx.SetRootToken(); err != nil {
//		panic(err)
//	}
//
//	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {
//		grpcServer = srv
//
//		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
//		if c.Mode == service.DevMode || c.Mode == service.TestMode {
//			reflection.Register(grpcServer)
//		}
//	})
//	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
//	defer s.Stop()
//
//	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
//	s.Start()
//}
