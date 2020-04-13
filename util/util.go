package util

import (
	"crypto/rsa"
	"math/big"
	"reflect"
	"unsafe"
)

func Int32Ptr(i int32) *int32 { return &i }

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func RsaPublicDecrypt(pubKey *rsa.PublicKey, data []byte) []byte {
	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(data)
	e := big.NewInt(int64(pubKey.E))
	c.Exp(m, e, pubKey.N)
	out := c.Bytes()
	skip := 0
	for i := 2; i < len(out); i++ {
		if i+1 >= len(out) {
			break
		}
		if out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}
	return out[skip:]
}

// 两个struct之间的拷贝。只支持field为基础类型的拷贝，包括所有int、float、bool、string
func CopyStruct(dst, src interface{}) {
	if dst == nil || src == nil {
		return
	}
	dv, sv := reflect.ValueOf(dst), reflect.ValueOf(src)

	if !isStruct(sv) || !isStruct(dv) || dv.Kind() != reflect.Ptr || sv.Kind() != reflect.Ptr {
		return
	}
	dv = reflect.Indirect(dv)
	sv = reflect.Indirect(sv)
	for i, len := 0, dv.NumField(); i < len; i++ {
		fieldName := dv.Type().Field(i).Name
		dfv := dv.FieldByName(fieldName)
		sfv := sv.FieldByName(fieldName)
		if !sfv.IsValid() || !dfv.CanSet() || dfv.Kind() != sfv.Kind() {
			continue
		}
		switch dfv.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Bool, reflect.String:
			dfv.Set(sfv)
		}
	}
}

func isStruct(v reflect.Value) bool {
	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}
	pv := reflect.Indirect(v)
	// struct is not yet initialized
	if pv.Kind() == reflect.Invalid {
		return false
	}
	return pv.Kind() == reflect.Struct
}
