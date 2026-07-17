package call

import (
	"server/common"

	"github.com/liangdas/mqant/log"
)

var (
	reloadFuncs = map[int][]func() error{}
)

func LoadConfigs(funcs map[int][]func() error) error {
	CommonReload(funcs)
	for _, v := range funcs {
		for _, k := range v {
			if err := k(); err != nil {
				log.Error("err:%v", err)
				// return err
			}
		}
	}
	reloadFuncs = funcs
	return nil
}

func ReloadConfig(t int) {
	log.Debug("reload %v", t)
	if fs, ok := reloadFuncs[t]; ok {
		for _, v := range fs {
			if err := v(); err != nil {
				log.Error("err:%v", err)
				return
			}
		}
	}
}

func CommonReload(c map[int][]func() error) {
	if _, ok := c[common.ReloadPayChannel]; !ok {
		c[common.ReloadPayChannel] = []func() error{func() error {
			if err := LoadConfigPayChannels(); err != nil {
				log.Error("err:%v", err)
				return err
			}
			return nil
		}}
	}
	if _, ok := c[common.ReloadWithdrawChannel]; !ok {
		c[common.ReloadWithdrawChannel] = []func() error{func() error {
			if err := LoadConfigWithdrawChannels(); err != nil {
				log.Error("err:%v", err)
				return err
			}
			return nil
		}}
	}
	if _, ok := c[common.ReloadMerchats]; !ok {
		c[common.ReloadMerchats] = []func() error{func() error {
			if err := LoadConfigMerchats(); err != nil {
				log.Error("err:%v", err)
				return err
			}
			return nil
		}}
	}
}
