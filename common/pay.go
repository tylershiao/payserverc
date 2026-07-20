package common

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	StatusROrderWaitting = iota // 等待提交状态
	StatusROrderCreate          // 未支付
	StatusROrderPay             // 代付:支付中
	StatusROrderFinish          // 完成支付
	StatusROrderFail            // 失败
	StatusROrderReturn          // 退回
	StatusROrderPayBack         // 客诉
	StatusROrderAll
)

var (
	PayInPrefix  = "In"
	PayOutPrefix = "Out"
)

// 支付类型
const (
	WithdrawTypeInvalid = iota
	WithdrawTypeBank
	WithdrawTypeUPI
	WithdrawTypeAll
)

type PayReq struct {
	MchId    int
	OrderID  string  // 平台订单号
	Name     string  // 用户姓名
	Mobile   string  // 手机号
	Email    string  // 邮箱
	Amount   float64 // 金额（单位：最小币种，如分）
	UID      string
	Channel  int
	IP       string
	DeviceID string
	Merchat  *Merchats
}

type PayResp struct {
	APIOrderID   string
	URL          string
	Channel      int
	Err          error
	Message      string
	UseOriginURL bool
}

type PayCallbackReq struct {
	MerchantID      int    `json:"merchantId"`
	PlatformOrderID string `json:"platformOrderId"` // 平台订单号
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单号
	OrderAmount     int64  `json:"orderAmount"`     // 金额（单位：最小币种，如分）
	OrderStatus     int    `json:"orderStatus"`     // 订单状态
	StatusMessage   string `json:"statusMessage"`   // 错误信息
	Metadata        string `json:"metadata"`
	Signature       string `json:"signature"` // 签名
}

type WithdrawReq struct {
	OrderID string  // 平台订单号
	Amount  float64 // 金额（单位：最小币种，如分）
	UID     string
	Address string
	Channel int
	CommonPayData
}

type CommonPayData struct {
	Name     string `json:"Name,omitempty"`     // 收款人姓名
	Mobile   string `json:"Mobile,omitempty"`   // 收款人手机号码
	Email    string `json:"Email,omitempty"`    // 收款人邮箱
	PayType  int    `json:"PayType,omitempty"`  // 支付方式 1Bank 2Upi
	IFSC     string `json:"IFSC,omitempty"`     // IFSC 码
	Account  string `json:"Account,omitempty"`  // 收款人卡号/账号
	BankName string `json:"BankName,omitempty"` // 银行名称
	IP       string `json:"IP,omitempty"`       // ip
	DeviceID string `json:"DeviceID,omitempty"` // 设备号
}

type WithdrawResp struct {
	APIOrderID string
	Channel    int
	Err        error
	Message    string
}

type WithdrawCallbackReq struct {
	MerchantID           int    `json:"merchantId"`
	PlatformOrderID      string `json:"platformOrderId"` // 平台订单号
	MerchantOrderID      string `json:"merchantOrderId"` // 商户订单号
	OrderAmount          int64  `json:"orderAmount"`     // 金额（单位：最小币种，如分）
	OrderStatus          int    `json:"orderStatus"`     // 订单状态
	StatusMessage        string `json:"statusMessage"`   // 错误信息
	Metadata             string `json:"metadata"`
	TransactionReference string `json:"transactionReference"` // 打款utr
	Signature            string `json:"signature"`            // 签名
}

