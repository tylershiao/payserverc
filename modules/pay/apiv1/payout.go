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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PayOutCreateReq struct {
	MID        int    `json:"MID" binding:"required"`
	MOrderID   string `json:"MOrderID" binding:"required"`
	Amount     int    `json:"Amount" binding:"required"`
	IP         string `json:"IP" binding:"required"`
	UserMobile string `json:"UserMobile" binding:"required"`
	UserID     any    `json:"UserID"`
	UserName   string `json:"UserName"`
	UserEmail  string `json:"UserEmail"`
	DeviceID   string `json:"DeviceID"`
	Type       int    `json:"Type" binding:"required"` // 代付方式 1bank 2upi
	Account    string `json:"Account" binding:"required"`
	IFSC       string `json:"IFSC"`     // PayType为1时必填
	BankName   string `json:"BankName"` // PayType为1时必填
	NotifyURL  string `json:"NotifyURL" binding:"required"`
	Extra      string `json:"Extra"`
	Sign       string `json:"Sign" binding:"required"`
}

type PayOutCreateResp struct {
	MID      int    // 商户号
	OrderID  string // 平台订单号
	MOrderID string // 商户订单号
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
	if !util.IsValidIndianLocalMobile(req.UserMobile) {
		a.Code = app.CodeParam
		a.Msg = "Invalid Mobile"
		return
	}
	if !util.IsValidIP(req.IP) {
		a.Code = app.CodeParam
		a.Msg = "Invalid IP"
		return
	}
	if !util.IsValidHTTPURL(req.NotifyURL) {
		a.Code = app.CodeParam
		a.Msg = "Invalid NotifyURL"
		return
	}
	if req.Type <= common.WithdrawTypeInvalid || req.Type >= common.WithdrawTypeAll {
		a.Code = app.CodeParam
		a.Msg = "Invalid PayType"
		return
	}
	if req.Type == common.WithdrawTypeBank {
		if req.BankName == "" {
			a.Code = app.CodeParam
			a.Msg = "BankName cannot be empty"
			return
		}
		if req.IFSC == "" {
			a.Code = app.CodeParam
			a.Msg = "IFSC cannot be empty"
			return
		}
	}

	mer := call.GetConfigMerchatsByIDReal(req.MID)
	if mer != nil {
		a.Merchat = mer
	}
	if !a.CheckMerchat(req, req.MID, req.Sign, 3) {
		return
	}
	if len(a.Merchat.WithdrawRangeSub) == 2 {
		min := int(a.Merchat.WithdrawRangeSub[0] * 100)
		max := int(a.Merchat.WithdrawRangeSub[1] * 100)
		if req.Amount < min {
			a.Code = app.CodeParam
			a.Msg = fmt.Sprintf("The payment limit cannot be lower than %d", min)
			return
		}
		if max > 0 && req.Amount > max {
			a.Code = app.CodeParam
			a.Msg = fmt.Sprintf("The payment limit cannot be higher than %d", max)
			return
		}
	}
	data := &common.CommonPayData{Name: req.UserName, Mobile: req.UserMobile, Email: req.UserEmail, IP: req.IP, PayType: req.Type,
		DeviceID: req.DeviceID, IFSC: req.IFSC, BankName: req.BankName, Account: req.Account}
	if !a.Risk(data) {
		return
	}
	one := &common.OrderWithdraw{MOrderID: req.MOrderID, MID: req.MID}
	db.Mysql().Get(one)
	if one.ID > 0 {
		a.Code = app.CodeDuplicateOrder
		return
	}
	log.Debug("PayOutCreate:%+v", req)
	amount := util.Decimal(float64(req.Amount)/100, 2)
	fee := a.Merchat.CallFee(2, amount)
	total := amount + fee
	if a.Merchat.Balance < total {
		a.Code = app.CodeInsufficientBalance
		return
	}
	orderID := call.NewOrderIDV2(common.PayOutPrefix)
	riskAction := call.CheckWithdrawCardRisk(req.MID, orderID, req.MOrderID, req.Type, req.Account, req.IP, req.Amount)
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
	if req.MID == 10000 || riskAction == call.WithdrawCardRiskActionPendingReview { // 测试号或风控待审核
		status = common.StatusROrderCreate
	}
	order := &common.OrderWithdraw{
		OrderID:     orderID,
		MOrderID:    req.MOrderID,
		MID:         req.MID,
		Amount:      amount,
		PayAccount:  tmp,
		Status:      status,
		UID:         util.AnyToString(req.UserID),
		Extra:       req.Extra,
		CallbackURL: req.NotifyURL,
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
		MID:      req.MID,
		OrderID:  orderID,
		MOrderID: req.MOrderID,
	}
	if req.MID == 10000 { // 测试号
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
	MID      int    `json:"MID" binding:"required"`
	MOrderID string `json:"MOrderID" binding:"required"`
	Sign     string `json:"Sign" binding:"required"`
}

type PayOutQueryResp struct {
	MID      int    // 商户号
	OrderID  string // 平台订单号
	MOrderID string // 商户订单号
	Amount   int64  // 订单金额
	Status   int    // 订单状态
}

func PayOutQuery(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(PayInQueryReq)
	if !a.S(req) {
		return
	}
	mer := &common.Merchats{MID: req.MID}
	db.Mysql().Get(mer)
	if mer.ID > 0 {
		a.Merchat = mer
	}
	if !a.CheckMerchat(req, req.MID, req.Sign, 4) {
		return
	}
	order := &common.OrderWithdraw{MOrderID: req.MOrderID, MID: req.MID}
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
		MID:      req.MID,
		OrderID:  order.OrderID,
		MOrderID: req.MOrderID,
		Amount:   int64(order.Amount * 100),
		Status:   order.Status,
	}
}
