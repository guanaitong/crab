package format

import "encoding/json"

func AsString(v interface{}) string {
	r, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(r)
}

func AsJson(d string, v interface{}) error {
	return json.Unmarshal([]byte(d), v)
}
