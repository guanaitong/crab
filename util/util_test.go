package util_test

import (
	"fmt"
	"github.com/guanaitong/crab/util"
	"testing"
	"time"
)

func TestCopy(t *testing.T) {
	source := &User1{
		Int8:        1,
		Int16:       -16,
		Int32:       32,
		Int64:       64,
		Int:         -123,
		i:           123,
		UInt8:       8,
		UInt16:      16,
		UInt32:      32,
		UInt64:      64,
		UInt:        456,
		Float32:     2.718281828459,
		Float64:     3.141592654,
		Bool:        true,
		Byte:        34,
		String:      "xxx",
		StringArray: []string{"1", "2"},
		IntArray:    []int{1, -1},
		StringMap: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
		IntMap: map[int]int{
			1: -1,
			2: -2,
		},
		T: time.Now(),
	}
	dst := new(User2)
	util.CopyStruct(dst, source)
	if dst.T != source.T {
		t.Fail()
	}
	fmt.Println(dst)
}

type User1 struct {
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
	Int    int
	UInt8  uint8
	UInt16 uint16
	UInt32 uint32
	UInt64 uint64
	UInt   uint
	i      int
	T      time.Time

	Float32 float32
	Float64 float64
	Bool    bool
	Byte    byte
	String  string

	StringArray []string
	IntArray    []int
	StringMap   map[string]string
	IntMap      map[int]int
}

type User2 struct {
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
	Int    int
	UInt8  uint8
	UInt16 uint16
	UInt32 uint32
	UInt64 uint64
	UInt   uint
	i      int
	T      time.Time

	Float32 float32
	Float64 float64
	Bool    bool
	Byte    byte
	String  string

	StringArray []string
	IntArray    []int
	StringMap   map[string]string
	IntMap      map[int]int
}
