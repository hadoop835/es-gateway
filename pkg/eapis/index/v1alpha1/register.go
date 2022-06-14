package v1alpha1

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/es-gateway/pkg/apiserver/runtime"
	schema "github.com/es-gateway/pkg/schema/meta/v1"
)

var GroupVersion = schema.GroupVersion{Group: "gateway.io", Version: "v1alpha1"}

func AddToContainer(container *restful.Container) error {
	ws := runtime.NewWebService(GroupVersion)
	handler := New()

	ws.Route(ws.GET("").To(handler.WithCreateIndex))

	container.Add(ws)
	return nil
}
