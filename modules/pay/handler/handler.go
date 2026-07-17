package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/call"
	"server/common"
	"server/config"
	"server/db"
	"server/modules/pay/allpay"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/util"
	"time"

	"github.com/liangdas/mqant/log"
)

func Recharge(req *common.PayReq) (ret common.PayResp) {
	trys := []int{}
	tryCount := 10
	var s string

	// 过滤费率高于商户的通道
	for _, v := range call.ConfigPayChannels {
		if v.PayFee > req.Merchat.PayFee {
			trys = append(trys, v.ChannelID)
		}
		s += fmt.Sprintf("pay_channel_%d, ", v.ChannelID)
	}

	log.Info("订单%s,可选总渠道:%s,不可用渠道:%v", req.OrderID, s, trys)
	for i := 0; i < tryCount; i++ {
		ret, trys = TryPay(req, trys)
		if ret.Err == nil || ret.Err.Error() == "no payways" {
			return
		}
	}
	return
}

// TryPay 尝试支付函数，根据配置和随机选择支付渠道进行支付
// 参数:
//   - req: 支付请求结构体
//   - trys: 可用支付渠道列表
//
// 返回值:
//   - ret: 支付响应结构体
//   - try: 已尝试的支付渠道列表
func TryPay(req *common.PayReq, trys []int) (ret common.PayResp, try []int) {
	// 初始化支付方式为无效值
	var payWay values.PayWay = -1

	// 存储可用支付渠道的切片
	payChannels := []*common.ConfigPayChannels{}

	// 根据商户ID获取商户配置信息
	m := call.GetConfigMerchatsByID(req.MchId)

	// 如果商户配置不存在，返回错误
	if m == nil {
		ret.Err = errors.New("invalid mid")
		return
	}

	// 如果商户配置了子支付渠道，则获取对应的支付渠道配置
	if len(m.PayChannelsSub) > 0 {
		payChannels = call.GetConfigPayChannelByCids(m.PayChannelsSub, req.Amount, trys...)

		// 如果存在可用支付渠道，则从其中随机选择一个
		if len(payChannels) > 0 { // 只能从可用渠道中选
			index := util.RandomOneEleFromSlice(payChannels, "Weight")
			if index >= 0 {
				payWay = values.PayWay(payChannels[index].ChannelID)
			}
		}
	} else {
		// 如果没有配置可用的支付渠道，则使用原来的逻辑选择支付方式
		if payWay < 0 {
			payWay = values.ChoosePayWay(req.UID, req.Channel, req.Amount, trys)
		}
	}

	// 记录调试日志，输出请求信息和选择的支付方式
	log.Debug("req:%+v,payway:%v", req, payWay)

	// 如果仍然没有选择到有效的支付方式，返回错误
	if payWay < 0 {
		ret.Err = errors.New("no payways")
		return
	}

	// 如果配置了选择支付方式，则将支付方式记录到Redis中
	if config.GetConfig().Pay.SelectPayWay {
		call.Go(func() {
			// 构造用户支付记录的Redis键
			key := common.GetRedisKeyPlayerPay(req.UID)

			// 将支付方式添加到Redis有序集合中
			db.Redis().ZAdd(key, 0, int(payWay))

			// 设置Redis键的过期时间为10分钟
			db.Redis().Expire(key, 10*time.Minute) // 保留10分钟
		})
	}

	// 更新请求中的支付渠道为选中的支付方式
	req.Channel = int(payWay)

	// 创建充值基础对象
	base := base.NewRechargeBase(req)

	// 创建子支付对象
	allpay.NewSub(base, int(payWay))

	// 如果子支付对象创建成功，则执行支付请求
	if base.Sub != nil {
		ret.Err = base.Req()
	} else {
		// 如果子支付对象创建失败，返回内部错误
		ret.Err = errors.New("inner error")
	}

	// 根据支付结果处理成功或失败的情况
	if ret.Err != nil {
		// 支付失败时，记录支付失败日志
		values.PayFail(payWay)
	} else {
		// 支付成功时，将支付响应赋值给返回结果
		if base.PayResp != nil {
			ret = *base.PayResp
		} else {
			// 增加一个 else 分支来处理异常情况
			// 这表明支付流程逻辑上成功了，但没有得到有效的响应数据
			// 你应该记录这个异常，并可能需要将 ret.Err 设置为一个明确的错误
			log.Error("错误：代收发送成功但 base.PayResp 为 nil。支付方式：%s", payWay)
			ret.Err = fmt.Errorf("payment response is empty, cannot process result")
		}
	}

	// 将已尝试的支付方式添加到尝试列表中
	try = append(try, int(payWay))
	ret.UseOriginURL = base.UseOriginURL
	// 返回支付结果和尝试列表
	return
}

// TryWithdraw 尝试执行提现操作
// 该函数负责处理订单提现逻辑，包括支付账户解析、渠道选择、请求发送等步骤
func TryWithdraw(or *common.OrderWithdraw) {
	// 初始化支付数据结构
	send := common.CommonPayData{}
	originStatus := or.Status

	// 尝试解析支付账户信息，如果解析失败则记录错误并回调
	if err := json.Unmarshal(or.PayAccount, &send); err != nil {
		log.Error("withdraw unmarshal err %v", err)
		or.Status = common.StatusROrderFail
		or.Message = "Invalid pay account."
		call.WithdrawCallback(or, originStatus)
		return
	}

	// 构造提现请求参数
	req := &common.WithdrawReq{
		OrderID: or.OrderID,
		Amount:  or.Amount,
		UID:     or.UID,
	}

	// 将支付数据赋值给请求结构体
	req.CommonPayData = send

	// 选择合适的提现渠道，如果无可用渠道则返回失败
	log.Info("开始检查提现渠道")
	cid := values.ChooseWithdrawWay(req, []int{}, or.MID)
	if cid < 0 {
		or.Status = common.StatusROrderFail
		or.Message = "No payment channels available."
		call.WithdrawCallback(or, originStatus)
		return
	}

	// 更新订单中的支付渠道信息
	db.Mysql().Update(&common.OrderWithdraw{ID: or.ID}, map[string]interface{}{"pay_channel": int(cid)})
	or.PayChannel = cid
	req.Channel = cid

	// 初始化提现基础处理对象
	base := base.NewWithdrawBase(req)

	// 创建子支付处理器
	allpay.NewSub(base, int(cid))

	// 如果子处理器创建失败，返回错误并回调
	if base.Sub == nil {
		or.Status = common.StatusROrderFail
		or.Message = "No payment channels available."
		call.WithdrawCallback(or, originStatus)
		return
	}

	// 发送提现请求
	err := base.Req()

	// 如果请求失败，返回错误并回调
	if err != nil {
		backReq := &values.BackReq{
			Opt:      int(base.Opt),
			Channel:  int(base.Channel),
			ErrorMsg: string(base.Body),
		}
		util.Go(func() {
			util.HttpPostWithoutRet(config.GetConfig().Pay.BackSysAddr+"/errNotify", backReq, nil)
		})
		or.Status = common.StatusROrderFail
		or.Message = "No payment channels available."
		call.WithdrawCallback(or, originStatus)
		return
	}

	// 更新订单中的API订单ID
	u := map[string]interface{}{"api_order_id": base.WithdrawResp.APIOrderID}
	db.Mysql().Update(&common.OrderWithdraw{ID: or.ID}, u)
}
