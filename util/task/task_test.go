package task_test

import (
	"context"
	"fmt"
	"github.com/guanaitong/crab/util/task"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	t.Log("pass")
}

func TestStartRunner1(t *testing.T) {
	var end = 0
	r := task.NewWorkGroup(context.Background(), "test", time.Second, func(ctx context.Context) {
		fmt.Println(ctx.Value(task.GoroutineIndexKey))
		time.Sleep(5 * time.Second)
		fmt.Println("end")
		end++
	}, 4)
	r.Start()
	time.Sleep(time.Second)
	r.Stop()
	if end != 4 {
		t.Fail()
	}
}

func TestStartRunner2(t *testing.T) {
	var end = 0
	ch := make(chan int, 2)
	r := task.NewWorkGroup(context.Background(), "test", time.Second, func(ctx context.Context) {
		fmt.Println(ctx.Value(task.GoroutineIndexKey))
		for i := range ch {
			fmt.Printf("num_%d \n", i)
			time.Sleep(time.Second)
		}
		fmt.Println("end")
		end++
	}, 2)
	r.Start()
	time.Sleep(time.Second)
	ch <- 3
	ch <- 3
	ch <- 3
	ch <- 3
	ch <- 2
	ch <- 2
	ch <- 2
	ch <- 1
	close(ch)

	r.Stop()
	if end != 2 {
		t.Fail()
	}
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
		t.Errorf("not workerGroup")
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
		t.Errorf("not workerGroup")
	}
}*/
