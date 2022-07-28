package options

type ElasticOptions struct {
	Addresses []string
	Username  string
	Password  string
}

func NewElasticConfig() *ElasticOptions {
	return &ElasticOptions{
		Addresses: []string{"http://172.22.80.237:9200"},
		Username:  "",
		Password:  "",
	}
}
