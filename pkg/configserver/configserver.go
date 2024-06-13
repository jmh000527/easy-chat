package configserver

import (
	"errors"
	"github.com/zeromicro/go-zero/core/conf"
)

// OnChange 是配置变更时的回调函数类型
type OnChange func([]byte) error

// ErrNotSetConfig 表示配置未设置的错误
var ErrNotSetConfig = errors.New("config not set")

// ConfigServer 接口定义了配置服务器的行为
type ConfigServer interface {
	Build() error
	SetOnChange(OnChange)
	FromJsonBytes() ([]byte, error)
}

// configServer 是一个配置服务器的代理类
type configServer struct {
	ConfigServer
	configFile string
}

// NewConfigServer 创建一个新的 configServer 实例
func NewConfigServer(configFile string, s ConfigServer) *configServer {
	return &configServer{
		ConfigServer: s,
		configFile:   configFile,
	}
}

// MustLoad 加载配置，支持本地配置文件和热加载配置中心
func (s *configServer) MustLoad(v any, onChange OnChange) error {
	// 检查配置文件路径和配置服务器是否都未设置
	if s.configFile == "" && s.ConfigServer == nil {
		return ErrNotSetConfig
	}

	// 如果没有配置服务器，使用go-zero的默认配置加载方式
	if s.ConfigServer == nil {
		conf.MustLoad(s.configFile, v)
		return nil
	}

	// 如果设置了热加载回调函数，则进行配置热加载
	if onChange != nil {
		s.ConfigServer.SetOnChange(onChange)
	}

	// 构建配置中心
	if err := s.ConfigServer.Build(); err != nil {
		return err
	}

	// 使用配置中心拉取的配置
	data, err := s.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}

	// 从JSON字节加载配置
	return conf.LoadFromJsonBytes(data, v)
}

// LoadFromJsonBytes 从JSON字节加载配置
func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}
