package natsClient

import (
	"runtime"
	"time"

	"github.com/liangdas/mqant/log"
	"github.com/nats-io/nats.go"
)

// CommonNatsImp 通用回调natsClient
type CommonNatsImp struct {
	f func(data []byte)
}

// g 是否用协程处理
func NewCommonNatsImp(conn *nats.Conn, topic string, f func(data []byte)) {
	base := NewNatsClient()
	base.Topic = topic
	r := &CommonNatsImp{f: f}
	base.NatsImp = r
	base.Conn = conn
	base.GO = true
	go func() {
		err := base.OnRequestHandle()
		if err != nil {
			log.Error("newClientDissconnectNatsImp OnRequestHandle error:%v", err)
		}
	}()
}

func (r *CommonNatsImp) OnMsgCallBack(m *nats.Msg) {
	r.f(m.Data)
}

// NewReplyNatsImp 回复监听
func NewReplyNatsImp(conn *nats.Conn, topic string) error {
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
	subs, err := conn.SubscribeSync(topic)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	defer subs.Unsubscribe()
	_, err = subs.NextMsg(5 * time.Second)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
	// one := &pb.InnerReply{}
	// if err := proto.Unmarshal(m.Data, one); err != nil {
	// 	log.Error("err:%v", err)
	// 	return err
	// }
	// if one.Code == 0 {
	// 	return nil
	// }
	// return errors.New(one.Msg)
}
