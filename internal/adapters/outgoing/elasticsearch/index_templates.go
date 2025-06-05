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

// DefaultDocumentMapping returns the default mapping structure for document indices
func DefaultDocumentMapping() IndexMapping {
	return IndexMapping{
		Settings: IndexSettings{
			NumberOfShards:   1,
			NumberOfReplicas: 1,
			Analysis: map[string]interface{}{
				"analyzer": map[string]interface{}{
					"html_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"char_filter": []string{
							"html_strip",
						},
						"filter": []string{
							"lowercase",
							"stop",
							"snowball",
							"unique",
						},
					},
					"url_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "uax_url_email",
						"filter": []string{
							"lowercase",
							"unique",
						},
					},
					"keyword_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
							"trim",
						},
					},
				},
				"filter": map[string]interface{}{
					"snowball": map[string]interface{}{
						"type":     "snowball",
						"language": "english",
					},
				},
			},
		},
		Mappings: map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"url": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "url_analyzer",
						},
					},
				},
				"title": map[string]interface{}{
					"type":     "text",
					"analyzer": "html_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":     "text",
							"analyzer": "keyword_analyzer",
						},
					},
				},
				"content": map[string]interface{}{
					"type":     "text",
					"analyzer": "html_analyzer",
				},
				"summary": map[string]interface{}{
					"type":     "text",
					"analyzer": "html_analyzer",
				},
				"language": map[string]interface{}{
					"type": "keyword",
				},
				"contentType": map[string]interface{}{
					"type": "keyword",
				},
				"metadata": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
				},
				"pageRank": map[string]interface{}{
					"type": "float",
				},
				"incomingLinks": map[string]interface{}{
					"type": "integer",
				},
				"outgoingLinks": map[string]interface{}{
					"type": "integer",
				},
				"lastCrawled": map[string]interface{}{
					"type": "date",
				},
				"lastModified": map[string]interface{}{
					"type": "date",
				},
				"createdAt": map[string]interface{}{
					"type": "date",
				},
				"keywords": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "html_analyzer",
						},
					},
				},
				"entities": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "keyword",
							"fields": map[string]interface{}{
								"text": map[string]interface{}{
									"type":     "text",
									"analyzer": "html_analyzer",
								},
							},
						},
						"type": map[string]interface{}{
							"type": "keyword",
						},
						"relevance": map[string]interface{}{
							"type": "float",
						},
					},
				},
				"EnhancedKeywords": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type": "keyword",
							"fields": map[string]interface{}{
								"analyzed": map[string]interface{}{
									"type":     "text",
									"analyzer": "keyword_analyzer",
								},
							},
						},
						"score": map[string]interface{}{
							"type": "float",
						},
						"isDomainSpecific": map[string]interface{}{
							"type": "boolean",
						},
						"category": map[string]interface{}{
							"type": "keyword",
						},
						"position": map[string]interface{}{
							"type": "integer",
						},
					},
				},
				"KeywordTexts": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "keyword_analyzer",
						},
					},
				},
				"ParsedContent": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
					"properties": map[string]interface{}{
						"keywords": map[string]interface{}{
							"type":    "object",
							"dynamic": true,
							"properties": map[string]interface{}{
								"texts": map[string]interface{}{
									"type": "keyword",
									"fields": map[string]interface{}{
										"analyzed": map[string]interface{}{
											"type":     "text",
											"analyzer": "keyword_analyzer",
										},
									},
								},
								"high_impact": map[string]interface{}{
									"type": "keyword",
									"fields": map[string]interface{}{
										"analyzed": map[string]interface{}{
											"type":     "text",
											"analyzer": "keyword_analyzer",
										},
									},
								},
								"by_domain": map[string]interface{}{
									"type":    "object",
									"dynamic": true,
								},
							},
						},
						"entities_simple": map[string]interface{}{
							"type":    "nested",
							"dynamic": true,
						},
					},
				},
			},
		},
	}
}
