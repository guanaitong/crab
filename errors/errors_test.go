package errors_test

import (
	"errors"
	"fmt"
	errors2 "github.com/guanaitong/crab/errors"
	"strconv"
	"testing"
)

func TestNewError(t *testing.T) {
	f := 11222221.231231111111
	fmt.Println(f)
	fmt.Println(strconv.FormatFloat(f, 'g', -1, 32))
	errors2.NewDbError(errors.New("x"))
}
