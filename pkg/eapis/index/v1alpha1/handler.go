package v1alpha1

import (
	"github.com/emicklei/go-restful/v3"
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
	print("1111")
}
