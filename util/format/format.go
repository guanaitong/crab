package format

import "encoding/json"

func AsString(v interface{}) string {
	r, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(r)
}

func AsDefaultString(v interface{}) string {
	r, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(r)
}

func AsBytes(v interface{}) []byte {
	r, err := json.Marshal(v)
	if err != nil {
		return []byte("")
	}
	return r
}

func AsDefaultBytes(v interface{}) []byte {
	r, err := json.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return r
}

func AsJson(d string, v interface{}) error {
	return json.Unmarshal([]byte(d), v)
}
