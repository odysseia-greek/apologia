package aletheia

func quizIndex(policyName string, numberOfShards, numberOfReplicas int) map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   numberOfShards,
				"number_of_replicas": numberOfReplicas,
				"lifecycle.name":     policyName,
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"theme": map[string]interface{}{
					"type": "keyword",
				},
				"segment": map[string]interface{}{
					"type": "keyword",
				},
				"set": map[string]interface{}{
					"type": "integer",
				},
			},
		},
	}
}
