package main

import (
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/resultx"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"log"

	"easy-chat/apps/user/api/internal/config"
	"easy-chat/apps/user/api/internal/handler"
	"easy-chat/apps/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// 加载go-zero的配置文件
	//conf.MustLoad(*configFile, &c)
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",             // Etcd 端点
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目密钥
		Namespace:      "user",                             // Etcd 命名空间
		Configs:        "user-api.yaml",                    // 配置文件名
		ConfigFilePath: "./etc/conf",                       // 配置文件路径（先删除再加载）
		LogLevel:       "DEBUG",                            // 日志级别
	})).MustLoad(&c)
	if err != nil {
		log.Fatal(err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
