package web

import (
	"github.com/gin-gonic/gin"
	"github.com/guanaitong/crab/errors"
	"github.com/guanaitong/crab/system"
	"net/http"
	"strconv"
)

var (
	codeKey        = http.CanonicalHeaderKey("x-error-code")
	msgKey         = http.CanonicalHeaderKey("x-error-msg")
	appNameKey     = http.CanonicalHeaderKey("x-app-name")
	appInstanceKey = http.CanonicalHeaderKey("x-app-instance")
)

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	Write(c, &ApiResponse{Code: errors.OK, Msg: "OK", Data: data})
}

func Fail(c *gin.Context, err errors.Error) {
	Write(c, &ApiResponse{Code: err.ErrorCode(), Msg: err.ErrorMsg(), Data: nil})
}

func Write(c *gin.Context, apiResponse *ApiResponse) {
	c.Header(codeKey, strconv.Itoa(apiResponse.Code))
	c.Header(msgKey, apiResponse.Msg)
	c.Header(appNameKey, system.GetAppName())
	c.Header(appInstanceKey, system.GetAppInstance())
	c.JSON(http.StatusOK, apiResponse)
	c.Abort()
}
