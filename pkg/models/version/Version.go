package version

type ElasticInfo struct {
	Name        string  `json:"name,omitempty"`
	ClusterName string  `json:"cluster_name,omitempty"`
	ClusterUuid string  `json:"cluster_uuid,omitempty"`
	Tagline     string  `json:"tagline,omitempty"`
	Version     version `json:"version,omitempty"`
}

type version struct {
	Number                           string `json:"number,omitempty"`
	BuildFlavor                      string `json:"build_flavor,omitempty"`
	BuildType                        string `json:"build_type,omitempty"`
	BuildHash                        string `json:"build_hash,omitempty"`
	BuildDate                        string `json:"build_date,omitempty"`
	BuildSnapshot                    bool   `json:"build_snapshot"`
	LuceneVersion                    string `json:"lucene_version,omitempty"`
	MinimumWireCompatibilityVersion  string `json:"minimum_wire_compatibility_version,omitempty"`
	MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version,omitempty"`
}
