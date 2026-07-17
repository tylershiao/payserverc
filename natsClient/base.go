package natsClient

import (
	"runtime"
	"server/util"
	"time"

	"github.com/liangdas/mqant/log"
	"github.com/nats-io/nats.go"
)

type NatsImp interface {
	// 收到消息后的回调
	OnMsgCallBack(*nats.Msg)
}

type NatsClient struct {
	NatsImp NatsImp
	*nats.Conn
	done    chan bool
	Topic   string
	Queue   string
	isClose bool
	GO      bool
}

func NewNatsClient() *NatsClient {
	n := new(NatsClient)
	n.done = make(chan bool)
	return n
}

func (nc *NatsClient) OnRequestHandle() error {
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r := r.(type) {
			case string:
				rn = r
			case error:
				rn = r.Error()
			}
			buf := make([]byte, 1024)
			l := runtime.Stack(buf, false)
			errstr := string(buf[:l])
			log.Error("%s\n ----Stack----\n%s", rn, errstr)
		}
	}()
	var subs *nats.Subscription
	var err error
	if len(nc.Queue) == 0 {
		subs, err = nc.SubscribeSync(nc.Topic)
	} else {
		subs, err = nc.QueueSubscribeSync(nc.Topic, nc.Queue)
	}
	if err != nil {
		return err
	}
	util.Go(func() {
		//服务关闭
		<-nc.done
		subs.Unsubscribe()
	})
	for !nc.isClose {
		m, err := subs.NextMsg(time.Minute)
		if err != nil && err == nats.ErrTimeout {
			if !subs.IsValid() {
				//订阅已关闭，需要重新订阅
				subs, err = nc.SubscribeSync(nc.Topic)
				if err != nil {
					log.Error("NatsClient SubscribeSync[1] error with '%v'", err)
					continue
				}
			}
			continue
		} else if err != nil {
			if !subs.IsValid() {
				//订阅已关闭，需要重新订阅
				subs, err = nc.SubscribeSync(nc.Topic)
				if err != nil {
					log.Error("NatsClient SubscribeSync[1] error with '%v'", err)
					continue
				}
			}
			log.Warning("NatsServer error with '%v'", err)
			continue
		}
		if nc.GO {
			util.Go(func() {
				nc.NatsImp.OnMsgCallBack(m)
			})
		} else {
			nc.NatsImp.OnMsgCallBack(m)
		}
	}
	return nil
}
