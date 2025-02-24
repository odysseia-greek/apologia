package aletheia

func quizIndex(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   2,
				"number_of_replicas": 1,
				"lifecycle.name":     policyName,
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"quizType": map[string]interface{}{
					"type": "keyword",
				},
				"theme": map[string]interface{}{
					"type": "keyword",
				},
				"segment": map[string]interface{}{
					"type": "keyword",
				},
				"set": map[string]interface{}{
					"type": "integer",
				},
				// 'content' field is not defined here as it won't be queried
			},
		},
	}
}
