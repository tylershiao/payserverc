package apiv1

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"server/call"
	"server/common"
	"server/db"
	"server/modules/pay/app"
	"server/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PayOutCreateReq struct {
	MerchantID         int    `json:"merchantId" binding:"required"`
	MerchantOrderID    string `json:"merchantOrderId" binding:"required"`
	OrderAmount        int    `json:"orderAmount" binding:"required"`
	ClientIP           string `json:"clientIp" binding:"required"`
	CustomerPhone      string `json:"customerPhone" binding:"required"`
	CustomerID         any    `json:"customerId"`
	CustomerName       string `json:"customerName"`
	CustomerEmail      string `json:"customerEmail"`
	DeviceFingerprint  string `json:"deviceFingerprint"`
	PayoutMethod       int    `json:"payoutMethod" binding:"required"` // 代付方式 1bank 2upi
	BeneficiaryAccount string `json:"beneficiaryAccount" binding:"required"`
	IFSCCode           string `json:"ifscCode"`        // PayType为1时必填
	BeneficiaryBank    string `json:"beneficiaryBank"` // PayType为1时必填
	CallbackURL        string `json:"callbackUrl" binding:"required"`
	Metadata           string `json:"metadata"`
	Signature          string `json:"signature" binding:"required"`
}

type PayOutCreateResp struct {
	MerchantID      int    `json:"merchantId"`      // 商户号
	PlatformOrderID string `json:"platformOrderId"` // 平台订单号
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单号
}

func PayOutCreate(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(PayOutCreateReq)
	if !a.S(req) {
		return
	}
	if !util.IsValidIndianLocalMobile(req.CustomerPhone) {
		a.Code = app.CodeParam
		a.Msg = "Invalid Mobile"
		return
	}
	if !util.IsValidIP(req.ClientIP) {
		a.Code = app.CodeParam
		a.Msg = "Invalid IP"
		return
	}
	if !util.IsValidHTTPURL(req.CallbackURL) {
		a.Code = app.CodeParam
		a.Msg = "Invalid NotifyURL"
		return
	}
	if req.PayoutMethod <= common.WithdrawTypeInvalid || req.PayoutMethod >= common.WithdrawTypeAll {
		a.Code = app.CodeParam
		a.Msg = "Invalid PayType"
		return
	}
	if req.PayoutMethod == common.WithdrawTypeBank {
		if req.BeneficiaryBank == "" {
			a.Code = app.CodeParam
			a.Msg = "BankName cannot be empty"
			return
		}
		if req.IFSCCode == "" {
			a.Code = app.CodeParam
			a.Msg = "IFSC cannot be empty"
			return
		}
	}

	mer := call.GetConfigMerchatsByIDReal(req.MerchantID)
	if mer != nil {
		a.Merchat = mer
	}
	if !a.CheckMerchat(req, req.MerchantID, req.Signature, 3) {
		return
	}
	if len(a.Merchat.WithdrawRangeSub) == 2 {
		min := int(a.Merchat.WithdrawRangeSub[0] * 100)
		max := int(a.Merchat.WithdrawRangeSub[1] * 100)
		if req.OrderAmount < min {
			a.Code = app.CodeParam
			a.Msg = fmt.Sprintf("The payment limit cannot be lower than %d", min)
			return
		}
		if max > 0 && req.OrderAmount > max {
			a.Code = app.CodeParam
			a.Msg = fmt.Sprintf("The payment limit cannot be higher than %d", max)
			return
		}
	}
	data := &common.CommonPayData{Name: req.CustomerName, Mobile: req.CustomerPhone, Email: req.CustomerEmail, IP: req.ClientIP, PayType: req.PayoutMethod,
		DeviceID: req.DeviceFingerprint, IFSC: req.IFSCCode, BankName: req.BeneficiaryBank, Account: req.BeneficiaryAccount}
	if !a.Risk(data) {
		return
	}
	one := &common.OrderWithdraw{MOrderID: req.MerchantOrderID, MID: req.MerchantID}
	db.Mysql().Get(one)
	if one.ID > 0 {
		a.Code = app.CodeDuplicateOrder
		return
	}
	log.Debug("PayOutCreate:%+v", req)
	amount := util.Decimal(float64(req.OrderAmount)/100, 2)
	fee := a.Merchat.CallFee(2, amount)
	total := amount + fee
	if a.Merchat.Balance < total {
		a.Code = app.CodeInsufficientBalance
		return
	}
	orderID := call.NewOrderIDV2(common.PayOutPrefix)
	riskAction := call.CheckWithdrawCardRisk(req.MerchantID, orderID, req.MerchantOrderID, req.PayoutMethod, req.BeneficiaryAccount, req.ClientIP, req.OrderAmount)
	if riskAction == call.WithdrawCardRiskActionRejected {
		a.Code = app.CodeRisk
		return
	}
	// 更新余额
	tx := db.Mysql().Begin()
	defer a.MCommit(tx)
	newMer := &common.Merchats{}
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("mid = ?", a.Merchat.MID).First(newMer).Error
	if err != nil {
		log.Error("err:%v", err)
		a.Code = app.CodeInnerError
		return
	}
	rows, err := db.Mysql().UpdateWRes(&common.Merchats{}, map[string]interface{}{"balance": gorm.Expr("balance - ?", total)},
		fmt.Sprintf("id = %d and balance >= %v", a.Merchat.ID, total), tx)
	if rows == 0 || err != nil {
		a.Code = app.CodeInsufficientBalance
		return
	}

	tmp, _ := json.Marshal(data)

	status := common.StatusROrderWaitting
	if req.MerchantID == 888 || riskAction == call.WithdrawCardRiskActionPendingReview { // 测试号或风控待审核
		status = common.StatusROrderCreate
	}
	order := &common.OrderWithdraw{
		OrderID:     orderID,
		MOrderID:    req.MerchantOrderID,
		MID:         req.MerchantID,
		Amount:      amount,
		PayAccount:  tmp,
		Status:      status,
		UID:         util.AnyToString(req.CustomerID),
		Extra:       req.Metadata,
		CallbackURL: req.CallbackURL,
		Fee:         fee,
	}
	err = tx.Create(order).Error
	if err != nil {
		a.Code = app.CodeInnerError
		return
	}
	// 写入资金变动明细
	fund := &common.FundDetails{
		MID:             mer.MID,
		OrderID:         orderID,
		Type:            common.FundDetailsTypeWithdraw,
		TypeName:        common.GetFundDetailsTypeName(common.FundDetailsTypeWithdraw),
		Amount:          -total,
		BalanceType:     common.FundDetailsBalanceTypeAvailable,
		BalanceTypeName: common.GetFundDetailsBalanceTypeName(common.FundDetailsBalanceTypeAvailable),
		Balance:         newMer.Balance - total,
	}
	if err := tx.Create(fund).Error; err != nil {
		log.Error("err:%v", err)
		a.Code = app.CodeInnerError
		return
	}

	a.Data = &PayOutCreateResp{
		MerchantID:      req.MerchantID,
		PlatformOrderID: orderID,
		MerchantOrderID: req.MerchantOrderID,
	}
	if req.MerchantID == 888 { // 测试号
		order.Status = common.StatusROrderFail
		if rand.Intn(2) == 0 {
			order.Status = common.StatusROrderFinish
		}
		time.AfterFunc(2*time.Second, func() {
			call.WithdrawCallback(order, status)
		})
		return
	}
}

