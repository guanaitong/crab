package web

import (
	"github.com/gin-gonic/gin"
	"github.com/guanaitong/crab/system"
	"net/http"
	"strings"
)

var (
	appNameKey     = http.CanonicalHeaderKey("x-app-name")
	appInstanceKey = http.CanonicalHeaderKey("x-app-instance")
)

const (
	defaultLoggerPath            = "logs"
	defaultApplicationLoggerFile = "application.log"
	defaultAccessLoggerFile      = "access.log"
	defaultLoggerFormat          = "2006-01-02"
)

var (
	defaultSkipPaths = []string{"/isLive"}
)

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	Write(c, http.StatusOK, data)
}

func Write(c *gin.Context, code int, data interface{}) {
	c.Header(appNameKey, system.GetAppName())
	c.Header(appInstanceKey, system.GetAppInstance())
	c.JSON(code, data)
	c.Abort()
}

func Default() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	logger := NewLogger()

	if err := logger.Init(); err != nil {
		engine.Use(gin.Logger())
	} else {
		engine.Use(logger.Logger())
	}
	return engine
}

type Controller interface {
	RequestMappings() []Handler
}

type GinGrouper interface {
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}

type Handler struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

func Setup(rootPath string, r GinGrouper, controller Controller) {
	rg := r.Group(rootPath)
	handleMethods := controller.RequestMappings()
	for _, handleMethod := range handleMethods {
		method := handleMethod.Method
		if method == "*" {
			rg.Any(handleMethod.Path, handleMethod.HandlerFunc)
			continue
		}
		if method == "" {
			method = "GET"
		}
		methods := strings.Split(method, ",")
		for _, m := range methods {
			rg.Handle(m, handleMethod.Path, handleMethod.HandlerFunc)
		}
	}
}
