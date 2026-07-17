package call

import (
	"fmt"
	"runtime"
	"time"

	"github.com/liangdas/mqant/log"
)

// Go 用协程处理f
func Go(f func()) {
	go func() {
		defer Recover()
		f()
	}()
}

func IndexTryCallback(f func() error, cb func()) {
	Go(func() {
		for i := 0; i <= 6; i++ {
			if i > 0 {
				next := time.Duration(i*i) * time.Minute
				time.Sleep(next)
			}
			err := f()
			if err == nil {
				break
			}
			log.Error("err:%v next:%v", err, i)
			if i == 6 {
				cb()
			}
		}
	})
}

func IndexTry(f func() error) {
	Go(func() {
		for i := 0; i <= 6; i++ {
			if i > 0 {
				next := time.Duration(i*i) * time.Minute
				time.Sleep(next)
			}
			err := f()
			if err == nil {
				break
			}
			log.Error("err:%v next:%v", err, i)
		}
	})
}

func IndexTryS(f func() error) {
	Go(func() {
		for i := 0; i <= 6; i++ {
			if i > 0 {
				next := time.Duration(i*i) * time.Second
				time.Sleep(next)
			}
			err := f()
			if err == nil {
				break
			}
			log.Error("err:%v next:%v", err, i)
		}
	})
}

func IndexTryWithParam(interval []int, f func() error) {
	Go(func() {
		err := f()
		if err == nil {
			return
		}
		for i := 0; i < len(interval); i++ {
			if i >= 0 {
				next := time.Duration(interval[i]) * time.Second
				time.Sleep(next)
			}
			err := f()
			if err == nil {
				break
			}
			log.Error("err:%v next:%v", err, interval[i])
		}
	})
}

// 捕获异常并打日志
// Usage: defer Recover()
func Recover() {
	if err := recover(); err != nil {
		buf := make([]byte, 1024)
		runtime.Stack(buf, false)
		str := fmt.Sprintf("%+v", err)
		log.Error("panic(%s), stack:\n%s", str, string(buf))
	}
}
