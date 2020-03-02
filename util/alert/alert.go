package alert

import (
	"bytes"
	"encoding/json"
	"github.com/guanaitong/crab/hc"
	"github.com/guanaitong/crab/system"
	"github.com/guanaitong/crab/util/runtime"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl = "http://message.frigate.devops.wuxingdev.cn"

	// way = 0
	byAppNameUrl = baseUrl + "/v2/message/sendMsgByAppNames"

	// way = 1
	byGroupUrl = baseUrl + "/v2/message/sendMsgByGroups"

	// way = 2
	byQiWeiXinUrl = baseUrl + "/v2/message/sendMsgByWeChatIds"

	bufferLen = 4096
)

var (
	ch = make(chan *FrigateMessage, bufferLen)
)

func init() {
	go func() {
		for {
			func() {
				defer runtime.HandleCrashWithConfig(false)
				for data := range ch {
					e := send(data)
					if e != nil {
						log.Printf("send fail,data:%v, error:%s", data, e.Error())
					}
				}

			}()
		}
	}()
}

type FrigateMessage struct {
	//发送渠道，默认通过1为企业微信通知
	Channel int `json:"channel"`
	//消息标题
	Title string `json:"title"`
	//消息内容
	Content string `json:"content"`
	//当有异常堆栈时，堆栈内容
	Stack string `json:"stack"`
	//模块
	Module string `json:"module"`
	//标签
	Tags map[string]string `json:"tags"`

	// ------------------以下属于系统变量------------------------
	TraceId     string `json:"traceId"`
	HostIp      string `json:"hostIp"`
	AppName     string `json:"appName"`
	AppInstance string `json:"appInstance"`
	WorkEnv     string `json:"workEnv"`
	WorkIdc     string `json:"workIdc"`
	//发送时间
	Time   int64 `json:"time"`
	Format bool  `json:"format"`

	receiveInfo pair
	way         int8
}

func (m *FrigateMessage) Send() {
	ch <- m
}

type pair struct {
	Key   string
	Value string
}

func NewMessage() *FrigateMessage {
	return &FrigateMessage{
		Title:       "frigate 消息通知",
		AppName:     system.GetAppName(),
		AppInstance: system.GetAppInstance(),
		HostIp:      system.GetHostIp(),
		WorkEnv:     system.GetWorkEnv(),
		WorkIdc:     system.GetWorkIdc(),
		Time:        time.Now().Unix() * 1000,
	}
}

// 通过工号发送消息
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByCorpCodes(channel int, msg string, corpCodes ...string) {
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	m := NewMessage()
	m.receiveInfo = pair{
		Key:   "receiveWeChatIds",
		Value: strings.Join(corpCodes, ","),
	}
	m.Channel = channel
	m.Content = msg
	m.way = 2
	m.Send()
}

// 通过组发送消息
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByGroupId(channel int, msg string, groupId int) {
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	m := NewMessage()
	m.receiveInfo = pair{
		Key:   "receiveGroups",
		Value: strconv.Itoa(groupId),
	}
	m.Channel = channel
	m.Content = msg
	m.way = 1
	m.Send()
}

// 通过应用名发送消息，自动获取应用名，不需要传递
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByAppName(channel int, msg string) {
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	m := NewMessage()
	m.receiveInfo = pair{
		Key:   "receiveAppNames",
		Value: system.GetAppName(),
	}
	m.Channel = channel
	m.Content = msg
	m.way = 0
	m.Send()
}

func send(message *FrigateMessage) error {
	way := message.way

	byUrl := ""
	if way == 0 {
		byUrl = byAppNameUrl
	} else if way == 1 {
		byUrl = byGroupUrl
	} else if way == 2 {
		byUrl = byQiWeiXinUrl
	}

	byUrl = byUrl + "?" + url.QueryEscape(message.receiveInfo.Key) + "=" + url.QueryEscape(message.receiveInfo.Value)

	data, _ := json.Marshal(message)

	req, err := http.NewRequest("POST", byUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "GOLANG_UTIL")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := hc.Request(req)

	if err != nil {
		log.Printf("send request fail, error:%s", err)
		return err
	}

	if resp.StatusCode() == 200 {
		return nil
	}

	log.Printf("send request fail, status:%d,body:%s", resp.StatusCode(), resp.AsString())

	return nil
}
