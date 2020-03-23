package web

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"github.com/guanaitong/crab/errors"
	"github.com/guanaitong/crab/system"
)

var (
	codeKey        = http.CanonicalHeaderKey("x-error-code")
	msgKey         = http.CanonicalHeaderKey("x-error-msg")
	appNameKey     = http.CanonicalHeaderKey("x-app-name")
	appInstanceKey = http.CanonicalHeaderKey("x-app-instance")
)

const (
	defaultLoggerPath            = "logs"
	defaultApplicationLoggerFile = "application.log"
	defaultAccessLoggerFile      = "access.log"
	defaultLoggerFormat          = "2006-01-02"
)

type logger struct {
	startTime int64
	filename  string
	filePath  string
	format    string
	writer    *os.File
	lock      sync.Mutex
}

type Logger interface {
	Init() error
	Write(s string) error
	Logger() gin.HandlerFunc
}

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

func NewLogger() Logger {
	return &logger{
		filename: defaultAccessLoggerFile,
		filePath: defaultLoggerPath,
		format:   defaultLoggerFormat,
		lock:     sync.Mutex{},
	}
}

func (l *logger) Init() error {
	var (
		writer *os.File
		err    error
	)

	if _, err = os.Stat(l.filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(l.filePath, 0755); err != nil {
			return err
		}
	}
	writer, err = os.OpenFile(path.Join(l.filePath, l.filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	l.writer = writer
	l.startTime = time.Now().Unix()
	return nil
}

func (l *logger) Write(s string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	var (
		filename       = l.filename
		filePath       = l.filePath
		filenameSuffix = path.Ext(l.filename)
		startTime      = time.Unix(l.startTime, 0)
		now            = time.Now()
		format         = l.format
		oldFilename    string
	)

	if startTime.Format(format) == now.Format(format) {
		l.writer.Close()
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "." + startTime.Format(l.format) + filenameSuffix
		err := syscall.Rename(path.Join(filePath, filename), path.Join(filePath, oldFilename))
		if err != nil {
			return err
		}
		if err = l.Init(); err != nil {
			return err
		}
	}
	_, err := l.writer.Write([]byte(s))
	return err
}

// access log
func (l *logger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		urlPath := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			urlPath = urlPath + "?" + raw
		}

		param.Path = urlPath
		l.Write(defaultLogFormatter(param))
	}
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// application log
func ApplicationLogger() {
	klog.InitFlags(nil)
	_ = flag.Set("logtostderr", "false")
	_ = flag.Set("alsologtostderr", "true")
	_ = flag.Set("log_file", path.Join(defaultLoggerPath, defaultApplicationLoggerFile))
}

func ApplicationLoggerFlush() {
	klog.Flush()
}
