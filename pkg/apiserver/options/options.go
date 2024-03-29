package options

import (
	"fmt"
	"github.com/es-gateway/pkg/utils/net"
	"github.com/spf13/pflag"
	"os"
)

type ApiServerConfig struct {
	// server bind address
	BindAddress string

	// insecure port number
	InsecurePort int

	// secure port number
	SecurePort int

	// tls cert file
	TlsCertFile string

	// tls private key file
	TlsPrivateKey string
	//admin
	AdminServerConfig *ApiServerConfig
}

func NewApiServerConfig() *ApiServerConfig {
	s := ApiServerConfig{
		BindAddress:   "0.0.0.0",
		InsecurePort:  9200,
		SecurePort:    0,
		TlsCertFile:   "",
		TlsPrivateKey: "",
		AdminServerConfig: &ApiServerConfig{
			BindAddress:   "0.0.0.0",
			InsecurePort:  8000,
			SecurePort:    0,
			TlsCertFile:   "",
			TlsPrivateKey: "",
		},
	}

	return &s
}

func (s *ApiServerConfig) withValidate() []error {
	errs := []error{}

	if s.SecurePort == 0 && s.InsecurePort == 0 {
		errs = append(errs, fmt.Errorf("insecure and secure port can not be disabled at the same time"))
	}

	if net.IsValidPort(s.SecurePort) {
		if s.TlsCertFile == "" {
			errs = append(errs, fmt.Errorf("tls cert file is empty while secure serving"))
		} else {
			if _, err := os.Stat(s.TlsCertFile); err != nil {
				errs = append(errs, err)
			}
		}

		if s.TlsPrivateKey == "" {
			errs = append(errs, fmt.Errorf("tls private key file is empty while secure serving"))
		} else {
			if _, err := os.Stat(s.TlsPrivateKey); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func (s *ApiServerConfig) AddFlags(fs *pflag.FlagSet, c *ApiServerConfig) {
	fs.StringVar(&s.BindAddress, "bind-address", c.BindAddress, "server bind address")
	fs.IntVar(&s.InsecurePort, "insecure-port", c.InsecurePort, "insecure port number")
	fs.IntVar(&s.SecurePort, "secure-port", s.SecurePort, "secure port number")
	fs.StringVar(&s.TlsCertFile, "tls-cert-file", c.TlsCertFile, "tls cert file")
	fs.StringVar(&s.TlsPrivateKey, "tls-private-key", c.TlsPrivateKey, "tls private key")
}
