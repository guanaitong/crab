package alert_test

import (
	"github.com/guanaitong/crab/util/alert"
	"testing"
	"time"
)

func TestSendByAppName(t *testing.T) {
	alert.SendByAppName(1, "TestSendByAppName")
	time.Sleep(time.Second * 5)
}

func TestSendByCorpCodes(t *testing.T) {
	alert.SendByCorpCodes(1, "TestSendByCorpCodes", "HB266")
	alert.SendByCorpCodes(7, "TestSendByCorpCodes1", "HB266", "HB533")
	time.Sleep(time.Second * 5)
}

func TestSendByGroupId(t *testing.T) {
	alert.SendByGroupId(1, "TestSendByGroupId", 4)
	time.Sleep(time.Second * 5)
}
