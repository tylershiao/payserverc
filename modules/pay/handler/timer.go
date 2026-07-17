package handler

import (
	"fmt"
	"server/call"
	"server/common"
	"server/config"
	"server/db"
	"server/modules/pay/allpay"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"sort"
	"sync"
	"time"

	"github.com/liangdas/mqant/log"
)

var (
	WithdrawGroup sync.WaitGroup
)

func InitTimer() {
	// StartPayCheckTimer()
	// StartWithdrawCheckTimer()
	WithdrawTimer() // 扫描order表进行代付提交
	QueryPayTimer() // 扫描order表的支付订单进行查询
	QueryWithdrawTimer()
}

func StartWithdrawCheckTimer() {
	t := config.GetConfig().Pay.WithdrawCheckTime
	if t <= 0 {
		return
	}
	time.AfterFunc(time.Duration(t)*time.Minute, func() {
		WithdrawCheck()
		StartWithdrawCheckTimer()
	})
}

// 每隔一段时间检测渠道成功率，切换渠道
func WithdrawCheck() {
	log.Debug("withdraw check start...... ")
	defer func() {
		log.Debug("withdraw check finish...... ")
		for _, v := range call.ConfigWithdrawChannels {
			log.Debug("%+v", *v)
		}
	}()
	successPer := map[int]int64{}
	successTotalPer := map[int]int64{}
	for i, v := range call.GetConfigWithdrawChannels() {
		per := getWithdrawPer(v.ChannelID, i)
		successPer[v.ChannelID] = per
		tPer := getWithdrawTotal(v.ChannelID)
		successTotalPer[v.ChannelID] = tPer
	}
	log.Debug("successPer:%v,successTotalPer:%v", successPer, successTotalPer)
	values.WithdrawWeightLock.Lock()
	defer values.WithdrawWeightLock.Unlock()
	// 全部按成功率重置权重
	for _, v := range call.ConfigWithdrawChannels {
		sup := successPer[v.ChannelID]
		sut := successTotalPer[v.ChannelID]
		if successPer[v.ChannelID] == -1 {
			successPer[v.ChannelID] = config.GetConfig().Pay.BaseWithdrawSuccess // 一单都没有按50的成功率计算
			sup = 0
		}
		if successTotalPer[v.ChannelID] == -1 {
			successTotalPer[v.ChannelID] = config.GetConfig().Pay.BaseWithdrawSuccess // 一单都没有按50的成功率计算
			sut = 0
		}
		per := successPer[v.ChannelID]*100 + successTotalPer[v.ChannelID] + 1 // 近期成功率*100+今日成功率：近期有成功的优先于今日成功的
		cid := v.ChannelID
		v.Weight = int(per)
		call.Go(func() {
			db.Mysql().UpdateW(&common.ConfigWithdrawChannels{}, map[string]interface{}{"weight": per, "success_today": sut, "success": sup, "fails": 0}, fmt.Sprintf("channel_id = %v", cid))
		})
	}
	sort.Slice(call.ConfigWithdrawChannels, func(i, j int) bool {
		return call.ConfigWithdrawChannels[i].Weight > call.ConfigWithdrawChannels[j].Weight
	})
}

func getWithdrawPer(channel, index int) int64 {
	now := time.Now()
	start := now.Add(-6 * time.Minute).Format("2006-01-02 15:04:05")
	end := now.Add(-1 * time.Minute).Format("2006-01-02 15:04:05")
	total := getWithdrawCount(start, end, channel, false)
	conf := &common.ConfigWithdrawChannels{ChannelID: channel}
	db.Mysql().Get(conf)
	total += conf.Fails
	if total == 0 {
		return -1
	}
	success := getWithdrawCount(start, end, channel, true)
	per := success * 100 / total
	log.Debug("channel:%v,total:%v,success:%v,per:%v", channel, total, success, per)
	return per
}

