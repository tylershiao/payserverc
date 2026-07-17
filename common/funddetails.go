package common

import (
	"time"

	"gorm.io/gorm"
)

// 资金明细 结构体  FundDetails
type FundDetails struct {
	ID              uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt       time.Time      // 创建时间
	UpdatedAt       time.Time      // 更新时间
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`                                                                         // 删除时间
	MID             int            `json:"MID" form:"MID" gorm:"index;comment:商户id;column:MID;"`                                   //商户id
	OrderID         string         `json:"OrderID" form:"OrderID" gorm:"comment:订单号;column:OrderID;"`                              //订单号
	Type            int            `json:"Type" form:"Type" gorm:"comment:变动类型;column:Type;"`                                      //变动类型
	TypeName        string         `json:"TypeName" form:"TypeName" gorm:"comment:变动类型;column:TypeName;"`                          //变动类型名称
	Amount          float64        `json:"Amount" form:"Amount" gorm:"comment:变化金额;column:Amount;"`                                //变化金额
	BalanceType     int            `json:"BalanceType" form:"BalanceType" gorm:"comment:变动余额类型;column:BalanceType;"`               //变动余额类型
	BalanceTypeName string         `json:"BalanceTypeName" form:"BalanceTypeName" gorm:"comment:变动余额类型名称;column:BalanceTypeName;"` //变动类型名称
	Balance         float64        `json:"Balance" form:"Balance" gorm:"type:decimal(18,2);comment:变动后余额;column:Balance;"`         //变动后余额
}

// 流水变动类型
const (
	FundDetailsTypeZero = iota
	FundDetailsTypeRecharge
	FundDetailsTypeWithdraw
	FundDetailsTypeU
	FundDetailsTypeSettle
	FundDetailsTypeFreeze
	FundDetailsTypeSubtract
	FundDetailsTypeAdd
	FundDetailsTypeAll
)

func GetFundDetailsTypeName(t int) string {
	switch t {
	case FundDetailsTypeRecharge:
		return "充值"
	case FundDetailsTypeWithdraw:
		return "提现"
	case FundDetailsTypeU:
		return "回U"
	case FundDetailsTypeSettle:
		return "结算"
	case FundDetailsTypeFreeze:
		return "冻结"
	case FundDetailsTypeSubtract:
		return "扣除"
	case FundDetailsTypeAdd:
		return "增加"
	default:
		return ""
	}
}

// 资金类型
const (
	FundDetailsBalanceTypeZero      = iota
	FundDetailsBalanceTypeAvailable // 可用
	FundDetailsBalanceTypeSettle    // 待结算
	FundDetailsBalanceTypeFreeze    // 冻结
	FundDetailsBalanceTypeAll
)

func GetFundDetailsBalanceTypeName(t int) string {
	switch t {
	case FundDetailsBalanceTypeAvailable:
		return "可用余额"
	case FundDetailsBalanceTypeSettle:
		return "待结算余额"
	case FundDetailsBalanceTypeFreeze:
		return "冻结金额"
	default:
		return ""
	}
}
