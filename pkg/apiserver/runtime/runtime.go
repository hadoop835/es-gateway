package runtime

import (
	"github.com/emicklei/go-restful/v3"
	schema "github.com/es-gateway/pkg/schema/meta/v1"
)

const (
	//ApiRootPath = "/eapis"
	ApiRootPath = "/"
)

const MimeMergePatchJson = "application/merge-patch+json"
const MimeJsonPatchJson = "application/json-patch+json"

func init() {
	restful.RegisterEntityAccessor(MimeMergePatchJson, restful.NewEntityAccessorJSON(restful.MIME_JSON))
	restful.RegisterEntityAccessor(MimeJsonPatchJson, restful.NewEntityAccessorJSON(restful.MIME_JSON))
}

func NewWebService(gv schema.GroupVersion) *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path(ApiRootPath).
		Produces(restful.MIME_JSON)
	return &webservice
}
