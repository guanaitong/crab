package task_test

import "testing"

func TestDefault(t *testing.T) {
	t.Log("pass")
}

/*
import (
	"github.com/guanaitong/crab/util/task"
	"testing"
	"time"
)

func TestStartBackgroundTask(t *testing.T) {
	count := 0
	task.StartBackgroundTask("test1", time.Second, func() {
		count = count + 1
		t.Log("123")
	})

	time.Sleep(time.Second * 3)

	if count < 3 {
		t.Errorf("not work")
	}
}

func TestStartBackgroundTaskCrash(t *testing.T) {
	count := 0
	task.StartBackgroundTask("test1", time.Second, func() {
		count = count + 1
		panic("123")
	})
	time.Sleep(time.Second * 5)
	if count < 3 {
		t.Errorf("not work")
	}
}*/
