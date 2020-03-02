package gconf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	//"strconv"
	"reflect"
	"strconv"
)

const (
	//fileType
	unkown     = iota //0，未知或者原始的string
	properties        //1，properties文件格式
	jsons             //2，json形式
)

type configData struct {
	raw            string
	fileType       int
	structuredData map[string]*Field
}

func (c *configData) unmarshal(v interface{}) error {
	elem := reflect.ValueOf(v).Elem()

L:
	for i := 0; i < elem.NumField(); i++ {
		fieldInfo := elem.Field(i)
		if !fieldInfo.CanSet() {
			continue
		}

		structField := elem.Type().Field(i) // a reflect.StructField

		//优先匹配tag:config
		var configTagName = structField.Tag.Get("config")
		if configTagName != "" {
			var lowName = strings.ToLower(configTagName)
			for k, fieldData := range c.structuredData {
				if lowName == strings.ToLower(k) { //命中
					fieldData.Set(fieldInfo)
					continue L
				}
			}
			continue
		}

		//后匹配字段名:兼容Java驼峰命名和properties命名特性小写下划线
		var lowName = strings.ToLower(structField.Name)
		for k, fieldData := range c.structuredData {
			if lowName == strings.ToLower(strings.Replace(k, "_", "", -1)) { //命中
				fieldData.Set(fieldInfo)
				continue L
			}
		}

	}

	return nil
}

type Field struct {
	Name     string
	RawValue interface{}
}

func (f *Field) AsString() string {
	if reflect.ValueOf(f.RawValue).Kind() == reflect.String {
		return (f.RawValue).(string)
	}
	return fmt.Sprint(f.RawValue)
}

func (f *Field) AsInt() (int, error) {
	if reflect.ValueOf(f.RawValue).Kind() == reflect.Int {
		return (f.RawValue).(int), nil
	}
	return strconv.Atoi(fmt.Sprint(f.RawValue))
}

func (f *Field) AsInt64() (int64, error) {
	if reflect.ValueOf(f.RawValue).Kind() == reflect.Int64 {
		return (f.RawValue).(int64), nil
	}
	return strconv.ParseInt(fmt.Sprint(f.RawValue), 10, 0)
}

func (f *Field) AsUint64() (uint64, error) {
	if reflect.ValueOf(f.RawValue).Kind() == reflect.Uint64 {
		return (f.RawValue).(uint64), nil
	}
	return strconv.ParseUint(fmt.Sprint(f.RawValue), 10, 0)
}

func (f *Field) AsFloat64() (float64, error) {
	if reflect.ValueOf(f.RawValue).Kind() == reflect.Float64 {
		return (f.RawValue).(float64), nil
	}
	return strconv.ParseFloat(fmt.Sprint(f.RawValue), 10)
}

func (f *Field) AsBool() (bool, error) {
	if reflect.ValueOf(f.RawValue).Kind() == reflect.Bool {
		return (f.RawValue).(bool), nil
	}
	return strconv.ParseBool(fmt.Sprint(f.RawValue))
}

func (f *Field) Set(value reflect.Value) {
	switch value.Kind() {
	case reflect.String:
		value.Set(reflect.ValueOf(f.AsString()))
	case reflect.Int:
		i, err := f.AsInt64()
		if err == nil {
			value.Set(reflect.ValueOf(int(i)))
		}
	case reflect.Int32:
		i, err := f.AsInt64()
		if err == nil {
			value.Set(reflect.ValueOf(int32(i)))
		}
	case reflect.Int64:
		i, err := f.AsInt64()
		if err == nil {
			value.Set(reflect.ValueOf(i))
		}
	case reflect.Uint:
		i, err := f.AsUint64()
		if err == nil {
			value.Set(reflect.ValueOf(uint(i)))
		}
	case reflect.Uint32:
		i, err := f.AsUint64()
		if err == nil {
			value.Set(reflect.ValueOf(uint32(i)))
		}
	case reflect.Uint64:
		i, err := f.AsUint64()
		if err == nil {
			value.Set(reflect.ValueOf(i))
		}
	case reflect.Float64:
		i, err := f.AsFloat64()
		if err == nil {
			value.Set(reflect.ValueOf(i))
		}
	case reflect.Float32:
		i, err := f.AsFloat64()
		if err == nil {
			value.Set(reflect.ValueOf(float32(i)))
		}
	case reflect.Bool:
		i, err := f.AsBool()
		if err == nil {
			value.Set(reflect.ValueOf(i))
		}

	}

}

func toStructuredData(key, raw string) *configData {
	if raw == "" {
		return &configData{
			raw:            raw,
			fileType:       unkown,
			structuredData: map[string]*Field{},
		}
	}

	if strings.HasSuffix(key, "properties") {
		var (
			part   []byte
			prefix bool
			lines  []string
			err    error
		)
		reader := bufio.NewReader(strings.NewReader(raw))
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
		var j = make(map[string]*Field)

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") {
				continue
			}

			kv := strings.SplitN(line, "=", 2)
			if len(kv) != 2 {
				continue
			}
			j[kv[0]] = &Field{
				Name:     kv[0],
				RawValue: kv[1],
			}
		}

		return &configData{
			raw:            raw,
			fileType:       properties,
			structuredData: j,
		}

	}

	var rawJson = make(map[string]interface{})
	err := json.Unmarshal([]byte(raw), &rawJson)

	if err == nil {
		var sd = make(map[string]*Field)

		for k, v := range rawJson {
			sd[k] = &Field{
				Name:     k,
				RawValue: v,
			}
		}

		return &configData{
			raw:            raw,
			fileType:       jsons,
			structuredData: sd,
		}

	}

	return &configData{
		raw:            raw,
		fileType:       unkown,
		structuredData: map[string]*Field{},
	}
}
