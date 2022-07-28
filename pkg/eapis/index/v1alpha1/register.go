package v1alpha1

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/es-gateway/pkg/apiserver/runtime"
	"github.com/es-gateway/pkg/client/es"
	schema "github.com/es-gateway/pkg/schema/meta/v1"
	"net/http"
)

var GroupVersion = schema.GroupVersion{Group: "gateway.io", Version: "v1alpha1"}

func AddToContainer(container *restful.Container, elastic *es.Elastic) error {
	ws := runtime.NewWebService(GroupVersion)
	handler := New(elastic)
	ws.Route(ws.GET("/").To(handler.WithESInfo).Returns(http.StatusOK, "", nil))

	ws.Route(ws.HEAD("/{path}").To(handler.WithExistsIndex))

	ws.Route(ws.POST("/{path}/_refresh").To(handler.WithRefreshIndex))

	ws.Route(ws.PUT("/{path}").To(handler.WithCreateIndex))

	ws.Route(ws.PUT("/{path}/_mapping").To(handler.WithCreateMapping))

	container.Add(ws)
	return nil
}
