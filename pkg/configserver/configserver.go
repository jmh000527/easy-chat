package configserver

import (
	"errors"
	"github.com/zeromicro/go-zero/core/conf"
)

type OnChange func([]byte) error

var ErrNotSetConfig = errors.New("config not set")

type ConfigServer interface {
	Build() error
	SetOnChange(OnChange)
	FromJsonBytes() ([]byte, error)
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

func (s *configServer) MustLoad(v any, onChange OnChange) error {
	if s.configFile == "" && s.ConfigServer == nil {
		return ErrNotSetConfig
	}
	if s.ConfigServer == nil {
		// 使用go-zero的默认配置
		conf.MustLoad(s.configFile, v)
		return nil
	}
	// 如果使用热加载
	if onChange != nil {
		s.ConfigServer.SetOnChange(onChange)
	}
	// 构建配置中心
	if err := s.ConfigServer.Build(); err != nil {
		return err
	}
	// 使用Sail拉取的配置
	data, err := s.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}
	return conf.LoadFromJsonBytes(data, v)
}

func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}
