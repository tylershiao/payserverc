package common

import "time"

const (
	WithdrawCardRiskReviewPending = iota
	WithdrawCardRiskReviewApproved
	WithdrawCardRiskReviewRejected
)

// WithdrawCardRiskLog stores payout card risk trigger and review records.
type WithdrawCardRiskLog struct {
	ID           uint64     `gorm:"primarykey;autoIncrement" json:"ID"`
	CreatedAt    time.Time  `json:"CreatedAt"`
	UpdatedAt    time.Time  `json:"UpdatedAt"`
	RiskDate     string     `json:"RiskDate" gorm:"column:risk_date;type:varchar(8);not null;uniqueIndex:uk_risk_date_card,priority:1"`
	OrderID      string     `json:"OrderID" gorm:"column:order_id;type:varchar(128);not null;default:''"`
	MID          int        `json:"MID" gorm:"column:mid;not null;default:0"`
	MOrderID     string     `json:"MOrderID" gorm:"column:m_order_id;type:varchar(128);not null;default:''"`
	CardNo       string     `json:"CardNo" gorm:"column:card_no;type:varchar(128);not null;default:'';uniqueIndex:uk_risk_date_card,priority:2;index:idx_card_review,priority:1"`
	PayType      int        `json:"PayType" gorm:"column:pay_type;not null;default:0"`
	Amount       int        `json:"Amount" gorm:"column:amount;not null;default:0"`
	IP           string     `json:"IP" gorm:"column:ip;type:varchar(64);not null;default:''"`
	RedisCount   int64      `json:"RedisCount" gorm:"column:redis_count;not null;default:0"`
	RiskReason   string     `json:"RiskReason" gorm:"column:risk_reason;type:varchar(255);not null;default:''"`
	ReviewStatus int        `json:"ReviewStatus" gorm:"column:review_status;type:tinyint;not null;default:0;index:idx_card_review,priority:2"`
	ReviewedAt   *time.Time `json:"ReviewedAt" gorm:"column:reviewed_at"`
}

func (WithdrawCardRiskLog) TableName() string {
	return "withdraw_card_risk_log"
}
