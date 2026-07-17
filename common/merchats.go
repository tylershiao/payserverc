package common

import "time"

type Merchats struct {
	ID                uint      `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt         time.Time // 创建时间
	UpdatedAt         time.Time // 更新时间
	MID               int       `json:"MID" form:"MID" gorm:"uniqueIndex;comment:商户id;column:mid;" binding:"required"`                //商户id
	Name              string    `json:"Name" form:"Name" gorm:"uniqueIndex;comment:商户名称;column:name;size:256;" binding:"required"`    //商户名称
	PayFee            float64   `json:"PayFee" form:"PayFee" gorm:"comment:代收费率;column:pay_fee;" binding:"required"`                  //代收费率
	WithdrawFee       float64   `json:"WithdrawFee" form:"WithdrawFee" gorm:"comment:代付费率;column:withdraw_fee;"`                      //代付费率
	WithdrawOne       float64   `json:"WithdrawOne" form:"WithdrawOne" gorm:"comment:代付单笔费率;column:withdraw_one;"`                    //代付单笔费率
	PayChannels       string    `json:"PayChannels" form:"PayChannels" gorm:"comment:可用通道;column:pay_channels;"`                      //可用通道
	PayoutPayChannels string    `json:"PayoutPayChannels" form:"PayoutPayChannels" gorm:"comment:代付可用通道;column:payout_pay_channels;"` // 代付可用通道
	PayChannelsSub    []int     `json:"-" form:"-" gorm:"-"`                                                                          //可用通道
	PayoutChannelsSub []int     `json:"-" form:"-" gorm:"-"`                                                                          //代付可用通道
	Key               string    `json:"Key" form:"Key" gorm:"comment:密钥;column:key;" binding:"required"`                              //密钥
	Balance           float64   `json:"Balance" form:"Balance" gorm:"default:0;comment:可用余额;column:balance;"`                         //可用余额
	BalanceSettle     float64   `json:"BalanceSettle" form:"BalanceSettle" gorm:"default:0;comment:待结算余额;column:balance_settle;"`     //待结算余额
	BalanceFreeze     float64   `json:"BalanceFreeze" form:"BalanceFreeze" gorm:"default:0;comment:冻结金额;column:balance_freeze;"`      //冻结金额
	PaySwitch         bool      `json:"PaySwitch" form:"PaySwitch" gorm:"comment:是否开启代收;column:pay_switch;"`                          //代收开关
	WithdrawSwitch    bool      `json:"WithdrawSwitch" form:"WithdrawSwitch" gorm:"comment:是否开启代付;column:withdraw_swith;"`            //代付开关
	Open              bool      `json:"Open" form:"Open" gorm:"default:0;comment:是否开启;column:open;"`                                  //是否开启
	WhiteIP           string    `json:"WhiteIP" form:"WhiteIP" gorm:"comment:ip白名单;column:white_ip;size:256;"`                        //ip白名单
	WhiteIPSub        []string  `json:"-" form:"-" gorm:"-"`                                                                          //ip白名单
	PayRange          string    `json:"PayRange" form:"PayRange" gorm:"comment:代收区间;column:pay_range;size:256;"`                      //代收区间
	PayRangeSub       []float64 `json:"-" form:"-" gorm:"-"`                                                                          //代收区间
	WithdrawRange     string    `json:"WithdrawRange" form:"WithdrawRange" gorm:"comment:代付区间;column:withdraw_range;size:256;"`       //代付区间
	WithdrawRangeSub  []float64 `json:"-" form:"-" gorm:"-"`                                                                          //代付区间
}

// TableName 商户 Merchats自定义表名 merchats
func (m *Merchats) TableName() string {
	return "merchats"
}

// 计算费率,1是代收 2是代付
func (m *Merchats) CallFee(opt int, amount float64) (fee float64) {
	if opt == 1 {
		fee = m.PayFee * amount / 100
	} else {
		fee = m.WithdrawFee*amount/100 + m.WithdrawOne
	}
	return
}
