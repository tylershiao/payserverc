package call

import (
	"errors"
	"fmt"
	"math/rand"
	"server/common"
	"server/db"
	"server/sign"
	"server/util"
	"strconv"
	"strings"
	"time"

	"github.com/liangdas/mqant/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	NotifyInerval = []int{30, 120, 270, 480} // 秒
)

func PayCallback(order *common.OrderPay, success bool) error {
	_, err := payCallback(order, success, nil)
	return err
}

// PayCallbackFromUpstream applies the 5/1000 drop rate to successful upstream callbacks.
// It reports whether normal callback processing completed.
func PayCallbackFromUpstream(order *common.OrderPay, success bool) (bool, error) {
	return payCallback(order, success, rand.Intn)
}

func payCallback(order *common.OrderPay, success bool, draw func(int) int) (bool, error) {
	log.Debug("PayCallback:%+v,success:%v", order, success)
	update := map[string]interface{}{}
	status := common.StatusROrderCreate
	if !success {
		status = common.StatusROrderFail
	} else {
		status = common.StatusROrderFinish
	}
	update["status"] = status
	update["callback_time"] = time.Now().Unix()
	update["api_order_id"] = order.APIOrderID

	var settledOrder *common.OrderPay
	err := db.Mysql().C().Transaction(func(tx *gorm.DB) error {
		current := &common.OrderPay{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(current, order.ID).Error; err != nil {
			return err
		}
		// Recheck under the row lock so concurrent callbacks and queries cannot settle a dropped order.
		if success && current.IsDropped {
			return nil
		}
		// An earlier callback may have read a pending order before a successful payment committed.
		if current.Status == status || (!success && current.Status != order.Status) {
			return nil
		}
		if success && draw != nil && draw(1000) < common.OrderPayDropRate {
			// UpdateColumn deliberately leaves status, callback_time and updated_at untouched.
			return tx.Model(&common.OrderPay{}).Where("id = ?", current.ID).UpdateColumn("is_dropped", true).Error
		}
		// 更新订单状态
		res := tx.Model(&common.OrderPay{}).Where("id = ? and status <> ?", order.ID, status).Updates(update)
		if res.Error != nil {
			log.Error("err:%v", res.Error)
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		mer := &common.Merchats{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("mid = ?", current.MID).First(mer).Error
		if err != nil {
			log.Error("err:%v", err)
			return err
		}
		col := "balance_settle"
		balTyp := common.FundDetailsBalanceTypeSettle
		bal := mer.BalanceSettle
		if current.RealTimeSettle { // 实时结算的通道直接增加余额
			col = "balance"
			balTyp = common.FundDetailsBalanceTypeAvailable
			bal = mer.Balance
		}
		var addAmount float64
		// 更新商户余额
		if success {
			addAmount = current.Amount - current.Fee
		} else if current.Status == common.StatusROrderFinish { // 原本成功的单转失败了,要扣除余额
			addAmount = -(current.Amount - current.Fee)
		}
		if addAmount != 0 {
			res := tx.Model(&common.Merchats{}).Where("id = ?", mer.ID).Updates(map[string]interface{}{col: gorm.Expr(fmt.Sprintf("%s + %.2f", col, addAmount))})
			if res.Error != nil {
				log.Error("err:%v", res.Error)
				return res.Error
			}
			if res.RowsAffected == 0 {
				err = errors.New("update fail")
				return err
			}
			// 写入资金变动明细
			fund := &common.FundDetails{
				MID:             mer.MID,
				OrderID:         current.OrderID,
				Type:            common.FundDetailsTypeRecharge,
				TypeName:        common.GetFundDetailsTypeName(common.FundDetailsTypeRecharge),
				Amount:          addAmount,
				BalanceType:     balTyp,
				BalanceTypeName: common.GetFundDetailsBalanceTypeName(balTyp),
				Balance:         bal + addAmount,
			}
			if err := tx.Create(fund).Error; err != nil {
				log.Error("err:%v", err)
				return err
			}
		}

		current.Message = order.Message
		settledOrder = current
		return nil
	})
	if err != nil || settledOrder == nil {
		return false, err
	}

	// Start merchant notifications only after the order and accounting transaction commits.
	merchat := GetConfigMerchatsByID(settledOrder.MID)
	if merchat == nil {
		log.Error("invalid merchat:%d", settledOrder.MID)
		return true, nil
	}
	send := &common.PayCallbackReq{MerchantID: settledOrder.MID, PlatformOrderID: settledOrder.OrderID, MerchantOrderID: settledOrder.MOrderID,
		OrderAmount: decimal.NewFromFloat(settledOrder.Amount).Shift(2).IntPart(), OrderStatus: status, StatusMessage: settledOrder.Message, Metadata: settledOrder.Extra}
	send.Signature = sign.SignMD5(send, merchat.Key, "Sign", "signature")
	callbackURL := settledOrder.CallbackURL
	IndexTryWithParam(NotifyInerval, func() error {
		return util.HttpPostWithoutRet(callbackURL, send, nil)
	})
	return true, nil

	// var err error
	// tx := db.Mysql().Begin()
	// defer func() {
	// 	if err == nil {
	// 		if err := tx.Commit().Error; err != nil {
	// 			tx.Rollback()
	// 			return
	// 		}
	// 	} else {
	// 		log.Error("err:%v", err)
	// 		tx.Rollback()
	// 		return
	// 	}
	// }()
	// // 更新订单状态
	// rows, err := db.Mysql().UpdateWRes(&common.OrderPay{}, update, fmt.Sprintf("id = %d and status <> %d", order.ID, status), tx)
	// if rows == 0 {
	// 	err = errors.New("update fail")
	// 	return nil
	// }
	// if err != nil {
	// 	return err
	// }
	// col := "balance_settle"
	// if order.RealTimeSettle { // 实时结算的通道直接增加余额
	// 	col = "balance"
	// }
	// // 更新商户余额
	// if success {
	// 	rows, err = db.Mysql().UpdateRes(&common.Merchats{MID: order.MID}, map[string]interface{}{col: gorm.Expr(fmt.Sprintf("%s + %.2f", col, order.Amount-order.Fee))}, tx)
	// 	if rows == 0 {
	// 		err = errors.New("update fail")
	// 		return nil
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// // 回调商户
	// merchat := GetConfigMerchatsByID(order.MID)
	// if merchat == nil {
	// 	log.Error("invalid merchat:%d", order.MID)
	// 	return nil
	// }
	// send := &common.PayCallbackReq{MerchantID: order.MID, PlatformOrderID: order.OrderID, MerchantOrderID: order.MOrderID,
	// 	OrderAmount: int64(order.Amount * 100), OrderStatus: status, StatusMessage: order.Message, Metadata: order.Extra}
	// send.Signature = sign.SignMD5(send, merchat.Key, "Sign", "signature")
	// IndexTryWithParam(NotifyInerval, func() error {
	// 	return util.HttpPostWithoutRet(order.CallbackURL, send, nil)
	// })

	// return nil
}

func WithdrawCallback(order *common.OrderWithdraw, originStatus int) error {
	log.Debug("WithdrawCallback:%+v", order)
	success := order.Status == common.StatusROrderFinish
	// 修改订单状态,给商户回调订单状态,修改商户余额
	wsql := ""
	updates := map[string]interface{}{"status": order.Status, "callback_time": time.Now().Unix()}
	if order.Message != "" {
		msg := order.Message
		if len(msg) > 256 {
			msg = msg[:256]
		}
		updates["message"] = msg
	}
	if order.Utr != "" {
		utr := order.Utr
		if len(utr) > 256 {
			utr = utr[:256]
		}
		updates["utr"] = utr
	}
	if success {
		wsql = fmt.Sprintf("id = %d and status < %d", order.ID, common.StatusROrderFinish)
	} else {
		wsql = fmt.Sprintf("id = %d and status <= %d", order.ID, common.StatusROrderFinish)
	}

	updated := false
	err := db.Mysql().C().Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&common.OrderWithdraw{}).Where(wsql).Updates(updates)
		if res.RowsAffected == 0 {
			// err := errors.New("update fail")
			return nil
		}
		if res.Error != nil {
			log.Error("Withdraw callback err:%v", res.Error)
			return res.Error
		}
		var addAmount float64
		// 2.返还余额
		if !success {
			addAmount = order.Fee + order.Amount
		} else if originStatus > common.StatusROrderFinish { // 失败转成功,扣除余额
			addAmount = -(order.Fee + order.Amount)
		}
		if addAmount != 0 {
			mer := &common.Merchats{}
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("mid = ?", order.MID).First(mer).Error
			if err != nil {
				log.Error("err:%v", err)
				return err
			}
			res = tx.Model(&common.Merchats{}).Where("mid = ?", order.MID).
				Updates(map[string]interface{}{"balance": gorm.Expr("balance + ?", addAmount)})
			if res.RowsAffected == 0 {
				err := errors.New("update fail")
				return err
			}
			if res.Error != nil {
				log.Error("Withdraw callback err:%v", res.Error)
				return res.Error
			}
			// 写入资金变动明细
			fund := &common.FundDetails{
				MID:             order.MID,
				OrderID:         order.OrderID,
				Type:            common.FundDetailsTypeWithdraw,
				TypeName:        common.GetFundDetailsTypeName(common.FundDetailsTypeWithdraw),
				Amount:          addAmount,
				BalanceType:     common.FundDetailsBalanceTypeAvailable,
				BalanceTypeName: common.GetFundDetailsBalanceTypeName(common.FundDetailsBalanceTypeAvailable),
				Balance:         mer.Balance + addAmount,
			}
			if err := tx.Create(fund).Error; err != nil {
				log.Error("err:%v", err)
				return err
			}
		}

		updated = true

		// 3.回调商户
		merchat := GetConfigMerchatsByID(order.MID)
		if merchat == nil {
			log.Error("invalid merchat:%d", order.MID)
			return nil
		}
		send := &common.WithdrawCallbackReq{MerchantID: order.MID, PlatformOrderID: order.OrderID, MerchantOrderID: order.MOrderID,
			OrderAmount: decimal.NewFromFloat(order.Amount).Shift(2).IntPart(), OrderStatus: order.Status, StatusMessage: order.Message, Metadata: order.Extra, TransactionReference: order.Utr}
		send.Signature = sign.SignMD5(send, merchat.Key, "Sign", "signature")
		IndexTryWithParam(NotifyInerval, func() error {
			return util.HttpPostWithoutRet(order.CallbackURL, send, nil)
		})
		return nil
	})
	if err != nil {
		return err
	}
	if !success && updated {
		rollbackWithdrawCardRiskCountOnFailure(order)
	}
	return nil
	// // 1.修改订单状态
	// tx := db.Mysql().Begin()
	// var err error
	// defer func() {
	// 	if err == nil {
	// 		if err := tx.Commit().Error; err != nil {
	// 			tx.Rollback()
	// 			return
	// 		}
	// 	} else {
	// 		log.Error("err:%v", err)
	// 		tx.Rollback()
	// 		return
	// 	}
	// }()
	// res := tx.Model(&common.OrderWithdraw{}).Where(wsql).Updates(updates)
	// if res.Error != nil || res.RowsAffected == 0 {
	// 	log.Error("Withdraw callback err:%v,rows:%d", res.Error, res.RowsAffected)
	// 	tx.Rollback()
	// 	err = res.Error
	// 	return err
	// }

	// // 2.返还余额
	// if !success {
	// 	res = tx.Model(&common.Merchats{}).Where("mid = ?", order.MID).Updates(map[string]interface{}{"balance": gorm.Expr("balance + ?", order.Fee+order.Amount)})
	// 	if res.Error != nil || res.RowsAffected == 0 {
	// 		log.Error("Withdraw callback err:%v,rows:%d", res.Error, res.RowsAffected)
	// 		tx.Rollback()
	// 		err = res.Error
	// 		return err
	// 	}
	// }

	// // 3.回调商户
	// merchat := GetConfigMerchatsByID(order.MID)
	// if merchat == nil {
	// 	log.Error("invalid merchat:%d", order.MID)
	// 	tx.Rollback()
	// 	err = errors.New("invalid merchat")
	// 	return err
	// }
	// send := &common.WithdrawCallbackReq{MerchantID: order.MID, PlatformOrderID: order.OrderID, MerchantOrderID: order.MOrderID,
	// 	OrderAmount: int64(order.Amount * 100), OrderStatus: order.Status, StatusMessage: order.Message, Metadata: order.Extra, TransactionReference: order.Utr}
	// send.Signature = sign.SignMD5(send, merchat.Key, "Sign", "signature")
	// IndexTryWithParam(NotifyInerval, func() error {
	// 	return util.HttpPostWithoutRet(order.CallbackURL, send, nil)
	// })

	return nil
}

