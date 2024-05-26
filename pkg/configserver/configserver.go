package configserver

import (
	"errors"
	"github.com/zeromicro/go-zero/core/conf"
)

var ErrNotSetConfig = errors.New("config not set")

type ConfigServer interface {
	FromJsonBytes() ([]byte, error)
	Error() error
}

// 代理类
type configServer struct {
	ConfigServer
	configFile string
}

func NewConfigServer(configFile string, s ConfigServer) *configServer {
	return &configServer{
		ConfigServer: s,
		configFile:   configFile,
	}
}

func (s *configServer) MustLoad(v any) error {
	if s.ConfigServer.Error() != nil {
		return s.ConfigServer.Error()
	}
	if s.configFile == "" && s.ConfigServer == nil {
		return ErrNotSetConfig
	}
	if s.ConfigServer == nil {
		// 使用go-zero的默认配置
		conf.MustLoad(s.configFile, v)
		return nil
	}
	// 使用Sail拉取的配置
	data, err := s.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}
	return conf.LoadFromJsonBytes(data, v)
}

func (s *configServer) Error() error {
	return s.ConfigServer.Error()
}
