package main

import (
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/apps/task/mq/internal/handler"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/pkg/configserver"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"log"
	"sync"
)

var configFile = flag.String("f", "/task/conf/task-mq.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "task",
		Configs:        "task-mq.yaml",
		ConfigFilePath: "/task/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		err := configserver.LoadFromJsonBytes(bytes, &c)
		if err != nil {
			log.Println("load config err:", err)
			return err
		}
		log.Println("load config success, config info:", c)

		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()

			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	wg.Wait()
}

func Run(c config.Config) {
	// 配置初始化，确保服务配置正确无误。
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	// 创建服务上下文。
	ctx := svc.NewServiceContext(c)

	// 初始化监听器，用于处理服务请求。
	listen := handler.NewListen(ctx)

	// 创建服务组，用于统一管理和启动服务。
	serviceGroup := service.NewServiceGroup()

	// 将所有服务添加到服务组中。
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}

	// 启动服务组，开始监听和处理请求。
	fmt.Println("start mqueue server at ", c.ListenOn, " ..... ")
	serviceGroup.Start()
}