/**
 * NewOrderID 根据传入的 mid 生成一个唯一的订单ID字符串。
 * 生成规则如下：
 * 1. 将 mid 转换为字符串，逐位映射为 StringMap 中对应的字符；
 * 2. 生成一个5位随机数（范围10000-99999），逐位映射为 StringMap 中对应的字符；
 * 3. 获取当前时间戳（毫秒级），逐位映射为 StringMap 中对应的字符；
 * 4. 如果配置中的 ServerFlag 不为 "a"，则将 ServerFlag 追加到结果末尾；
 * 5. 最终返回拼接后的字符串作为订单ID。
 *
 * @param mid 用于生成订单ID的中间标识符（通常为用户ID或商户ID）
 * @return 生成的唯一订单ID字符串
 */
func NewOrderID(mid int) string {
	var orderID strings.Builder
	str := strconv.Itoa(mid)
	for _, v := range str {
		orderID.WriteString(StringMap[string(v)])
	}
	ran := rand.Intn(90000) + 10000
	ranStr := strconv.Itoa(ran)
	for _, v := range ranStr {
		orderID.WriteString(StringMap[string(v)])
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	for _, v := range timestamp {
		orderID.WriteString(StringMap[string(v)])
	}
	return orderID.String()
}

var StringMap = map[string]string{
	"1": "T",
	"2": "L",
	"3": "P",
	"4": "B",
	"5": "C",
	"6": "X",
	"7": "Z",
	"8": "T",
	"9": "Y",
	"0": "W",
}

func NewOrderIDV2(typ string) string {
	return typ + SnowNode().Generate().String()
}