func getWithdrawTotal(channel int) int64 {
	now := time.Now()
	start := now.Format("2006-01-02")
	end := now.AddDate(0, 0, 1).Format("2006-01-02")
	total := getWithdrawCount(start, end, channel, false)
	conf := &common.ConfigWithdrawChannels{ChannelID: channel}
	db.Mysql().Get(conf)
	total += conf.Fails
	if total == 0 {
		return -1
	}
	success := getWithdrawCount(start, end, channel, true)
	per := success * 100 / total
	log.Debug("total:channel:%v,total:%v,success:%v,per:%v", channel, total, success, per)
	return per
}

func getWithdrawCount(start, end string, channel int, success bool) int64 {
	sql := fmt.Sprintf("created_at >= '%v' and created_at < '%v'", start, end)
	if channel >= 0 {
		sql += fmt.Sprintf(" and pay_channel = %v", channel)
	}
	if success {
		sql += fmt.Sprintf(" and status = %v", common.StatusROrderFinish)
	} else {
		sql += fmt.Sprintf(" and status > %d", common.StatusROrderCreate)
	}
	return db.Mysql().Count(&common.OrderWithdraw{}, sql)
}

func WithdrawTimer() {
	list := []*common.OrderWithdraw{}
	db.Mysql().QueryAll(fmt.Sprintf("status = %d", common.StatusROrderWaitting), "", &common.OrderWithdraw{}, &list)
	l := len(list)
	log.Debug("WithdrawTimer orders:%v", l)
	if l > 0 {
		WithdrawGroup.Add(l)
		for _, v := range list {
			// 提交订单
			res, err := db.Mysql().UpdateRes(&common.OrderWithdraw{OrderID: v.OrderID, Status: common.StatusROrderWaitting},
				map[string]interface{}{"status": common.StatusROrderPay})
			if err != nil {
				log.Error("err:%v", err)
				WithdrawGroup.Done()
				continue
			}
			if res == 0 {
				WithdrawGroup.Done()
				continue
			}
			or := v
			log.Debug("trying:%v", or.OrderID)
			call.Go(func() {
				TryWithdraw(or)
				WithdrawGroup.Done()
			})
		}
	}
	t := config.GetConfig().Pay.WithdrawScanTime
	if t <= 0 {
		t = 10
	}
	time.AfterFunc(time.Duration(t)*time.Second, WithdrawTimer)
}

func StartPayCheckTimer() {
	values.PayStatus, _ = db.Redis().GetInt(common.RedisKeyPayStatus)
	t := config.GetConfig().Pay.PayCheckTime
	if t <= 0 {
		return
	}
	time.AfterFunc(time.Duration(t)*time.Minute, func() {
		PayCheck()
		StartPayCheckTimer()
	})
}

var (
	totalAll int64 = 20 // 当前时间段计算的总单数
)

