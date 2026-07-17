package call

var (
	execQueue = NewQueue()
)

// 顺序执行并发操作
func InitExecQueue() {
	execQueue.Run()
}

func AddQueue(f func()) {
	execQueue.AddQueue(f)
}

type ExecQueue struct {
	c chan func()
}

func (e *ExecQueue) Run() {
	Go(func() {
		for {
			f := <-e.c
			func() {
				defer Recover()
				f()
			}()
		}
	})
}

func (e *ExecQueue) AddQueue(f func()) {
	e.c <- f
}

func NewQueue(size ...int) *ExecQueue {
	e := &ExecQueue{}
	if size != nil {
		e.c = make(chan func(), size[0])
	} else {
		e.c = make(chan func(), 10000)
	}
	e.Run()
	return e
}
