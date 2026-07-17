package pay

import (
	"server/call"
	"server/common"
	"server/db"
	"server/modules/pay/values"
)

func loadConfig() error {
	c := map[int][]func() error{}
	c[common.ReloadPayChannel] = []func() error{func() error { return LoadPayWeight() }}
	c[common.ReloadWithdrawChannel] = []func() error{func() error { return LoadWithdrawWeight() }}
	call.LoadConfigs(c)
	return nil
}

func LoadPayMethod() error {
	values.PayStatus, _ = db.Redis().GetInt(common.RedisKeyPayStatus)
	return nil
}

func LoadPayWeight() error {
	values.PayWeightLock.Lock()
	defer values.PayWeightLock.Unlock()
	call.LoadConfigPayChannels()
	return nil
}

func LoadWithdrawWeight() error {
	values.WithdrawWeightLock.Lock()
	defer values.WithdrawWeightLock.Unlock()
	call.LoadConfigWithdrawChannels()
	return nil
}