// 每隔一段时间检测渠道成功率，切换渠道
func PayCheck() {
	log.Debug("pay check start...... ")
	defer func() {
		log.Debug("pay check finish...... ")
		log.Debug("status:%v", values.PayStatus)
		for _, v := range call.ConfigPayChannels {
			log.Debug("%+v", *v)
		}
	}()
	log.Debug("limit:%v", config.GetConfig().Pay.CheckLimit)
	checkCount := config.GetConfig().Pay.PayCheckCount
	now := time.Now()
	start := now.Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
	end := now.Format("2006-01-02 15:04:05")
	total := getPayCount(start, end, -1, false)
	log.Debug("checkCount:%d,total:%d", checkCount, total)
	// 总订单过少
	if total < int64(checkCount) {
		return
	}

	success := getPayCount(start, end, -1, true)
	log.Debug("success:%d", success)
	// shouldKeep := success*100/total >= int64(config.GetConfig().Pay.CheckLimit)
	// 检查是否有渠道的成功率达标
	pays := []int{}
	values.PayWeightLock.RLock()
	for _, v := range call.ConfigPayChannels {
		if v.Weight > 0 {
			pays = append(pays, v.ChannelID)
		}
	}
	values.PayWeightLock.RUnlock()
	choosePay := -1
	var bestPer int64
	successPer := map[int]int64{}
	successTotalPer := map[int]int64{}
	// var totalSuccess int64
	for i, v := range pays {
		per := getSuccessPer(v, i)
		successPer[v] = per
		if per >= int64(config.GetConfig().Pay.CheckLimit) && per > bestPer {
			choosePay = v
			bestPer = per
		}
		tPer := getSuccessPerTotal(v)
		successTotalPer[v] = tPer
		// totalSuccess += tPer
	}
	log.Debug("successPer:%v,successTotalPer:%v", successPer, successTotalPer)
	// if shouldKeep {
	// 	for _, v := range pays {
	// 		st := successTotalPer[v]
	// 		sp := successPer[v]
	// 		cid := v
	// 		util.Go(func() {
	// 			db.Mysql().C().Model(&common.ConfigPayChannels{}).Where("channel_id = ?", cid).Updates(map[string]interface{}{"success_today": st, "success": sp})
	// 		})
	// 	}
	// 	return
	// }
	// 当前有达标的渠道选择，直接切换渠道权重
	values.PayWeightLock.Lock()
	defer values.PayWeightLock.Unlock()
	// 全部按成功率重置权重
	for _, v := range pays {
		st := successTotalPer[v]
		sp := successPer[v]
		var per int64
		if choosePay == -1 { // 所有通道成功率都很低时,自动选择近期成功率降低加入权重的占比,原因是其他通道会因此订单数过少导致难以尝试其他通道
			per = sp*10 + st + 1
		} else {
			per = sp*100 + st + 1
		}
		if st == 0 {
			per += int64(config.GetConfig().Pay.BaseSuccess)
		}
		for _, c := range call.ConfigPayChannels {
			if c.ChannelID == v {
				c.Weight = int(per)
				break
			}
		}
		cid := v
		call.Go(func() {
			db.Mysql().C().Model(&common.ConfigPayChannels{}).Where("channel_id = ?", cid).Updates(map[string]interface{}{"weight": per, "success_today": st, "success": sp})
		})
	}
	if values.PayStatus != values.PayStatusForceAuto {
		if choosePay > -1 {
			// 切换为主渠道模式
			values.PayStatus = values.PayStatusChannel
			call.Go(func() {
				db.Redis().SetData(common.RedisKeyPayStatus, values.PayStatusChannel)
			})
		} else {
			values.PayStatus = values.PayStatusAuto
			call.Go(func() {
				db.Redis().SetData(common.RedisKeyPayStatus, values.PayStatusAuto)
			})
		}
	}
	sort.Slice(call.ConfigPayChannels, func(i, j int) bool {
		return call.ConfigPayChannels[i].Weight > call.ConfigPayChannels[j].Weight
	})
}

func getPayCount(start, end string, channel int, success bool) int64 {
	sql := fmt.Sprintf("created_at >= '%v' and created_at < '%v'", start, end)
	if channel >= 0 {
		sql += fmt.Sprintf(" and pay_channel = %v", channel)
	}
	if success {
		sql += fmt.Sprintf(" and status = %v", common.StatusROrderFinish)
	}
	return db.Mysql().Count(&common.OrderPay{}, sql)
}

func getSuccessPer(channel, index int) int64 {
	now := time.Now()
	start := now.Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
	end := now.Format("2006-01-02 15:04:05")
	total := getPayCount(start, end, channel, false)
	success := getPayCount(start, end, channel, true)
	if index == 0 { // 当前最高权重的渠道
		if total < 20 { // 如果最高权重的渠道总单数少于20，降低最低单数判断标准
			totalAll = 10
		} else {
			totalAll = 20
		}
	}
	printTotal := total
	if total < totalAll {
		total = totalAll
	}
	per := success * 100 / total
	log.Debug("channel:%v,total:%v,totalAll:%v,success:%v,per:%v", channel, printTotal, totalAll, success, per)
	return per
}

