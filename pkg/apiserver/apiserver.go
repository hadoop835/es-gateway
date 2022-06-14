package apiserver

import (
	"context"
	"github.com/emicklei/go-restful/v3"
	"github.com/es-gateway/pkg/apiserver/config"
	indexv1alpha1 "github.com/es-gateway/pkg/eapis/index/v1alpha1"
	"net/http"
)

type ApiServer struct {
	container *restful.Container
	Server    *http.Server
	Config    *config.Config
}

func (s *ApiServer) PrepareRun(stopCh <-chan struct{}) error {
	container := restful.NewContainer()
	s.container = container
	s.container.Router(restful.CurlyRouter{})
	//拦截器
	s.installEsGatewayAPIs(stopCh)

	s.Server.Handler = s.container
	return nil
}

func (s *ApiServer) installEsGatewayAPIs(stopCh <-chan struct{}) {
	indexv1alpha1.AddToContainer(s.container)

}

func (s *ApiServer) Run(ctx context.Context) (err error) {
	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-ctx.Done()
		_ = s.Server.Shutdown(shutdownCtx)
	}()
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}
	return err
}
