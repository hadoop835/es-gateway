package v1alpha1

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/es-gateway/pkg/client/es"
	"github.com/es-gateway/pkg/log"
	"golang.org/x/net/context"
	"io/ioutil"
	"time"
)

/**
创建索引
*/
type Handler struct {
	elastic *es.Elastic
}

func New(elastic *es.Elastic) *Handler {
	return &Handler{elastic: elastic}
}

/**

 */
func (h *Handler) WithESInfo(request *restful.Request, response *restful.Response) {
	version := h.elastic.Info()
	response.AddHeader("X-Elastic-Product", "Elasticsearch")
	response.WriteEntity(version)
}

//判断索引是否存在
func (h *Handler) WithExistsIndex(request *restful.Request, response *restful.Response) {
	uri := request.PathParameters()["path"]
	req := esapi.IndicesExistsRequest{
		Index: []string{uri},
	}
	res, err := req.Do(context.Background(), h.elastic.Client)
	if err != nil {
		log.INFO("查询索引是否存在：%s，失败：%s", uri[1:len(uri)], err)
	}
	result := h.elastic.Response(res)
	log.InfoJSONTOString("查询索引是否存在：%s", result)
	response.WriteHeaderAndEntity(res.StatusCode, result)
}

//刷新索引
func (h Handler) WithRefreshIndex(request *restful.Request, response *restful.Response) {
	uri := request.Request.RequestURI
	req := esapi.IndicesRefreshRequest{
		Index: []string{uri[1:len(uri)]},
	}
	res, err := req.Do(context.Background(), h.elastic.Client)
	if err != nil {
		log.INFO("刷新索引值：%s，失败：%s", uri[1:len(uri)], err)
	}
	defer res.Body.Close()
	result := h.elastic.Response(res)
	log.InfoJSONTOString("刷新索引是否存在：%s", result)
	response.WriteHeaderAndEntity(res.StatusCode, result)
}

//创建索引
func (h Handler) WithCreateIndex(request *restful.Request, response *restful.Response) {
	index := request.PathParameters()["path"]
	req_master_timeout := request.QueryParameter("master_timeout")
	master_timeout, _ := time.ParseDuration(req_master_timeout)
	req_timeout := request.QueryParameter("timeout")
	body, _ := ioutil.ReadAll(request.Request.Body)
	log.INFO("创建索引参数：%s", body)
	timeout, _ := time.ParseDuration(req_timeout)
	req := esapi.IndicesCreateRequest{
		Index:         index,
		MasterTimeout: master_timeout,
		Pretty:        true,
		Timeout:       timeout,
		Body:          request.Request.Body,
	}
	res, err := req.Do(context.Background(), h.elastic.Client)
	if err != nil {
		log.INFO("创建索引值：%s，失败：%s", index, err)
	}
	defer res.Body.Close()
	result := h.elastic.Response(res)
	log.InfoJSONTOString("创建索引成功：%s", result)
	response.WriteHeaderAndEntity(res.StatusCode, result)
}

//创建mapping
func (h Handler) WithCreateMapping(request *restful.Request, response *restful.Response) {
	index := request.PathParameters()["path"]
	req_master_timeout := request.QueryParameter("master_timeout")
	master_timeout, _ := time.ParseDuration(req_master_timeout)
	req_timeout := request.QueryParameter("timeout")
	timeout, _ := time.ParseDuration(req_timeout)
	req := esapi.IndicesPutMappingRequest{
		Index:         []string{index},
		MasterTimeout: master_timeout,
		Timeout:       timeout,
		Pretty:        true,
		Body:          request.Request.Body,
	}
	res, err := req.Do(context.Background(), h.elastic.Client)
	if err != nil {
		log.INFO("创建mapping：%s，失败：%s", index, err)
	}
	defer res.Body.Close()
	body := h.elastic.Response(res)
	log.InfoJSONTOString("创建mapping成功：%s", body)
	response.WriteHeaderAndEntity(res.StatusCode, body)

}
