package call

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	mqrpc "github.com/liangdas/mqant/rpc"
	"github.com/liangdas/mqant/selector"
)

var (
	caller module.Module
)

// NewCaller 新建一个调用对象
func NewCaller(m module.Module) {
	caller = m
}

func GetCaller() module.Module {
	return caller
}

// Publish 发布消息
func Publish(topic string, data proto.Message) error {
	send, _ := proto.Marshal(data)
	err := caller.GetApp().Transport().Publish(topic, send)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// Publish 发布消息
func PublishRequest(topic, reply string, data proto.Message) error {
	send, _ := proto.Marshal(data)
	return caller.GetApp().Transport().PublishRequest(topic, reply, send)
}

// RPC 服务器内部调用
func RPC(module, method string, param proto.Message, opts ...selector.SelectOption) (interface{}, string) {
	return caller.GetApp().Call(context.Background(), module, method, mqrpc.Param(param), opts...)
}
