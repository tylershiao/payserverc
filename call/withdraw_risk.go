package call

import (
	"encoding/json"
	"strings"
	"time"

	"server/common"
	"server/config"
	"server/db"

	"github.com/liangdas/mqant/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	withdrawCardDailyLimit     = int64(5)
	WithdrawCardRiskActionPass = iota
	WithdrawCardRiskActionPendingReview
	WithdrawCardRiskActionRejected
)

// CheckWithdrawCardRisk checks whether a payout card exceeds today's usage limit.
// If the card is approved today it passes, if rejected today it blocks, otherwise
// requests beyond the daily limit stay in create status waiting for review.
func CheckWithdrawCardRisk(mid int, orderID, mOrderID string, payType int, account, ip string, amount int) int {
	if !config.GetConfig().Pay.WithdrawRiskAudit {
		return WithdrawCardRiskActionPass
	}
	account = strings.TrimSpace(account)
	if account == "" {
		return WithdrawCardRiskActionPass
	}

	now := time.Now()
	day := now.Format("20060102")

	reviewStatus, ok := getWithdrawCardReviewStatus(day, account, now)
	if !ok {
		return WithdrawCardRiskActionPass
	}
	if reviewStatus == common.WithdrawCardRiskReviewApproved {
		return WithdrawCardRiskActionPass
	}
	if reviewStatus == common.WithdrawCardRiskReviewRejected {
		return WithdrawCardRiskActionRejected
	}

	count, ok := incrWithdrawCardDailyCount(day, account, now)
	if !ok {
		return WithdrawCardRiskActionPass
	}
	if count <= withdrawCardDailyLimit {
		return WithdrawCardRiskActionPass
	}

	insertWithdrawCardRiskLog(day, account, mid, orderID, mOrderID, payType, ip, amount, count)
	return WithdrawCardRiskActionPendingReview
}

func rollbackWithdrawCardRiskCountOnFailure(order *common.OrderWithdraw) {
	if !config.GetConfig().Pay.WithdrawRiskAudit {
		return
	}
	if order == nil || order.Status != common.StatusROrderFail {
		return
	}

	payData := &common.CommonPayData{}
	if err := json.Unmarshal(order.PayAccount, payData); err != nil {
		log.Error("withdraw risk rollback unmarshal pay account fail, order:%s err:%v", order.OrderID, err)
		return
	}
	account := strings.TrimSpace(payData.Account)
	if account == "" {
		return
	}

	riskDay := order.CreatedAt
	if riskDay.IsZero() {
		riskDay = time.Now()
	}
	day := riskDay.Format("20060102")
	reviewStatus, ok := getWithdrawCardReviewStatus(day, account, time.Now())
	if !ok || reviewStatus == common.WithdrawCardRiskReviewApproved {
		return
	}

	key := common.GetRedisKeyWithdrawRiskCardCount(day, account)
	current := db.Redis().GetInt64(key)
	if current <= 0 {
		return
	}
	next, err := db.Redis().Incr(key, -1)
	if err != nil {
		log.Error("withdraw risk rollback redis fail, key:%s err:%v", key, err)
		return
	}
	if next < 0 {
		_ = db.Redis().SetData(key, 0)
	}
}

func incrWithdrawCardDailyCount(day, account string, now time.Time) (int64, bool) {
	key := common.GetRedisKeyWithdrawRiskCardCount(day, account)
	count, err := db.Redis().Incr(key, 1)
	if err != nil {
		log.Error("withdraw risk incr redis fail, key:%s err:%v", key, err)
		return 0, false
	}
	_ = db.Redis().Expire(key, getTodayRemainTTL(now))
	return count, true
}

func getWithdrawCardReviewStatus(day, account string, now time.Time) (int, bool) {
	passKey := common.GetRedisKeyWithdrawRiskCardPass(day, account)
	if db.Redis().GetIntWithoutErr(passKey) == 1 {
		return common.WithdrawCardRiskReviewApproved, true
	}

	one := &common.WithdrawCardRiskLog{}
	err := db.Mysql().C().Where("risk_date = ? AND card_no = ?", day, account).Take(one).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.WithdrawCardRiskReviewPending, true
		}
		if isMissingWithdrawRiskTable(err) {
			log.Error("withdraw risk table missing, skip risk check day:%s account:%s err:%v", day, account, err)
			return common.WithdrawCardRiskReviewPending, false
		}
		log.Error("withdraw risk query review fail, day:%s account:%s err:%v", day, account, err)
		return common.WithdrawCardRiskReviewPending, true
	}
	if one.ReviewStatus == common.WithdrawCardRiskReviewApproved {
		_ = db.Redis().SetData(passKey, 1, getTodayRemainTTL(now))
		return common.WithdrawCardRiskReviewApproved, true
	}
	if one.ReviewStatus == common.WithdrawCardRiskReviewRejected {
		return common.WithdrawCardRiskReviewRejected, true
	}
	return common.WithdrawCardRiskReviewPending, true
}

func insertWithdrawCardRiskLog(day, account string, mid int, orderID, mOrderID string, payType int, ip string, amount int, count int64) {
	one := &common.WithdrawCardRiskLog{
		RiskDate:     day,
		OrderID:      orderID,
		MID:          mid,
		MOrderID:     mOrderID,
		CardNo:       account,
		PayType:      payType,
		Amount:       amount,
		IP:           ip,
		RedisCount:   count,
		RiskReason:   "当日代付次数超过限制",
		ReviewStatus: common.WithdrawCardRiskReviewPending,
	}
	err := db.Mysql().C().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "risk_date"}, {Name: "card_no"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"updated_at":  time.Now(),
			"order_id":    orderID,
			"mid":         mid,
			"m_order_id":  mOrderID,
			"pay_type":    payType,
			"amount":      amount,
			"ip":          ip,
			"redis_count": count,
			"risk_reason": one.RiskReason,
		}),
	}).Create(one).Error
	if err != nil {
		log.Error("withdraw risk insert mysql fail, day:%s account:%s err:%v", day, account, err)
	}
}

func isMissingWithdrawRiskTable(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "doesn't exist") || strings.Contains(msg, "does not exist") || strings.Contains(msg, "unknown table")
}

func getTodayRemainTTL(now time.Time) time.Duration {
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	ttl := end.Sub(now) + time.Second
	if ttl <= 0 {
		return 24 * time.Hour
	}
	return ttl
}
