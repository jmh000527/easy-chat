package main

import (
	"easy-chat/apps/im/api/internal/config"
	"easy-chat/apps/im/api/internal/handler"
	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/pkg/configserver"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"log"
	"sync"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "/im/conf/im-api.yaml", "the config file")
var wg sync.WaitGroup
var mu sync.Mutex

var signalChan = make(chan struct{})

func main() {
	flag.Parse()

	var c config.Config

	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "im",
		Configs:        "im-api.yaml",
		ConfigFilePath: "/im/conf",
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
		proc.WrapUp()
		proc.Shutdown()

		// 等待任务完成
		<-signalChan

		// 另外启动一个服务
		go func(c config.Config) {
			defer func() {
				wg.Add(1)
				wg.Done()
				// 发送信号通知任务完成
				signalChan <- struct{}{}
			}()
			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// 首次启动服务
	wg.Add(1)
	go func(c config.Config) {
		defer func() {
			wg.Add(1)
			wg.Done()
			// 发送信号通知任务完成
			signalChan <- struct{}{}
		}()
		Run(c)
	}(c)

	wg.Wait()
}

func Run(c config.Config) {
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
