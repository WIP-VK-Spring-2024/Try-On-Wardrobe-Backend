package utils

import (
	"encoding/json"
)

const EmptyJson = "{}"

func GetJson[T any](bytes []byte) (*T, error) {
	result := new(T)

	err := json.Unmarshal(bytes, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
