package common

import "time"

// 重新加载配置的id
const (
	ReloadInvalid = iota
	ReloadPayChannel
	ReloadWithdrawChannel
	ReloadMerchats
	ReloadAll
)

type ConfigPayChannels struct {
	ID             uint      `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt      time.Time // 创建时间
	UpdatedAt      time.Time // 更新时间
	ChannelID      int       `json:"ChannelID" form:"ChannelID" gorm:"uniqueIndex;comment:通道id;column:channel_id;" binding:"required"`       //通道id
	Name           string    `json:"Name" form:"Name" gorm:"uniqueIndex;comment:通道名称;column:name;" binding:"required"`                       //通道名称
	PayFee         float64   `json:"PayFee" form:"PayFee" gorm:"default:8;comment:代收费率;column:pay_fee;" binding:"required"`                  //代收费率
	WithdrawFee    float64   `json:"WithdrawFee" form:"WithdrawFee" gorm:"default:3;comment:代付费率;column:withdraw_fee;" binding:"required"`   //代付费率
	WithdrawOne    float64   `json:"WithdrawOne" form:"WithdrawOne" gorm:"default:6;comment:代付单笔费率;column:withdraw_one;" binding:"required"` //代付单笔费率
	Weight         int       `json:"Weight" form:"Weight" gorm:"comment:权重;column:weight;" binding:"required"`                               //权重
	PayDown        float64   `json:"PayDown" form:"PayDown" gorm:"default:100;comment:支付下限;column:pay_down;" binding:"required"`             //支付下限
	PayUp          float64   `json:"PayUp" form:"PayUp" gorm:"default:10000;comment:支付上限;column:pay_up;" binding:"required"`                 //支付上限
	SuccessToday   int       `json:"SuccessToday" form:"SuccessToday" gorm:"comment:今日成功率;column:success_today;"`                            //今日成功率
	Success        int       `json:"Success" form:"Success" gorm:"comment:近5分钟成功率;column:success;"`                                          //近5分钟成功率
	RealTimeSettle bool      `json:"RealTimeSettle" form:"RealTimeSettle" gorm:"default:0;comment:实时结算;column:realtime_settle;"`             //实时结算
}

// TableName 支付通道 PayChannel自定义表名 pay_channel
func (ConfigPayChannels) TableName() string {
	return "pay_channel"
}

type ConfigWithdrawChannels struct {
	ID           uint      `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt    time.Time // 创建时间
	UpdatedAt    time.Time // 更新时间
	ChannelID    int       `json:"ChannelID" form:"ChannelID" gorm:"uniqueIndex;comment:通道id;column:channel_id;" binding:"required"`       //通道id
	Name         string    `json:"Name" form:"Name" gorm:"uniqueIndex;comment:通道名称;column:name;" binding:"required"`                       //通道名称
	PayFee       float64   `json:"PayFee" form:"PayFee" gorm:"default:8;comment:代收费率;column:pay_fee;" binding:"required"`                  //代收费率
	WithdrawFee  float64   `json:"WithdrawFee" form:"WithdrawFee" gorm:"default:3;comment:代付费率;column:withdraw_fee;" binding:"required"`   //代付费率
	WithdrawOne  float64   `json:"WithdrawOne" form:"WithdrawOne" gorm:"default:6;comment:代付单笔费率;column:withdraw_one;" binding:"required"` //代付单笔费率
	Weight       int       `json:"Weight" form:"Weight" gorm:"comment:权重;column:weight;" binding:"required"`                               //权重
	PayDown      float64   `json:"PayDown" form:"PayDown" gorm:"default:100;comment:支付下限;column:pay_down;" binding:"required"`             //支付下限
	PayUp        float64   `json:"PayUp" form:"PayUp" gorm:"default:10000;comment:支付上限;column:pay_up;" binding:"required"`                 //支付上限
	UPI          bool      `json:"UPI" form:"UPI" gorm:"default:0;comment:是否开放UPI代付;column:upi;"`                                          //是否开放UPI代付
	Bank         bool      `json:"Bank" form:"Bank" gorm:"default:0;comment:是否开放Bank代付;column:bank;"`                                      //是否开放Bank代付
	SuccessToday int       `json:"SuccessToday" form:"SuccessToday" gorm:"comment:今日成功率;column:success_today;"`                            //今日成功率
	Success      int       `json:"Success" form:"Success" gorm:"comment:近5分钟成功率;column:success;"`                                          //近5分钟成功率
	Fails        int64     `json:"Fails" form:"Fails" gorm:"comment:近5分钟提交直接失败的单数;column:fails;"`                                          //近5分钟提交直接失败的单数
}

// TableName 支付通道 WithdrawChannel自定义表名 pay_channel
func (ConfigWithdrawChannels) TableName() string {
	return "withdraw_channel"
}
