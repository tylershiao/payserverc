package call

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"server/common"
	"server/db"
	"server/util"
	"strings"

	"github.com/liangdas/mqant/log"
	"github.com/oschwald/geoip2-golang"
)

var (
	region                 *geoip2.Reader
	ConfigPayChannels      []*common.ConfigPayChannels
	ConfigWithdrawChannels []*common.ConfigWithdrawChannels
	ConfigMerchats         []*common.Merchats
)

// LoadIpDB 初始化ip数据库
func LoadIpDB() (err error) {
	if region != nil {
		region.Close()
	}
	region, err = geoip2.Open("../city.mmdb")
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	return
}

// SearchIP 查询ip
func SearchIP(ip string) (ipInfo *geoip2.City, err error) {
	if region == nil {
		if err = LoadIpDB(); err != nil {
			return
		}
	}
	pip := net.ParseIP(ip)
	return region.City(pip)
}

// LoadConfigPayChannels 加载代收渠道配置
func LoadConfigPayChannels() (err error) {
	var one []*common.ConfigPayChannels
	if _, err = db.Mysql().QueryAll("", "weight desc", &common.ConfigPayChannels{}, &one); err != nil {
		log.Error("err:%v", err)
		return err
	}
	ret := []*common.ConfigPayChannels{}
	for _, v := range one {
		if v.Weight <= 0 {
			continue
		}
		ret = append(ret, v)
		log.Debug("finish load ConfigPayChannels:%+v", *v)
	}
	ConfigPayChannels = ret
	return nil
}

// GetConfigPayChannels 获取代收渠道配置
func GetConfigPayChannels() []*common.ConfigPayChannels {
	ret := []*common.ConfigPayChannels{}
	for i, v := range ConfigPayChannels {
		if v.Weight <= 0 {
			continue
		}
		ret = append(ret, ConfigPayChannels[i])
	}
	return ret
}

func GetConfigPayChannelByCid(cid int) *common.ConfigPayChannels {
	for _, v := range ConfigPayChannels {
		if v.ChannelID == cid {
			return v
		}
	}
	return nil
}

// GetConfigPayChannelByCids 获取代收渠道配置
func GetConfigPayChannelByCids(ids []int, amount float64, bans ...int) []*common.ConfigPayChannels {
	ret := []*common.ConfigPayChannels{}
	for _, v := range ConfigPayChannels {
		if util.SliceContain(ids, v.ChannelID) {
			// 检查金额
			if amount < v.PayDown || amount > v.PayUp {
				continue
			}
			if bans != nil && util.SliceContain(bans, v.ChannelID) {
				continue
			}
			ret = append(ret, v)
		}
	}
	return ret
}

// LoadConfigWithdrawChannels 加载代付渠道配置
func LoadConfigWithdrawChannels() (err error) {
	var one []*common.ConfigWithdrawChannels
	if _, err = db.Mysql().QueryAll("", "weight desc", &common.ConfigWithdrawChannels{}, &one); err != nil {
		log.Error("err:%v", err)
		return err
	}
	ret := []*common.ConfigWithdrawChannels{}
	for _, v := range one {
		if v.Weight <= 0 {
			continue
		}
		ret = append(ret, v)
		log.Debug("finish load ConfigWithdrawChannels:%+v", *v)
	}
	ConfigWithdrawChannels = ret
	return nil
}

// GetConfigWithdrawChannels 获取代付渠道配置
func GetConfigWithdrawChannels() []*common.ConfigWithdrawChannels {
	ret := []*common.ConfigWithdrawChannels{}
	for i, v := range ConfigWithdrawChannels {
		if v.Weight <= 0 {
			continue
		}
		ret = append(ret, ConfigWithdrawChannels[i])
	}
	return ret
}

// GetConfigWithdrawChannelsByID 获取代付渠道配置
func GetConfigWithdrawChannelsByID(id int) *common.ConfigWithdrawChannels {
	for i, v := range ConfigWithdrawChannels {
		if v.ChannelID == id {
			return ConfigWithdrawChannels[i]
		}
	}
	return nil
}

// GetConfigWithdrawChannelsByID 获取代付渠道配置
func GetConfigWithdrawChannelsByIDs(ids []int) []*common.ConfigWithdrawChannels {
	ret := []*common.ConfigWithdrawChannels{}
	for _, v := range ConfigWithdrawChannels {
		if util.SliceContain(ids, v.ChannelID) {
			ret = append(ret, v)
		}
	}
	return ret
}

// GetConfigWithdrawChannelsBest 选择推荐的代付渠道
func GetConfigWithdrawChannelsBest(id int) *common.ConfigWithdrawChannels {
	if len(ConfigWithdrawChannels) > 0 {
		return ConfigWithdrawChannels[0]
	}
	return nil
}

