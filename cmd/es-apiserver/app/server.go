package app

import (
	"context"
	"github.com/es-gateway/cmd/es-apiserver/app/options"
	"github.com/es-gateway/pkg/apiserver/config"
	"github.com/es-gateway/pkg/signals"
	"github.com/spf13/cobra"
	"net/http"
)

/**
create api server
*/
func WithApiServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	_config, err := config.TryLoadFromDisk()
	if err == nil {
		s = &options.ServerRunOptions{
			ApiServerConfig: s.ApiServerConfig,
			Config:          _config,
			ProxyRunConfig:  s.ProxyRunConfig,
			ElasticConfig:   s.ElasticConfig,
		}
	} else {

	}
	cmd := &cobra.Command{
		Use: "es-gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(s, config.WatchConfigChange(), signals.SetupSignalHandler())
		},
		SilenceUsage: true,
	}
	return cmd
}

func Run(s *options.ServerRunOptions, configCh <-chan config.Config, ctx context.Context) error {
	ictx, cancelFunc := context.WithCancel(context.TODO())
	errCh := make(chan error)
	defer close(errCh)
	go func() {
		if err := run(s, ictx); err != nil {
			errCh <- err
		}
	}()
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			return nil
		case cfg := <-configCh:
			cancelFunc()
			s.Config = &cfg
			ictx, cancelFunc = context.WithCancel(context.TODO())
			go func() {
				if err := run(s, ictx); err != nil {
					errCh <- err
				}
			}()
		case err := <-errCh:
			cancelFunc()
			return err
		}
	}
}

func run(opt *options.ServerRunOptions, ctx context.Context) error {
	// start proxy
	//proxyServer, _ := proxy.NewProxy(opt.ProxyRunConfig)
	//proxyServer.Run()

	apiserver, err := opt.NewAPIServer(ctx.Done())
	if err != nil {
		return err
	}
	err = apiserver.PrepareRun(ctx.Done())
	if err != nil {
		return err
	}
	err = apiserver.Run(ctx)

	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
