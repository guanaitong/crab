package system_test

import (
	"github.com/guanaitong/crab/system"
	"testing"
)

func TestSetup(t *testing.T) {
	system.Setup("TestApp", 99)
	t.Log(system.GetAppName())
	t.Log(system.GetErrCodePrefix())
	t.Log(system.GetEnv("none"))
}
