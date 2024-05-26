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
	s := sail.New(&sail.MetaConfig{
		ETCDEndpoints:  cfg.ETCDEndpoints,  // Etcd 端点
		ProjectKey:     cfg.ProjectKey,     // 项目密钥
		Namespace:      cfg.Namespace,      // Etcd 命名空间
		Configs:        cfg.Configs,        // 配置文件名
		ConfigFilePath: cfg.ConfigFilePath, // 配置文件路径（先删除再加载）
		LogLevel:       cfg.LogLevel,       // 日志级别
	})
	return &Sail{
		Sail: s,
	}
}

func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}
	v, err := s.MergeVipers()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	data := v.AllSettings()
	return json.Marshal(data)
}

func (s *Sail) Error() error {
	return s.Err()
}
