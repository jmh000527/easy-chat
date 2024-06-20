package main

// Config 结构体用于存储应用程序的配置
type Config struct {
	Name     string
	Host     string
	Port     string
	Mode     string
	Database string

	UserRpc struct {
		Etcd struct {
			Hosts []string
			Key   string
		}
	}
	Redisx struct {
		Host string
		Pass string
	}
	JwtAuth struct {
		AccessSecret string
	}
}

//func main() {
//	var cfg Config
//	// 初始化 Sail 客户端
//	s := sail.New(&sail.MetaConfig{
//		ETCDEndpoints:  "192.168.199.138:3379",             // Etcd 端点
//		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目密钥
//		Namespace:      "im",                               // Etcd 命名空间
//		Configs:        "im-api.yaml",                      // 配置文件名
//		ConfigFilePath: "./conf",                           // 配置文件路径（先删除再加载）
//		LogLevel:       "DEBUG",                            // 日志级别
//	}, sail.WithOnConfigChange(func(configFileKey string, s *sail.Sail) {
//		if err := s.Err(); err != nil {
//			fmt.Println(err)
//			return
//		}
//		// 拉取初始配置
//		fmt.Println(s.Pull())
//		// 合并配置
//		v, err := s.MergeVipers()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		if err := v.Unmarshal(&cfg); err != nil {
//			fmt.Println(err)
//			return
//		}
//		fmt.Println(cfg, "\n", cfg.Database)
//	}))
//	if err := s.Err(); err != nil {
//		fmt.Println(err)
//		return
//	}
//	// 拉取初始配置
//	fmt.Println(s.Pull())
//	// 合并配置
//	v, err := s.MergeVipers()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if err := v.Unmarshal(&cfg); err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(cfg, "\n", cfg.Database)
//
//	// 无限循环以保持程序运行
//	for {
//		time.Sleep(time.Second)
//	}
//}
