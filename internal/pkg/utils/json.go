package utils

import (
	"encoding/json"

	"github.com/mailru/easyjson"
)

func EasyJsonMarshal(value interface{}) ([]byte, error) {
	marshaler, ok := value.(easyjson.Marshaler)
	if ok {
		return easyjson.Marshal(marshaler)
	}
	return json.Marshal(value)
}

func EasyJsonUnmarshal(data []byte, value interface{}) error {
	unmarshaler, ok := value.(easyjson.Unmarshaler)
	if ok {
		return easyjson.Unmarshal(data, unmarshaler)
	}
	return json.Unmarshal(data, value)
}
