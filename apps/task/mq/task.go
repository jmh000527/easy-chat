package main

import (
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/apps/task/mq/internal/handler"
	"easy-chat/apps/task/mq/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	// 添加服务组
	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	// 启动注册的所有服务
	fmt.Println("Starting MQueue at ...")
	serviceGroup.Start()
}