func getSuccessPerTotal(channel int) int64 {
	now := time.Now()
	start := now.Format("2006-01-02")
	end := now.AddDate(0, 0, 1).Format("2006-01-02")
	total := getPayCount(start, end, channel, false)
	if total == 0 {
		return 0
	}
	success := getPayCount(start, end, channel, true)
	per := success * 100 / total
	log.Debug("total:channel:%v,total:%v,success:%v,per:%v", channel, total, success, per)
	return per
}

func QueryPayTimer() {
	t := config.GetConfig().Pay.QueryPayTime
	if t <= 0 {
		return
	}
	call.Go(func() {
		list := []*common.OrderPay{}
		now := time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
		db.Mysql().QueryAll(fmt.Sprintf("status = %d and created_at >= '%s'",
			common.StatusROrderCreate, now), "", &common.OrderPay{}, &list)
		l := len(list)
		log.Debug("QueryPayTimer orders:%v", l)
		if l > 0 {
			WithdrawGroup.Add(l)
			for i := range list {
				TryQueryPay(list[i])
			}
		}
		time.AfterFunc(time.Duration(t)*time.Second, QueryPayTimer)
	})
}

func TryQueryPay(or *common.OrderPay) {
	cid := or.PayChannel
	queryReq := &base.QueryReq{
		OrderID:    or.OrderID,
		APIOrderID: or.APIOrderID,
		Channel:    cid,
	}
	queryBase := base.NewQueryPayBase(queryReq)
	allpay.NewSub(queryBase, int(cid))
	if queryBase.Sub == nil {
		return
	}
	queryErr := queryBase.Req()
	if queryErr != nil {
		log.Error("orderid:%v,cid:%v,err:%v", or.OrderID, cid, queryErr)
		return
	}
	data := queryBase.QueryPayResp
	log.Debug("queryPay orderid:%v, resp:%+v", or.OrderID, data)
	// 支付订单只有成功的情况会去处理,其他情况不管
	if data.Status == base.QueryStatusOrderSuccess {
		call.PayCallback(or, true)
	}
}

func QueryWithdrawTimer() {
	t := config.GetConfig().Pay.QueryWithdrawTime
	if t <= 0 {
		return
	}
	call.Go(func() {
		list := []*common.OrderWithdraw{}
		db.Mysql().QueryAll(fmt.Sprintf("status = %d", common.StatusROrderPay), "", &common.OrderWithdraw{}, &list)
		l := len(list)
		log.Debug("QueryWithdrawTimer orders:%v", l)
		if l > 0 {
			for i := range list {
				TryQueryWithdraw(list[i])
			}
		}
		time.AfterFunc(time.Duration(t)*time.Second, QueryWithdrawTimer)
	})
}

func TryQueryWithdraw(or *common.OrderWithdraw) {
	cid := or.PayChannel
	queryReq := &base.QueryReq{
		OrderID:    or.OrderID,
		APIOrderID: or.APIOrderID,
	}
	queryBase := base.NewQueryWithdrawBase(queryReq)
	allpay.NewSub(queryBase, int(cid))
	if queryBase.Sub == nil || !queryBase.ShouldQueryWithdraw {
		return
	}
	queryErr := queryBase.Req()
	if queryErr != nil {
		log.Error("err:%v", queryErr)
		return
	}
	data := queryBase.QueryWithdrawResp
	log.Debug("queryWithdraw orderid:%v, resp:%+v", or.OrderID, data)
	u := map[string]interface{}{}
	if data.APIOrderID != "" {
		u["apipayid"] = data.APIOrderID
		or.APIOrderID = data.APIOrderID
	}
	originStatus := or.Status
	if data.Status == base.QueryStatusFail { // 提单失败,重新提交
		or.Status = common.StatusROrderFail
		call.WithdrawCallback(or, originStatus)
		return
	} else if data.Status == base.QueryStatusOrderSuccess { // 支付成功
		or.Status = common.StatusROrderFinish
		or.Message = data.Msg
		call.WithdrawCallback(or, originStatus)
		return
	}
}
