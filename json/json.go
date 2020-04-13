package json

import (
	"github.com/guanaitong/crab/util"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func MarshalToString(v interface{}) (string, error) {
	return json.MarshalToString(v)
}
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalFromString(str string, v interface{}) error {
	return json.UnmarshalFromString(str, v)
}
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func AsString(v interface{}) string {
	r, err := Marshal(v)
	if err != nil {
		return ""
	}
	return string(r)
}

func AsBytes(v interface{}) []byte {
	r, err := Marshal(v)
	if err != nil {
		return []byte("")
	}
	return r
}

func AsJson(d string, v interface{}) error {
	return Unmarshal(util.StringToBytes(d), v)
}
