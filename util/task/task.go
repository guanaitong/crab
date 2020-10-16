package task

import (
	"context"
	"fmt"
	"github.com/guanaitong/crab/alert"
	"github.com/guanaitong/crab/util/runtime"
	"time"
)

// 任务会永远的运行下去
func StartBackgroundTask(name string, period time.Duration, task func()) {
	go func() {
		for {
			func() {
				defer runtime.HandleCrashWithConfig(false, func(r interface{}) {
					callers := runtime.GetCallers(r)
					msg := fmt.Sprintf("GoroutineName:%s,\nObserved a panic: %#v (%v)\n%v", name, r, r, callers)
					alert.SendByAppName(4, msg)
				})
				task()
				time.Sleep(period)
			}()
		}
	}()
}
// when stopped,stoppedSignal will receive signal
func StartStopAbleBackgroundTask(ctx context.Context, stoppedSignal chan struct{}, name string, period time.Duration, task func()) {
	go func() {
		for {
			select {
			default:
				func() {
					defer runtime.HandleCrashWithConfig(false, func(r interface{}) {
						callers := runtime.GetCallers(r)
						msg := fmt.Sprintf("GoroutineName:%s,\nObserved a panic: %#v (%v)\n%v", name, r, r, callers)
						alert.SendByAppName(4, msg)
					})
					task()
					time.Sleep(period)
				}()
			case <-ctx.Done():
				stoppedSignal <- struct{}{}
				return
			}
		}
	}()
}

// 任务只会执行一次
func StartAsyncTask(name string, task func()) {
	go func() {
		func() {
			defer runtime.HandleCrashWithConfig(true, func(r interface{}) {
				callers := runtime.GetCallers(r)
				msg := fmt.Sprintf("GoroutineName:%s,\nObserved a panic: %#v (%v)\n%v", name, r, r, callers)
				alert.SendByAppName(4, msg)
			})
			task()
		}()
	}()
}
