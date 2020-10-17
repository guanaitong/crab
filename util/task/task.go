package task

import (
	"context"
	"fmt"
	"github.com/guanaitong/crab/alert"
	"github.com/guanaitong/crab/util/runtime"
	"k8s.io/klog"
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

type WorkGroup interface {
	Start()
	Stop()
}

// 因为：context的CancelFunc不会等待任务真正停止
// 所以有workerGroup的封装：workerGroup的stop方法，会同步等待任务真正stop
// 一个workerGroup对应一组goroutine，数量为parallelNum
type workerGroup struct {
	name         string
	period       time.Duration
	taskFunc     func(ctx context.Context)
	stopSignalCh chan struct{}
	ctx          context.Context
	cancel       context.CancelFunc
	parallelNum  int
}

func (worker *workerGroup) Start() {
	klog.Infof("workerGroup [%s] starting", worker.name)
	for i := 0; i < worker.parallelNum; i++ {
		worker.startOneGoroutine(i)
	}
	klog.Infof("workerGroup [%s] started", worker.name)
}

const GoroutineIndexKey = "goroutine_index"

func (worker *workerGroup) startOneGoroutine(i int) {
	go func() {
		goroutineName := fmt.Sprintf("%s__%d", worker.name, i)
		for {
			select {
			default:
				func() {
					defer runtime.HandleCrashWithConfig(false, func(r interface{}) {
						callers := runtime.GetCallers(r)
						msg := fmt.Sprintf("GoroutineName:%s,\nObserved a panic: %#v (%v)\n%v", goroutineName, r, r, callers)
						alert.SendByAppName(4, msg)
					})
					// 这个把goroutine的index放到ctx里
					ctx := context.WithValue(worker.ctx, GoroutineIndexKey, i)
					worker.taskFunc(ctx)
					time.Sleep(worker.period)
				}()
			case <-worker.ctx.Done():
				// 收到ctx的停止信号后，给stoppedSignal发一个信号
				worker.stopSignalCh <- struct{}{}
				klog.Infof("workerGroup [%s] goroutine_index_[%d] stopped", worker.name, i)
				return
			}
		}
	}()
	klog.Infof("workerGroup [%s] goroutine_index_[%d] started", worker.name, i)
}

// stop方法会等待taskFunc彻底退出
func (worker *workerGroup) Stop() {
	klog.Infof("workerGroup [%s] stopping", worker.name)
	worker.cancel()
	// 等待真正停止信号
	stopNum := 0
	for range worker.stopSignalCh {
		stopNum++
		if stopNum == worker.parallelNum {
			break
		}
	}
	close(worker.stopSignalCh)
	klog.Infof("workerGroup [%s] stopped", worker.name)
}

func NewWorkGroup(
	parentCtx context.Context,
	name string,
	period time.Duration,
	taskFunc func(ctx context.Context),
	parallelNum int) WorkGroup {
	if parallelNum < 1 {
		panic("parallelNum must be positive")
	}
	ctx, cancel := context.WithCancel(parentCtx)
	r := &workerGroup{name: name,
		period:       period,
		taskFunc:     taskFunc,
		stopSignalCh: make(chan struct{}, parallelNum),
		ctx:          ctx,
		cancel:       cancel,
		parallelNum:  parallelNum}

	return r
}
