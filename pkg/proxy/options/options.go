package options

type ProxyRunConfig struct {
	// server bind address
	BindAddress string
	// insecure port number
	InsecurePort int
	//the read timeout of connections.
	ReadTimeout int
	//idle timeout of connections.
	DleTimeout int
	// tls cert file
	TlsCertFile string

	// tls private key file
	TlsPrivateKey string
}

func NewProxyRunConfig() *ProxyRunConfig {
	return &ProxyRunConfig{
		BindAddress:   "0.0.0.0",
		InsecurePort:  3307,
		ReadTimeout:   3000,
		DleTimeout:    1000,
		TlsPrivateKey: "",
		TlsCertFile:   "",
	}
}
