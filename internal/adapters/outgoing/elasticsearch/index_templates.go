package elasticsearch

// IndexMapping represents the structure of an Elasticsearch index mapping
type IndexMapping struct {
	Settings IndexSettings          `json:"settings"`
	Mappings map[string]interface{} `json:"mappings"`
}

// IndexSettings represents the settings section of an Elasticsearch index
type IndexSettings struct {
	NumberOfShards   int                    `json:"number_of_shards"`
	NumberOfReplicas int                    `json:"number_of_replicas"`
	Analysis         map[string]interface{} `json:"analysis,omitempty"`
}
