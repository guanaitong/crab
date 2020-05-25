package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/guanaitong/crab/util/alert"
	"github.com/guanaitong/crab/util/runtime"
)

type Executor interface {
	Name() string
	Period() int
	Execute()
	Release(wg *sync.WaitGroup)
}

func StartBackgroundJob(ctx context.Context, executor Executor) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(time.Second * time.Duration(executor.Period())):
				func() {
					defer runtime.HandleCrashWithConfig(false, func(r interface{}) {
						callers := runtime.GetCallers(r)
						msg := fmt.Sprintf("GoroutineName:%s,\nObserved a panic: %#v (%v)\n%v", executor.Name(), r, r, callers)
						alert.SendByAppName(4, msg)
					})
					executor.Execute()
				}()
			}
		}
	}()
}

func StartAsyncJob(executor Executor) {
	go func() {
		func() {
			defer runtime.HandleCrashWithConfig(true, func(r interface{}) {
				callers := runtime.GetCallers(r)
				msg := fmt.Sprintf("GoroutineName:%s,\nObserved a panic: %#v (%v)\n%v", executor.Name(), r, r, callers)
				alert.SendByAppName(4, msg)
			})
			executor.Execute()
		}()
	}()
}
