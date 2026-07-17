package sys

import (
	"fmt"
	"server/call"
	"server/common"
	"server/config"
	"server/db"
	"server/modules/pay/allpay"
	"server/modules/pay/app"
	"server/modules/pay/base"
	"server/util"

	"github.com/gin-gonic/gin"
)

type ReloadReq struct {
	Type int
}

func Reload(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(ReloadReq)
	if !a.S(req) {
		return
	}
	call.ReloadConfig(req.Type)
}

type TestPayReq struct {
	PayChannelID int
	Amount       float64
}

type TestPayResp struct {
	PayURL   string // 支付链接
	ErrorMsg string // 报错信息
}

// 测试拉单
func TestPay(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(TestPayReq)
	if !a.S(req) {
		return
	}
	name := util.GenerateRandomString(5)
	payReq := &common.PayReq{
		MchId:   1,
		Channel: req.PayChannelID,
		OrderID: call.NewOrderIDV2(common.PayInPrefix),
		Amount:  req.Amount,
		UID:     "1921",
		Name:    name,
		Email:   fmt.Sprintf("%s@gmail.com", name),
		Mobile:  util.CheckPhone(""),
	}
	resp := &TestPayResp{}
	a.Data = resp
	base := base.NewRechargeBase(payReq)
	allpay.NewSub(base, req.PayChannelID)
	if base.Sub == nil {
		resp.ErrorMsg = "该通道不支持代收"
		a.Code = 1
		return
	}
	err := base.Req()
	if err != nil {
		resp.ErrorMsg = base.Body
		a.Code = 1
	} else {
		order := &common.OrderPay{
			OrderID:    payReq.OrderID,
			MOrderID:   payReq.OrderID,
			APIOrderID: base.PayResp.APIOrderID,
			MID:        payReq.MchId,
			Amount:     req.Amount,
			// PayAccount:  tmp,
			Status: common.StatusROrderCreate,
			UID:    util.AnyToString(payReq.UID),
			// Extra:       req.Extra,
			// CallbackURL: req.NotifyURL,
			// Fee:         a.Merchat.CallFee(1, req.Amount),
			PayChannel: payReq.Channel,
			PayUrl:     base.PayResp.URL,
		}
		paymentUrl := base.PayResp.URL
		if !base.UseOriginURL {
			paymentUrl = fmt.Sprintf("%s?order=%s", config.GetConfig().Pay.PaymentURL, payReq.OrderID)
		}
		resp.PayURL = paymentUrl
		err := db.Mysql().Create(order)
		if err != nil {
			resp.ErrorMsg = "测试拉单失败"
			a.Code = 1
			return
		}
	}
}