type OrderPay struct {
	ID                uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt         time.Time      // 创建时间
	UpdatedAt         time.Time      // 更新时间
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`                                                                                     // 删除时间
	OrderID           string         `json:"OrderID" form:"OrderID" gorm:"uniqueIndex;comment:平台订单号;column:order_id;size:256;"`                  //平台订单号
	MOrderID          string         `json:"MOrderID" form:"MOrderID" gorm:"comment:商户订单号;column:m_order_id;size:256;"`                          //商户订单号
	APIOrderID        string         `json:"APIOrderID" form:"APIOrderID" gorm:"comment:三方订单号;column:api_order_id;size:256;"`                    //三方订单号
	MID               int            `json:"MID" form:"MID" gorm:"comment:商户id;column:mid;"`                                                     //商户id
	PayAccount        datatypes.JSON `json:"PayAccount" form:"PayAccount" gorm:"comment:账户信息;column:pay_account;size:512;" swaggertype:"object"` //账户信息
	Amount            float64        `json:"Amount" form:"Amount" gorm:"column:amount;"`                                                         //金额
	PayChannel        int            `json:"PayChannel" form:"PayChannel" gorm:"comment:支付通道;column:pay_channel;"`                               //支付通道
	Status            int            `json:"Status" form:"Status" gorm:"index;comment:状态;column:status;"`                                        //状态
	Message           string         `json:"Message" form:"Message" gorm:"comment:返回消息;column:message;"`                                         //返回消息
	UID               string         `json:"UID" form:"UID" gorm:"comment:玩家id;column:uid;"`                                                     //玩家id
	Extra             string         `json:"Extra" form:"Extra" gorm:"comment:透传参数;column:extra;size:256;"`                                      //透传参数
	ES                string         `json:"ES" form:"ES" gorm:"comment:额外字段;column:es;"`                                                        //额外字段
	CallbackURL       string         `json:"CallbackURL" form:"CallbackURL" gorm:"comment:回调地址;column:callback_url;"`                            //回调地址
	CallbackTime      int            `json:"CallbackTime" form:"CallbackTime" gorm:"comment:回调时间戳;column:callback_time;"`                        //回调时间戳
	Fee               float64        `json:"Fee" form:"Fee" gorm:"comment:手续费;column:fee;"`                                                      //手续费
	RealTimeSettle    bool           `json:"RealTimeSettle" form:"RealTimeSettle" gorm:"default:0;comment:实时结算;column:realtime_settle;"`         //实时结算
	QueryNotify       int64          `json:"query_notify" form:"query_notify" gorm:"column:query_notify"`
	QueryNotifyChatID int64          `json:"query_notify_chat_id" form:"query_notify_chat_id" gorm:"column:query_notify_chat_id"`
	PayUrl            string         `json:"pay_url" form:"pay_url" gorm:"column:pay_url"`
}

// TableName 代收订单 OrderPay自定义表名 order_pay
func (OrderPay) TableName() string {
	return "order_pay"
}

// 代付订单 结构体  OrderWithdraw
type OrderWithdraw struct {
	ID           uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt    time.Time      // 创建时间
	UpdatedAt    time.Time      // 更新时间
	OrderID      string         `json:"OrderID" form:"OrderID" gorm:"uniqueIndex;comment:平台订单号;column:order_id;size:256;"`                  //平台订单号
	MOrderID     string         `json:"MOrderID" form:"MOrderID" gorm:"comment:商户订单号;column:m_order_id;size:256;"`                          //商户订单号
	APIOrderID   string         `json:"APIOrderID" form:"APIOrderID" gorm:"comment:三方订单号;column:api_order_id;size:256;"`                    //三方订单号
	MID          int            `json:"MID" form:"MID" gorm:"comment:商户id;column:mid;"`                                                     //商户id
	PayAccount   datatypes.JSON `json:"PayAccount" form:"PayAccount" gorm:"comment:账户信息;column:pay_account;size:512;" swaggertype:"object"` //账户信息
	Amount       float64        `json:"Amount" form:"Amount" gorm:"column:amount;"`                                                         //金额
	PayChannel   int            `json:"PayChannel" form:"PayChannel" gorm:"index;comment:支付通道;column:pay_channel;"`                         //支付通道
	Status       int            `json:"Status" form:"Status" gorm:"index;comment:状态;column:status;"`                                        //状态
	Message      string         `json:"Message" form:"Message" gorm:"comment:返回消息;column:message;"`                                         //返回消息
	UID          string         `json:"UID" form:"UID" gorm:"comment:玩家id;column:uid;"`                                                     //玩家id
	Extra        string         `json:"Extra" form:"Extra" gorm:"comment:透传参数;column:extra;size:256;"`                                      //透传参数
	ES           string         `json:"ES" form:"ES" gorm:"comment:额外字段;column:es;"`                                                        //额外字段
	CallbackURL  string         `json:"CallbackURL" form:"CallbackURL" gorm:"comment:回调地址;column:callback_url;"`                            //回调地址
	CallbackTime int            `json:"CallbackTime" form:"CallbackTime" gorm:"comment:回调时间戳;column:callback_time;"`                        //回调时间戳
	Fee          float64        `json:"Fee" form:"Fee" gorm:"comment:手续费;column:fee;"`                                                      //手续费
	Utr          string         `json:"Utr" form:"Utr" gorm:"comment:打款utr;column:utr;size:256;"`                                           //打款utr
}

// TableName 代付订单 OrderWithdraw自定义表名 order_withdraw
func (OrderWithdraw) TableName() string {
	return "order_withdraw"
}
