package system_test

import (
	"github.com/guanaitong/crab/system"
	"testing"
)

func TestSetup(t *testing.T) {
	system.SetupAppName("TestApp")
	t.Log(system.GetAppName())
	t.Log(system.GetEnv("none"))
}
