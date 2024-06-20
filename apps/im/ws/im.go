package main

import (
	"easy-chat/apps/im/ws/internal/config"
	"easy-chat/apps/im/ws/internal/handler"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/pkg/configserver"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"log"
	"sync"
	"time"
)

var configFile = flag.String("f", "/im/conf/im-ws.yaml", "the config file")
var wg sync.WaitGroup

var signalChan = make(chan struct{})

//var configFile = flag.String("f", "C:/Users/jmh00/GolandProjects/easy-chat/apps/im/ws/etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "im",
		Configs:        "im-ws.yaml",
		ConfigFilePath: "/im/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		err := configserver.LoadFromJsonBytes(bytes, &c)
		if err != nil {
			log.Panicln("load config error:", err)
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
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithWebsocketAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerAck(websocket.NoAck),
		websocket.WithWebsocketMaxConnectionIdle(7*time.Hour),
		websocket.WithServerSendErrCount(3),
	)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()
}
