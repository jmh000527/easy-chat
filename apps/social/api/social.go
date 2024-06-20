package main

import (
	"easy-chat/apps/social/api/internal/config"
	"easy-chat/apps/social/api/internal/handler"
	"easy-chat/apps/social/api/internal/svc"
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/resultx"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"log"
	"sync"
)

var configFile = flag.String("f", "/social/conf/social-api.yaml", "the config file")

var wg sync.WaitGroup
var mu sync.Mutex

var signalChan = make(chan struct{})

func main() {
	flag.Parse()

	var c config.Config
	// 加载go-zero的配置文件
	//conf.MustLoad(*configFile, &c)
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.199.138:3379",             // Etcd 端点
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目密钥
		Namespace:      "social",                           // Etcd 命名空间
		Configs:        "social-api.yaml",                  // 配置文件名
		ConfigFilePath: "/social/conf",                     // 配置文件路径（先删除再加载）
		//ConfigFilePath: "./etc/conf", // 配置文件路径（先删除再加载）
		LogLevel: "DEBUG", // 日志级别
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

// Run 启动基于配置的RESTful API服务。
// c: 包含服务配置信息的配置对象。
func Run(c config.Config) {
	// 初始化REST服务器，应用配置中的REST配置，并启用CORS。
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	// 确保服务器在函数返回前停止，以进行资源清理。
	defer server.Stop()

	// 创建服务上下文，传递配置信息。
	ctx := svc.NewServiceContext(c)
	// 注册服务处理程序，将服务器与业务逻辑连接起来。
	handler.RegisterHandlers(server, ctx)

	// 设置错误处理程序和成功处理程序，以统一处理HTTP请求的结果。
	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	// 输出启动信息，指示服务正在启动。
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// 启动服务器，开始监听HTTP请求。
	server.Start()
}
