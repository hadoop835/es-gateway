package options

import (
	"crypto/tls"
	"fmt"
	"github.com/es-gateway/pkg/apiserver"
	"github.com/es-gateway/pkg/apiserver/config"
	apiServerConfig "github.com/es-gateway/pkg/apiserver/options"
	proxyRunConfig "github.com/es-gateway/pkg/proxy/options"
	"net/http"
)

type ServerRunOptions struct {
	ConfigFile      string
	ApiServerConfig *apiServerConfig.ApiServerConfig
	ProxyRunConfig  *proxyRunConfig.ProxyRunConfig
	*config.Config
}

//
func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		ApiServerConfig: apiServerConfig.NewApiServerConfig(),
		Config:          config.NewConfig(),
		ProxyRunConfig:  proxyRunConfig.NewProxyRunConfig(),
	}
	return s
}

/**
配置文件初始化
*/
func (s *ServerRunOptions) NewAPIServer(stopCh <-chan struct{}) (*apiserver.ApiServer, error) {
	_apiServer := &apiserver.ApiServer{
		Config: s.Config,
	}
	//http
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", s.ApiServerConfig.InsecurePort),
	}
	//https
	if s.ApiServerConfig.SecurePort != 0 {
		certificate, err := tls.LoadX509KeyPair(s.ApiServerConfig.TlsCertFile, s.ApiServerConfig.TlsPrivateKey)
		if err != nil {
			return nil, err
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}
		server.Addr = fmt.Sprintf(":%d", s.ApiServerConfig.SecurePort)
	}
	_apiServer.Server = server
	return _apiServer, nil
}
