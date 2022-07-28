package es

import (
	"crypto/tls"
	"github.com/bytedance/sonic/decoder"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/es-gateway/pkg/client/es/options"
	"github.com/es-gateway/pkg/log"
	"github.com/es-gateway/pkg/models/version"
	"net/http"
)

type Elastic struct {
	Client *elasticsearch.Client
}

func NewClient(opt *options.ElasticOptions) (*Elastic, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: opt.Addresses,
		Username:  opt.Username,
		Password:  opt.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	return &Elastic{Client: client}, err
}

/**
  获取集群信息
*/
func (elastic Elastic) Info() version.ElasticInfo {
	var version version.ElasticInfo
	res, err := elastic.Client.Info()
	if err != nil {
		log.INFO("get info response error: %s", err)
	}
	defer res.Body.Close()
	decoder.NewStreamDecoder(res.Body).Decode(&version)
	return version
}

func (elastic Elastic) Response(response *esapi.Response) map[string]interface{} {
	var body map[string]interface{} = map[string]interface{}{}
	if response.Body != nil {
		decoder.NewStreamDecoder(response.Body).Decode(&body)
	}
	body["StatusCode"] = response.StatusCode
	body["Header"] = response.Header
	return body
}
