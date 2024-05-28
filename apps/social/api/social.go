package main

import (
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/resultx"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"log"
	"sync"

	"easy-chat/apps/social/api/internal/config"
	"easy-chat/apps/social/api/internal/handler"
	"easy-chat/apps/social/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "social",
		Configs:        "social-api.yaml",
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
		proc.WrapUp()
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
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
