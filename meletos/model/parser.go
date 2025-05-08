package model

import (
	"bytes"
	"encoding/json"
)

func StructToMap(input interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(input)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.NewDecoder(&buf).Decode(&result)
	return result, err
}

func MapToStruct(input map[string]interface{}, output interface{}) error {
	raw, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, output)
}
