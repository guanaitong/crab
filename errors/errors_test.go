package errors

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestNewError(t *testing.T) {
	f := 11222221.231231111111
	fmt.Sprint("")
	fmt.Println(f)
	fmt.Println(strconv.FormatFloat(f, 'g', -1, 32))
	NewSystemDbError(errors.New("x"))
	NewSystemRedisError(errors.New("x"))
	NewBusinessError(123, "xxx")
	NewApiError(123, "xx", errors.New("x"))
}
