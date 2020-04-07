package router

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

var (
	defaultSkipPaths = []string{"/isLive"}
)

type msgFmt struct {
	RequestTime        string        // 请求时间
	ThreadName         string        // 处理线程
	AccessIp           string        // 访问者的IP
	AppName            string        // 应用名
	AppInstance        string        // 应用实例
	RequestAppName     string        // 请求应用名
	RequestAppInstance string        // 请求应用实例
	TraceId            string        // 分布式jaeger追踪系统里的traceId
	SpanId             string        // 分布式jaeger追踪系统里的spanId
	ParentId           string        // 分布式jaeger追踪系统里的parentId
	Method             string        // http请求方法
	URL                string        // http请求URL里的路径
	Protocol           string        // http请求协议版本
	StatusCode         int           // http请求返回码
	BodySize           int           // http请求发送信息的字节数，不包括http头，如果字节数为0的话，显示为-
	Latency            time.Duration // http请求处理消耗时间，单位毫秒
	Query              string        // http请求URL里的query
}

type logger struct {
	startTime  int64
	filename   string
	filePath   string
	format     string
	writer     *os.File
	skipPaths  []string
	lock       sync.Mutex
	logChan    chan *msgFmt
	signalChan chan struct{}
	wait       sync.WaitGroup
}

type Logger interface {
	Init() error
	Write(s *msgFmt)
	Flush()
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
	var (
		signalChan = make(chan struct{}, 1)
		logChan    = make(chan *msgFmt, 10)
		wait       = sync.WaitGroup{}
		lock       = sync.Mutex{}
	)

	l := &logger{
		filename:   defaultAccessLoggerFile,
		filePath:   defaultLoggerPath,
		format:     defaultLoggerFormat,
		skipPaths:  defaultSkipPaths,
		lock:       lock,
		logChan:    logChan,
		signalChan: signalChan,
		wait:       wait,
	}

	go func() {
		for {
			l.handle()
		}
	}()

	return l
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

func (l *logger) handle() {
	for {
		select {
		case msg := <-l.logChan:
			l.write(msg)
			l.wait.Done()
		case <-l.signalChan:
			l.flush()
		}
	}
}

func (l *logger) write(msg *msgFmt) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	var (
		filename       = l.filename
		filePath       = l.filePath
		filenameSuffix = path.Ext(l.filename)
		startTime      = time.Unix(l.startTime, 0)
		now            = time.Now()
		timeFormat     = l.format
		oldFilename    string
		err            error
	)

	if startTime.Format(timeFormat) != now.Format(timeFormat) {
		l.flush()
		oldFilename = strings.Replace(filename, filenameSuffix, "", 1) + "." + startTime.Format(timeFormat) + filenameSuffix
		if err = syscall.Rename(path.Join(filePath, filename), path.Join(filePath, oldFilename)); err != nil {
			return err
		}
		if err = l.Init(); err != nil {
			return err
		}
	}

	log := defaultLogFormatter(msg)
	_, err = l.writer.Write([]byte(log))

	return err
}

func (l *logger) flush() {
	for {
		if len(l.logChan) > 0 {
			msg := <-l.logChan
			l.write(msg)
			l.wait.Done()
		}
		break
	}
	l.writer.Close()
}

func (l *logger) Write(msg *msgFmt) {
	l.wait.Add(1)
	l.logChan <- msg
}

func (l *logger) Flush() {
	l.signalChan <- struct{}{}
	l.wait.Wait()
}

// access log
func (l *logger) Logger() gin.HandlerFunc {
	notlogged := l.skipPaths
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notlogged {
			skip[p] = struct{}{}
		}
	}

	appName := system.GetAppName()
	appInstance := system.GetAppInstance()

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		urlPath := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		if _, ok := skip[urlPath]; !ok {
			timestamp := time.Now()
			param := &msgFmt{
				RequestTime:        timestamp.Format("2006/01/02:15:04:05"),
				ThreadName:         "-",
				AccessIp:           c.ClientIP(),
				AppName:            appName,
				AppInstance:        appInstance,
				RequestAppName:     "-",
				RequestAppInstance: "-",
				TraceId:            "-",
				SpanId:             "-",
				ParentId:           "-",
				Method:             c.Request.Method,
				URL:                urlPath,
				Protocol:           c.Request.Proto,
				StatusCode:         c.Writer.Status(),
				BodySize:           c.Writer.Size(),
				Latency:            timestamp.Sub(start),
				Query:              raw,
			}

			l.Write(param)
		}
	}
}

// defaultLogFormatter is the default log format function Logger middleware uses.
func defaultLogFormatter(param *msgFmt) string {
	return fmt.Sprintf("[%s]^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^%d^%v^%d^%s\n",
		param.RequestTime,
		param.ThreadName,
		param.AccessIp,
		param.AppName,
		param.AppInstance,
		param.RequestAppName,
		param.RequestAppInstance,
		param.TraceId,
		param.SpanId,
		param.ParentId,
		param.Method,
		param.URL,
		param.Protocol,
		param.StatusCode,
		param.BodySize,
		param.Latency,
		param.Query,
	)
}

// application log
func ApplicationLogger() {
	_ = os.Mkdir(defaultLoggerPath, 0755)
	klog.InitFlags(nil)
	_ = flag.Set("logtostderr", "false")
	_ = flag.Set("alsologtostderr", "true")
	_ = flag.Set("log_file", path.Join(defaultLoggerPath, defaultApplicationLoggerFile))
}

func ApplicationLoggerFlush() {
	klog.Flush()
}
