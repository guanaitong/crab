package gconf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

const (
	//fileType
	unkown     = iota //0，未知或者原始的string
	properties        //1，properties文件格式
	jsons             //2，json形式
)

type valueHandler struct {
	cp          interface{} // config struct pointer
	handlerFunc func(value string, cp interface{}) error
}

func (vh *valueHandler) refresh(value string) error {
	return vh.handlerFunc(value, vh.cp)
}

type Value struct {
	key          string
	value        string
	fileType     int
	valueHandler *valueHandler
}

func newValue(key, value string) *Value {
	fileType := unkown
	if strings.HasSuffix(key, ".properties") {
		fileType = properties
	} else if strings.HasSuffix(key, ".json") {
		fileType = jsons
	}
	return &Value{
		key:          key,
		value:        value,
		fileType:     fileType,
		valueHandler: nil,
	}
}

func (v *Value) Raw() string {
	if v == nil {
		return ""
	}
	return v.value
}

func (v *Value) FileType() int {
	return v.fileType
}

func (v *Value) refresh(newValue string) bool {
	if v.value == newValue {
		return false
	}
	v.value = newValue
	if v.valueHandler != nil {
		v.valueHandler.refresh(newValue)
	}
	return true
}
func (v *Value) Register(x interface{}) error {
	if v.valueHandler != nil {
		panic("value has registered")
	}
	if v.FileType() == properties {
		v.valueHandler = &valueHandler{cp: x, handlerFunc: propFunc}
	} else if v.FileType() == jsons {
		v.valueHandler = &valueHandler{cp: x, handlerFunc: jsonFunc}
	} else {
		panic("unsupported filed type")
	}
	return v.valueHandler.refresh(v.value)
}

func jsonFunc(value string, cp interface{}) error {
	return json.Unmarshal([]byte(value), cp)
}

func propFunc(value string, cp interface{}) error {
	data := ReadMapFromProp(value)
	if len(data) == 0 {
		return nil
	}
	elem := reflect.ValueOf(cp).Elem()

	for i := 0; i < elem.NumField(); i++ {
		fieldInfo := elem.Field(i)
		if !fieldInfo.CanSet() {
			continue
		}

		structField := elem.Type().Field(i) // a reflect.StructField

		var value = ""
		//优先匹配tag:config
		if name := structField.Tag.Get("config"); name != "" {
			var lowName = strings.ToLower(name)
			for k, fieldData := range data {
				if lowName == strings.ToLower(k) { //命中
					value = fieldData
					break
				}
			}
		} else {
			//后匹配字段名:兼容Java驼峰命名和properties命名特性小写下划线
			var lowName = strings.ToLower(structField.Name)
			for k, fieldData := range data {
				if lowName == strings.ToLower(strings.Replace(k, "_", "", -1)) { //命中
					value = fieldData
					break
				}
			}
		}

		if value == "" {
			continue
		}
		setFieldValue(fieldInfo, value)

	}

	return nil
}

func setFieldValue(filed reflect.Value, value string) {
	switch filed.Kind() {
	case reflect.String:
		filed.Set(reflect.ValueOf(value))
	case reflect.Bool:
		i, err := strconv.ParseBool(value)
		if err == nil {
			filed.Set(reflect.ValueOf(i))
		}
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(int(i)))
		}
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(int8(i)))
		}
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(int16(i)))
		}
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(int32(i)))
		}
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(i))
		}
	case reflect.Uint:
		i, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(uint(i)))
		}
	case reflect.Uint8:
		i, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(uint8(i)))
		}
	case reflect.Uint16:
		i, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(uint16(i)))
		}
	case reflect.Uint32:
		i, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(uint32(i)))
		}
	case reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			filed.Set(reflect.ValueOf(i))
		}
	case reflect.Float64:
		i, err := strconv.ParseFloat(value, 10)
		if err == nil {
			filed.Set(reflect.ValueOf(i))
		}
	case reflect.Float32:
		i, err := strconv.ParseFloat(value, 10)
		if err == nil {
			filed.Set(reflect.ValueOf(float32(i)))
		}

	}

}
func ReadMapFromProp(value string) map[string]string {
	var (
		part   []byte
		prefix bool
		lines  []string
		err    error
	)
	reader := bufio.NewReader(strings.NewReader(value))
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	var j = make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		j[kv[0]] = kv[1]
	}
	return j
}
