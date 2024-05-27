package configserver

import (
	"encoding/json"
	"fmt"
	"github.com/HYY-yu/sail-client"
)

type Config struct {
	ETCDEndpoints  string `toml:"etcd_endpoints"` // 逗号分隔的ETCD地址，0.0.0.0:2379,0.0.0.0:12379,0.0.0.0:22379
	ProjectKey     string `toml:"project_key"`
	Namespace      string `toml:"namespace"`
	Configs        string `toml:"configs"`
	ConfigFilePath string `toml:"config_file_path"` // 本地配置文件存放路径，空代表不存储本都配置文件
	LogLevel       string `toml:"log_level"`        // 日志级别(DEBUG\INFO\WARN\ERROR)，默认 WARN
}

type Sail struct {
	*sail.Sail
	sail.OnConfigChange
	c *Config
}

func NewSail(cfg *Config) *Sail {
	return &Sail{
		c: cfg,
	}
}

func (s *Sail) Build() error {
	var opts []sail.Option
	if s.OnConfigChange != nil {
		opts = append(opts, sail.WithOnConfigChange(s.OnConfigChange))
	}
	s.Sail = sail.New(&sail.MetaConfig{
		ETCDEndpoints:  s.c.ETCDEndpoints,  // Etcd 端点
		ProjectKey:     s.c.ProjectKey,     // 项目密钥
		Namespace:      s.c.Namespace,      // Etcd 命名空间
		Configs:        s.c.Configs,        // 配置文件名
		ConfigFilePath: s.c.ConfigFilePath, // 配置文件路径（先删除再加载）
		LogLevel:       s.c.LogLevel,       // 日志级别
	}, opts...)
	return s.Sail.Err()
}

func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}
	return s.fromJsonBytes(s.Sail)
}

func (s *Sail) fromJsonBytes(sail *sail.Sail) ([]byte, error) {
	v, err := sail.MergeVipers()
	if err != nil {
		return nil, err
	}
	data := v.AllSettings()
	return json.Marshal(data)
}

func (s *Sail) SetOnChange(f OnChange) {
	// 设置热加载方法
	s.OnConfigChange = func(configFileKey string, sail *sail.Sail) {
		data, err := s.fromJsonBytes(sail)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = f(data); err != nil {
			fmt.Println("OnChange err: ", err)
		}
	}
}
