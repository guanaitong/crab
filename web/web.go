package web

import (
	"github.com/gin-gonic/gin"
	"github.com/guanaitong/crab/errors"
	"github.com/guanaitong/crab/system"
	"k8s.io/klog"
	"net/http"
	"reflect"
	"strconv"
)

var (
	codeKey        = http.CanonicalHeaderKey("x-error-code")
	msgKey         = http.CanonicalHeaderKey("x-error-msg")
	appNameKey     = http.CanonicalHeaderKey("x-app-name")
	appInstanceKey = http.CanonicalHeaderKey("x-app-instance")
)

func Success(c *gin.Context, data interface{}) {
	setHeader(c, errors.OK, "OK")
	if data == nil {
		c.String(http.StatusOK, "")
		c.Abort()
	}
	switch reflect.ValueOf(data).Kind() {
	case reflect.String:
		c.String(http.StatusOK, data.(string))
	case reflect.Bool:
		c.String(http.StatusOK, strconv.FormatBool(data.(bool)))
	case reflect.Int:
		c.String(http.StatusOK, strconv.FormatInt(int64(data.(int)), 10))
	case reflect.Int8:
		c.String(http.StatusOK, strconv.FormatInt(int64(data.(int8)), 10))
	case reflect.Int16:
		c.String(http.StatusOK, strconv.FormatInt(int64(data.(int16)), 10))
	case reflect.Int32:
		c.String(http.StatusOK, strconv.FormatInt(int64(data.(int32)), 10))
	case reflect.Int64:
		c.String(http.StatusOK, strconv.FormatInt(data.(int64), 10))
	case reflect.Uint:
		c.String(http.StatusOK, strconv.FormatUint(uint64(data.(uint)), 10))
	case reflect.Uint8:
		c.String(http.StatusOK, strconv.FormatUint(uint64(data.(uint8)), 10))
	case reflect.Uint16:
		c.String(http.StatusOK, strconv.FormatUint(uint64(data.(uint16)), 10))
	case reflect.Uint32:
		c.String(http.StatusOK, strconv.FormatUint(uint64(data.(uint32)), 10))
	case reflect.Uint64:
		c.String(http.StatusOK, strconv.FormatUint(data.(uint64), 10))
	case reflect.Float32:
		c.String(http.StatusOK, strconv.FormatFloat(float64(data.(float32)), 'g', -1, 32))
	case reflect.Float64:
		c.String(http.StatusOK, strconv.FormatFloat(data.(float64), 'g', -1, 64))
	default:
		c.JSON(http.StatusOK, data)
	}
	c.Abort()
}

func Fail(c *gin.Context, err errors.Error) {
	klog.Warningln("Fail--> " + err.Error())
	setHeader(c, err.ErrorCode(), err.ErrorMsg())
	c.Abort()
}

func setHeader(c *gin.Context, code int, message string) {
	c.Header(codeKey, strconv.Itoa(code))
	c.Header(msgKey, message)
	c.Header(appNameKey, system.GetAppName())
	c.Header(appInstanceKey, system.GetAppInstance())
	c.Status(http.StatusOK)
}
