package gconf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type es struct {
	ClusterName string
	hosts       string
	Port        int32
}

type impower struct {
	Path   string `config:"path"`
	Cmd    string `config:"cmd"`
	Period int    `config:"period"`
	HasX   bool   `config:"has"`
}

func TestGetConfigCollection(t *testing.T) {
	t.Log("[impower]-------------------------------")
	d1 := GetConfigCollection("impower")
	t.Log(d1, d1.AsMap())

	dm1 := d1.GetConfigAsStructuredMap("deny.properties")
	t.Log(dm1)

	imp := new(impower)
	if err := d1.GetConfigAsBean("deny.properties", imp); err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v\n", imp)
	}

	t.Log("[userdoor]------------------------------")
	c := GetConfigCollection("userdoor")
	v1 := c.GetConfig("es.properties")
	t.Log(v1, v1, c.AsMap())

	c1 := GetConfigCollection("userdoor")
	t.Log(c1.AsMap())

	p := new(es)
	v2 := c1.GetConfigAsStructuredMap("es.properties")
	c1.GetConfigAsBean("es.properties", p)
	t.Log(v2, p)

}

type testBean struct {
	a string `config:"a"`
	b int
	c uint64
}

func TestReflectBase(t *testing.T) {

	p := new(testBean)
	v := reflect.ValueOf(p).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		fieldInfo1 := v.Field(i)
		fmt.Println(i)
		fmt.Println(fieldInfo.Name)
		fmt.Println(fieldInfo.Type)
		fmt.Println(fieldInfo.Tag)
		fmt.Println(fieldInfo1.Type().Name())

	}
}

func TestSplit(tt *testing.T) {
	t := "mysql_url=amy:1qazxsw@@tcp(mdb.servers.dev.ofc:3306)/notifyagent?charset=utf8"
	ss := strings.SplitN(t, "=", 2)
	for _, s := range ss {
		fmt.Println(s)
	}
}

func TestSplit1(tt *testing.T) {
	t := "mysql_url="
	ss := strings.SplitN(t, "=", 2)
	for _, s := range ss {
		fmt.Println(s)
	}
}
