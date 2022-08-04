package apiserver

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/es-gateway/pkg/apiserver/config"
	"github.com/es-gateway/pkg/client/es"
	indexv1alpha1 "github.com/es-gateway/pkg/eapis/index/v1alpha1"
	"github.com/es-gateway/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

type ApiServer struct {
	container      *restful.Container
	Server         *http.Server
	Config         *config.Config
	Elastic        *es.Elastic
	adminContainer *restful.Container
	AdminServer    *http.Server
}

func (s *ApiServer) PrepareRun(stopCh <-chan struct{}) error {
	container := restful.NewContainer()
	s.container = container
	s.container.Router(restful.CurlyRouter{})
	//拦截器
	s.installEsGatewayAPIs(stopCh, context.Background())
	s.Server.Handler = s.container
	//
	s.adminContainer = restful.NewContainer()
	s.adminContainer.Router(restful.CurlyRouter{})
	s.installAdminAPIs(stopCh)
	s.AdminServer.Handler = s.adminContainer
	return nil
}

func (s *ApiServer) installEsGatewayAPIs(stopCh <-chan struct{}, ctx context.Context) {
	s.container.Filter(routeLogging)
	indexv1alpha1.AddToContainer(s.container, s.Elastic, ctx)

}

func (s *ApiServer) installAdminAPIs(stopCh <-chan struct{}) {
	s.adminContainer.Filter(routeLogging)

}

// Route Filter (defines FilterFunction)
func routeLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Log().Info(fmt.Sprintf("%s,%s", req.Request.Method, req.Request.RequestURI))
	chain.ProcessFilter(req, resp)
}

func (s *ApiServer) Run(ctx context.Context) (err error) {
	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-ctx.Done()
		_ = s.Server.Shutdown(shutdownCtx)
		_ = s.AdminServer.Shutdown(shutdownCtx)
	}()
	go func() {
		log.Log().Info("start listening on ", zap.String("addr", s.AdminServer.Addr))
		if s.AdminServer.TLSConfig != nil {
			err = s.AdminServer.ListenAndServeTLS("", "")
		} else {
			err = s.AdminServer.ListenAndServe()
		}
	}()
	log.Log().Info("start listening on ", zap.String("addr", s.Server.Addr))
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}

	return err
}
