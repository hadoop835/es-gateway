package v1alpha1

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/es-gateway/pkg/log"
)

/**
创建索引
*/
type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

/**

 */
func (i *Handler) WithCreateIndex(request *restful.Request, response *restful.Response) {
	log.Log().Info("创建所以")
}
