package natsClient

import (
	"strconv"

	"github.com/liangdas/mqant/log"
	"github.com/nats-io/nats.go"
)

type ReloadNatsImp struct {
	f func(int)
}

// 新建一个reload监听
func NewReloadNats(conn *nats.Conn, f func(int)) {
	base := NewNatsClient()
	base.Topic = TopicReloadConfig
	r := &ReloadNatsImp{f: f}
	base.NatsImp = r
	base.Conn = conn
	go func() {
		err := base.OnRequestHandle()
		if err != nil {
			log.Error("newClientDissconnectNatsImp OnRequestHandle error:%v", err)
		}
	}()
}

func (r *ReloadNatsImp) OnMsgCallBack(m *nats.Msg) {
	// req := &pb.ReloadGameConfig{}
	// err := proto.Unmarshal(m.Data, req)
	// if err != nil {
	// 	log.Error("ReloadNatsImp OnMsgCallBack error %v", err)
	// 	return
	// }
	data := string(m.Data)
	t, err := strconv.Atoi(data)
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	log.Debug("ReloadNatsImp:%d", t)
	r.f(t)
}