// LoadConfigMerchats 加载商户配置
func LoadConfigMerchats() (err error) {
	var one []*common.Merchats
	if _, err = db.Mysql().QueryAll("", "", &common.Merchats{}, &one); err != nil {
		log.Error("err:%v", err)
		return err
	}
	for _, v := range one {
		if v.PayChannels != "" {
			// 左右增加[]
			str := "[" + v.PayChannels + "]"
			err := json.Unmarshal([]byte(str), &v.PayChannelsSub)
			if err != nil {
				log.Error("err:%v", err)
				return err
			}
		}
		if v.PayoutPayChannels != "" {
			// 左右增加[]
			str := "[" + v.PayoutPayChannels + "]"
			err := json.Unmarshal([]byte(str), &v.PayoutChannelsSub)
			if err != nil {
				log.Error("err:%v", err)
				return err
			}
		}
		if v.WhiteIP != "" {
			ips := strings.Split(v.WhiteIP, ",")
			for _, ip := range ips {
				ip = strings.TrimSpace(ip) // 去除空格
				normalizedIP := util.NormalizeIP(ip)
				if normalizedIP == "" {
					log.Error("invalid IP: %s\n", ip)
					return errors.New("invalid ip")
				}
				v.WhiteIPSub = append(v.WhiteIPSub, normalizedIP)
			}
		}
		if v.PayRange != "" {
			var min, max float64
			count, err := fmt.Sscanf(v.PayRange, "%f-%f", &min, &max)
			if err != nil {
				log.Error("err:%v", err)
				return err
			}
			if count != 2 {
				err = errors.New("invalid count")
				log.Error("err:%v", err)
				return err
			}
			if min > max {
				err = errors.New("invalid min max")
				log.Error("err:%v", err)
				return err
			}
			v.PayRangeSub = append(v.PayRangeSub, min, max)
		}
		if v.WithdrawRange != "" {
			var min, max float64
			count, err := fmt.Sscanf(v.WithdrawRange, "%f-%f", &min, &max)
			if err != nil {
				log.Error("err:%v", err)
				return err
			}
			if count != 2 {
				err = errors.New("invalid count")
				log.Error("err:%v", err)
				return err
			}
			if min > max {
				err = errors.New("invalid min max")
				log.Error("err:%v", err)
				return err
			}
			v.WithdrawRangeSub = append(v.WithdrawRangeSub, min, max)
		}
	}
	ConfigMerchats = one
	return nil
}

// GetConfigMerchatsByID 获取商户
func GetConfigMerchatsByID(id int) *common.Merchats {
	for _, v := range ConfigMerchats {
		if v.MID == id {
			return v
		}
	}
	return nil
}

// GetConfigMerchatsByID 获取商户
func GetConfigMerchatsByIDReal(id int) *common.Merchats {
	mer := &common.Merchats{MID: id}
	db.Mysql().Get(mer)
	if mer.PayChannels != "" {
		// 左右增加[]
		str := "[" + mer.PayChannels + "]"
		err := json.Unmarshal([]byte(str), &mer.PayChannelsSub)
		if err != nil {
			log.Error("err:%v", err)
			return nil
		}
	}
	if mer.PayoutPayChannels != "" {
		// 左右增加[]
		str := "[" + mer.PayoutPayChannels + "]"
		err := json.Unmarshal([]byte(str), &mer.PayoutChannelsSub)
		if err != nil {
			log.Error("err:%v", err)
			return nil
		}
	}
	if mer.WhiteIP != "" {
		ips := strings.Split(mer.WhiteIP, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip) // 去除空格
			normalizedIP := util.NormalizeIP(ip)
			if normalizedIP == "" {
				log.Error("invalid IP: %s\n", ip)
				return nil
			}
			mer.WhiteIPSub = append(mer.WhiteIPSub, normalizedIP)
		}
	}
	if mer.PayRange != "" {
		var min, max float64
		count, err := fmt.Sscanf(mer.PayRange, "%f-%f", &min, &max)
		if err != nil {
			log.Error("err:%v", err)
			return nil
		}
		if count != 2 {
			err = errors.New("invalid count")
			log.Error("err:%v", err)
			return nil
		}
		if min > max {
			err = errors.New("invalid min max")
			log.Error("err:%v", err)
			return nil
		}
		mer.PayRangeSub = append(mer.PayRangeSub, min, max)
	}
	if mer.WithdrawRange != "" {
		var min, max float64
		count, err := fmt.Sscanf(mer.WithdrawRange, "%f-%f", &min, &max)
		if err != nil {
			log.Error("err:%v", err)
			return nil
		}
		if count != 2 {
			err = errors.New("invalid count")
			log.Error("err:%v", err)
			return nil
		}
		if min > max {
			err = errors.New("invalid min max")
			log.Error("err:%v", err)
			return nil
		}
		mer.WithdrawRangeSub = append(mer.WithdrawRangeSub, min, max)
	}
	return mer
}
