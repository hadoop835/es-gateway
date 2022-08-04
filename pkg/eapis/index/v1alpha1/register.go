package v1alpha1

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/es-gateway/pkg/apiserver/runtime"
	"github.com/es-gateway/pkg/client/es"
	schema "github.com/es-gateway/pkg/schema/meta/v1"
	"golang.org/x/net/context"
	"net/http"
)

var GroupVersion = schema.GroupVersion{Group: "gateway.io", Version: "v1alpha1"}

func AddToContainer(container *restful.Container, elastic *es.Elastic, ctx context.Context) error {
	ws := runtime.NewWebService(GroupVersion)
	handler := New(elastic, ctx)
	ws.Route(ws.GET("/").To(handler.WithESInfo).Returns(http.StatusOK, "", nil))

	ws.Route(ws.HEAD("/{path}").To(handler.WithExistsIndex))

	ws.Route(ws.POST("/{path}/_refresh").To(handler.WithRefreshIndex))

	ws.Route(ws.PUT("/{path}").To(handler.WithCreateIndex))

	ws.Route(ws.PUT("/{path}/_mapping").To(handler.WithCreateMapping))
	//保存数据
	ws.Route(ws.POST("/{path}/_doc/").To(handler.WithCreateDocument))
	//修改数据
	ws.Route(ws.PUT("/{path}/_doc/{id}").To(handler.WithCreateDocument))

	//删除数据
	ws.Route(ws.DELETE("/{path}/_doc/{id}").To(handler.WithDeleteDocument))
	container.Add(ws)
	return nil
}