type PayOutQueryReq struct {
	MerchantID      int    `json:"merchantId" binding:"required"`
	MerchantOrderID string `json:"merchantOrderId" binding:"required"`
	Signature       string `json:"signature" binding:"required"`
}

type PayOutQueryResp struct {
	MerchantID      int    `json:"merchantId"`      // 商户号
	PlatformOrderID string `json:"platformOrderId"` // 平台订单号
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单号
	OrderAmount     int64  `json:"orderAmount"`     // 订单金额
	OrderStatus     int    `json:"orderStatus"`     // 订单状态
}

func PayOutQuery(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(PayOutQueryReq)
	if !a.S(req) {
		return
	}
	mer := &common.Merchats{MID: req.MerchantID}
	db.Mysql().Get(mer)
	if mer.ID > 0 {
		a.Merchat = mer
	}
	if !a.CheckMerchat(req, req.MerchantID, req.Signature, 4) {
		return
	}
	order := &common.OrderWithdraw{MOrderID: req.MerchantOrderID, MID: req.MerchantID}
	err := db.Mysql().Get(order)
	if order.ID == 0 {
		a.Code = app.CodeOrderNotExist
		return
	}
	if err != nil {
		a.Code = app.CodeInnerError
		return
	}
	a.Data = &PayOutQueryResp{
		MerchantID:      req.MerchantID,
		PlatformOrderID: order.OrderID,
		MerchantOrderID: req.MerchantOrderID,
		OrderAmount:     decimal.NewFromFloat(order.Amount).Shift(2).IntPart(),
		OrderStatus:     order.Status,
	}
}
