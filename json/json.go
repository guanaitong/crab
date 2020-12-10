package json

import (
	"encoding/json"
	"github.com/guanaitong/crab/util/strings2"
)

func MarshalToString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalFromString(str string, v interface{}) error {
	return Unmarshal(strings2.StringToBytes(str), v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// 忽略异常
func AsString(v interface{}) string {
	r, err := Marshal(v)
	if err != nil {
		return ""
	}
	return string(r)
}

// 忽略异常
func AsBytes(v interface{}) []byte {
	r, err := Marshal(v)
	if err != nil {
		return []byte("")
	}
	return r
}

func AsJson(str string, v interface{}) error {
	return Unmarshal(strings2.StringToBytes(str), v)
}