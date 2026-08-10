package apiv1

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"server/call"
	"server/common"
	"server/config"
	"server/db"
	"server/modules/pay/app"
	"server/modules/pay/handler"
	"server/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
	"github.com/shopspring/decimal"
)

type PayInCreateReq struct {
	MerchantID        int    `json:"merchantId" binding:"required"`
	MerchantOrderID   string `json:"merchantOrderId" binding:"required"`
	OrderAmount       int    `json:"orderAmount" binding:"required"`
	ClientIP          string `json:"clientIp" binding:"required"`
	CallbackURL       string `json:"callbackUrl" binding:"required"`
	CustomerPhone     string `json:"customerPhone" binding:"required"`
	CustomerID        any    `json:"customerId"`
	CustomerName      string `json:"customerName"`
	CustomerEmail     string `json:"customerEmail"`
	DeviceFingerprint string `json:"deviceFingerprint"`
	Metadata          string `json:"metadata"`
	Signature         string `json:"signature" binding:"required"`
}

type PayInCreateResp struct {
	MerchantID      int    `json:"merchantId"`      // 商户号
	PlatformOrderID string `json:"platformOrderId"` // 平台订单号
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单号
	CheckoutURL     string `json:"checkoutUrl"`     // 支付链接
}

func PayInCreate(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(PayInCreateReq)
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
	if !a.CheckMerchat(req, req.MerchantID, req.Signature, 1) {
		return
	}
	if len(a.Merchat.PayRangeSub) == 2 {
		min := int(a.Merchat.PayRangeSub[0] * 100)
		max := int(a.Merchat.PayRangeSub[1] * 100)
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
	data := &common.CommonPayData{Name: req.CustomerName, Mobile: req.CustomerPhone, Email: req.CustomerEmail, IP: req.ClientIP, DeviceID: req.DeviceFingerprint}
	if !a.Risk(data) {
		return
	}
	one := &common.OrderPay{MOrderID: req.MerchantOrderID, MID: req.MerchantID}
	db.Mysql().Get(one)
	if one.ID > 0 {
		a.Code = app.CodeDuplicateOrder
		return
	}
	log.Debug("PayInCreate:%+v", req)

	tmp, _ := json.Marshal(data)
	amount := util.Decimal(float64(req.OrderAmount)/100, 2)
	orderID := call.NewOrderIDV2(common.PayInPrefix)

	if req.MerchantID == 888 { // 测试号
		order := &common.OrderPay{
			OrderID:     orderID,
			MOrderID:    req.MerchantOrderID,
			MID:         req.MerchantID,
			Amount:      amount,
			PayAccount:  tmp,
			Status:      common.StatusROrderCreate,
			UID:         util.AnyToString(req.CustomerID),
			Extra:       req.Metadata,
			CallbackURL: req.CallbackURL,
			Fee:         a.Merchat.CallFee(1, amount),
			// PayChannel:  ret.Channel,
			RealTimeSettle: true,
		}
		err := db.Mysql().Create(order)
		if err != nil {
			a.Code = app.CodeInnerError
			return
		}
		a.Data = &PayInCreateResp{
			MerchantID:      req.MerchantID,
			PlatformOrderID: orderID,
			MerchantOrderID: req.MerchantOrderID,
			CheckoutURL:     "https://www.google.com",
		}
		success := false
		if rand.Intn(2) == 0 {
			success = true
		}
		time.AfterFunc(2*time.Second, func() {
			call.PayCallback(order, success)
		})
		return
	}

	ret := handler.Recharge(&common.PayReq{
		MchId:    req.MerchantID,
		OrderID:  orderID,
		Name:     req.CustomerName,
		Mobile:   req.CustomerPhone,
		Email:    req.CustomerEmail,
		Amount:   amount,
		UID:      util.AnyToString(req.CustomerID),
		IP:       req.ClientIP,
		DeviceID: req.DeviceFingerprint,
		Merchat:  a.Merchat,
	})
	if ret.Err != nil {
		a.Code = app.CodeNoChannel
		log.Error("err:%v", ret.Err)
		return
	}

	order := &common.OrderPay{
		OrderID:     orderID,
		MOrderID:    req.MerchantOrderID,
		APIOrderID:  ret.APIOrderID,
		MID:         req.MerchantID,
		Amount:      amount,
		PayAccount:  tmp,
		Status:      common.StatusROrderCreate,
		UID:         util.AnyToString(req.CustomerID),
		Extra:       req.Metadata,
		CallbackURL: req.CallbackURL,
		Fee:         a.Merchat.CallFee(1, amount),
		PayChannel:  ret.Channel,
		PayUrl:      ret.URL,
	}
	if !ret.UseOriginURL {
		ret.URL = fmt.Sprintf("%s?order=%s", config.GetConfig().Pay.PaymentURL, orderID)
	}
	ch := call.GetConfigPayChannelByCid(order.PayChannel)
	if ch != nil && ch.RealTimeSettle {
		order.RealTimeSettle = true
	}
	err := db.Mysql().Create(order)
	if err != nil {
		a.Code = app.CodeInnerError
		return
	}
	a.Data = &PayInCreateResp{
		MerchantID:      req.MerchantID,
		PlatformOrderID: orderID,
		MerchantOrderID: req.MerchantOrderID,
		CheckoutURL:     ret.URL,
	}
}

type PayInQueryReq struct {
	MerchantID      int    `json:"merchantId" binding:"required"`
	MerchantOrderID string `json:"merchantOrderId" binding:"required"`
	Signature       string `json:"signature" binding:"required"`
}

type PayInQueryResp struct {
	MerchantID      int    `json:"merchantId"`      // 商户号
	PlatformOrderID string `json:"platformOrderId"` // 平台订单号
	MerchantOrderID string `json:"merchantOrderId"` // 商户订单号
	OrderAmount     int64  `json:"orderAmount"`     // 订单金额
	OrderStatus     int    `json:"orderStatus"`     // 订单状态
}

func PayInQuery(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(PayInQueryReq)
	if !a.S(req) {
		return
	}
	if !a.CheckMerchat(req, req.MerchantID, req.Signature, 2) {
		return
	}
	order := &common.OrderPay{MOrderID: req.MerchantOrderID, MID: req.MerchantID}
	err := db.Mysql().Get(order)
	if order.ID == 0 {
		a.Code = app.CodeOrderNotExist
		return
	}
	if err != nil {
		a.Code = app.CodeInnerError
		return
	}
	a.Data = &PayInQueryResp{
		MerchantID:      req.MerchantID,
		PlatformOrderID: order.OrderID,
		MerchantOrderID: req.MerchantOrderID,
		OrderAmount:     decimal.NewFromFloat(order.Amount).Shift(2).IntPart(),
		OrderStatus:     order.Status,
	}
}
