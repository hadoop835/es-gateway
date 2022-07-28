package options

import (
	"crypto/tls"
	"fmt"
	"github.com/es-gateway/pkg/apiserver"
	"github.com/es-gateway/pkg/apiserver/config"
	apiServerConfig "github.com/es-gateway/pkg/apiserver/options"
	"github.com/es-gateway/pkg/client/es"
	esConfig "github.com/es-gateway/pkg/client/es/options"
	proxyRunConfig "github.com/es-gateway/pkg/proxy/options"
	"net/http"
)

type ServerRunOptions struct {
	ConfigFile      string
	ApiServerConfig *apiServerConfig.ApiServerConfig
	ProxyRunConfig  *proxyRunConfig.ProxyRunConfig
	*config.Config
	ElasticConfig *esConfig.ElasticOptions
}

//
func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		ApiServerConfig: apiServerConfig.NewApiServerConfig(),
		Config:          config.NewConfig(),
		ProxyRunConfig:  proxyRunConfig.NewProxyRunConfig(),
		ElasticConfig:   esConfig.NewElasticConfig(),
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
	//es
	_apiServer.Elastic, _ = es.NewClient(s.ElasticConfig)
	//admin

	//http
	admin_server := &http.Server{
		Addr: fmt.Sprintf(":%d", s.ApiServerConfig.AdminServerConfig.InsecurePort),
	}
	//https
	if s.ApiServerConfig.AdminServerConfig.SecurePort != 0 {
		certificate, err := tls.LoadX509KeyPair(s.ApiServerConfig.AdminServerConfig.TlsCertFile, s.ApiServerConfig.AdminServerConfig.TlsPrivateKey)
		if err != nil {
			return nil, err
		}
		admin_server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}
		admin_server.Addr = fmt.Sprintf(":%d", s.ApiServerConfig.AdminServerConfig.SecurePort)
	}
	_apiServer.AdminServer = admin_server
	return _apiServer, nil
}
