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
)

type PayInCreateReq struct {
	MID        int    `json:"MID" binding:"required"`
	MOrderID   string `json:"MOrderID" binding:"required"`
	Amount     int    `json:"Amount" binding:"required"`
	IP         string `json:"IP" binding:"required"`
	NotifyURL  string `json:"NotifyURL" binding:"required"`
	UserMobile string `json:"UserMobile" binding:"required"`
	UserID     any    `json:"UserID"`
	UserName   string `json:"UserName"`
	UserEmail  string `json:"UserEmail"`
	DeviceID   string `json:"DeviceID"`
	Extra      string `json:"Extra"`
	Sign       string `json:"Sign" binding:"required"`
}

type PayInCreateResp struct {
	MID      int    // 商户号
	OrderID  string // 平台订单号
	MOrderID string // 商户订单号
	PayLink  string // 支付链接
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
	if !a.CheckMerchat(req, req.MID, req.Sign, 1) {
		return
	}
	if len(a.Merchat.PayRangeSub) == 2 {
		min := int(a.Merchat.PayRangeSub[0] * 100)
		max := int(a.Merchat.PayRangeSub[1] * 100)
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
	data := &common.CommonPayData{Name: req.UserName, Mobile: req.UserMobile, Email: req.UserEmail, IP: req.IP, DeviceID: req.DeviceID}
	if !a.Risk(data) {
		return
	}
	one := &common.OrderPay{MOrderID: req.MOrderID, MID: req.MID}
	db.Mysql().Get(one)
	if one.ID > 0 {
		a.Code = app.CodeDuplicateOrder
		return
	}
	log.Debug("PayInCreate:%+v", req)

	tmp, _ := json.Marshal(data)
	amount := util.Decimal(float64(req.Amount)/100, 2)
	orderID := call.NewOrderIDV2(common.PayInPrefix)

	if req.MID == 10000 { // 测试号
		order := &common.OrderPay{
			OrderID:     orderID,
			MOrderID:    req.MOrderID,
			MID:         req.MID,
			Amount:      amount,
			PayAccount:  tmp,
			Status:      common.StatusROrderCreate,
			UID:         util.AnyToString(req.UserID),
			Extra:       req.Extra,
			CallbackURL: req.NotifyURL,
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
			MID:      req.MID,
			OrderID:  orderID,
			MOrderID: req.MOrderID,
			PayLink:  "https://www.google.com",
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
		MchId:    req.MID,
		OrderID:  orderID,
		Name:     req.UserName,
		Mobile:   req.UserMobile,
		Email:    req.UserEmail,
		Amount:   amount,
		UID:      util.AnyToString(req.UserID),
		IP:       req.IP,
		DeviceID: req.DeviceID,
		Merchat:  a.Merchat,
	})
	if ret.Err != nil {
		a.Code = app.CodeNoChannel
		log.Error("err:%v", ret.Err)
		return
	}

	order := &common.OrderPay{
		OrderID:     orderID,
		MOrderID:    req.MOrderID,
		APIOrderID:  ret.APIOrderID,
		MID:         req.MID,
		Amount:      amount,
		PayAccount:  tmp,
		Status:      common.StatusROrderCreate,
		UID:         util.AnyToString(req.UserID),
		Extra:       req.Extra,
		CallbackURL: req.NotifyURL,
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
		MID:      req.MID,
		OrderID:  orderID,
		MOrderID: req.MOrderID,
		PayLink:  ret.URL,
	}
}

type PayInQueryReq struct {
	MID      int    `json:"MID" binding:"required"`
	MOrderID string `json:"MOrderID" binding:"required"`
	Sign     string `json:"Sign" binding:"required"`
}

type PayInQueryResp struct {
	MID      int    // 商户号
	OrderID  string // 平台订单号
	MOrderID string // 商户订单号
	Amount   int64  // 订单金额
	Status   int    // 订单状态
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
	if !a.CheckMerchat(req, req.MID, req.Sign, 2) {
		return
	}
	order := &common.OrderPay{MOrderID: req.MOrderID, MID: req.MID}
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
		MID:      req.MID,
		OrderID:  order.OrderID,
		MOrderID: req.MOrderID,
		Amount:   int64(order.Amount * 100),
		Status:   order.Status,
	}
}
