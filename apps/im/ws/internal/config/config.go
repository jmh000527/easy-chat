package config

import "github.com/zeromicro/go-zero/core/service"

// Config 结构体定义了服务的配置选项。
//
// 该结构体嵌套了服务配置和JWT认证相关的配置，以及 MongoDB 和消息传输的配置。
type Config struct {
	service.ServiceConf // 嵌套的服务通用配置

	ListenOn string // 服务监听的地址，例如 "0.0.0.0:8080"

	JwtAuth struct {
		AccessSecret string // JWT 认证的访问密钥，用于签名和验证 JWT 令牌
	}

	Mongo struct {
		Url string // MongoDB 连接的 URL，例如 "mongodb://localhost:27017"
		Db  string // 使用的 MongoDB 数据库名称
	}

	MsgChatTransfer struct {
		Topic string   // 消息聊天传输的主题名称
		Addrs []string // 消息传输服务的地址列表，例如 Kafka 或其他消息队列地址
	}

	MsgReadTransfer struct {
		Topic string   // 消息已读传输的主题名称
		Addrs []string // 消息传输服务的地址列表
	}
}
