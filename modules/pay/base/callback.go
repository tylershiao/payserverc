package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"server/call"
	"server/common"
	"server/db"
	"server/modules/pay/values"
	"server/util/rlock"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
)

func (b *Base) PayCallback(c *gin.Context) {
	str, err := b.Callback(c)
	defer func() {
		if b.CallbackType == CallbackTypeJson {
			c.JSON(http.StatusOK, b.CallbackRespJson)
		} else {
			c.String(http.StatusOK, b.CallbackResp.Msg)
		}
	}()
	if err != nil {
		log.Error("err:%v", err)
		b.CallbackResp.Msg = "callback fail"
		return
	}
	b.C = c
	if b.Route != nil && b.Route() {
		return
	}
	if !b.CheckSign(str) {
		if b.CallbackStatus == CallbackStatusPass {
			return
		}
		log.Error("check sign fail")
		b.CallbackResp.Msg = "callback fail"
		return
	}
	ctx := context.Background()
	rl, err := rlock.NewLockForRedis(ctx, db.Redis().GetRedis(), time.Duration(10), time.Duration(8))
	if err != nil {
		log.Error("rlock err:%v", err)
		return
	}

	url, err := rl.Lock(ctx, fmt.Sprintf("payin_order:%s", b.CallbackResp.OrderID))
	if err != nil {
		log.Error("call rl.Lock func err: %v, order_no: %s", err, b.CallbackResp.OrderID)
		return
	}
	defer url()
	orderID := b.CallbackResp.OrderID
	if orderID == "" {
		log.Error("order:%v not exist", orderID)
		return
	}
	or := &common.OrderPay{OrderID: orderID}
	db.Mysql().Get(or)
	if or.ID == 0 {
		log.Error("order:%v not exist", orderID)
		return
	}
	if or.Status == common.StatusROrderFinish {
		log.Error("order:%v exec done", orderID)
		return
	}
	success := b.CallbackResp.Success
	if or.APIOrderID == "" {
		or.APIOrderID = b.CallbackResp.APIOrderID
	}
	if err := call.PayCallback(or, success); err != nil {
		b.CallbackResp.Msg = "callback fail"
		return
	}
	if success && (time.Now().Unix()-or.CreatedAt.Unix() <= 300) { // 5分钟内成功的单才增加权重
		values.PayCallback(b.Channel)
	}
}

func (b *Base) WithdrawCallback(c *gin.Context) {
	str, err := b.Callback(c)
	defer func() {
		if b.CallbackType == CallbackTypeJson {
			c.JSON(http.StatusOK, b.CallbackRespJson)
		} else {
			c.String(http.StatusOK, b.CallbackResp.Msg)
		}
	}()
	if err != nil {
		log.Error("err:%v", err)
		b.CallbackResp.Msg = "callback fail"
		return
	}
	b.C = c
	if b.Route != nil && b.Route() {
		return
	}
	if !b.CheckSign(str) {
		if b.CallbackStatus == CallbackStatusPass {
			return
		}
		log.Error("check sign fail")
		b.CallbackResp.Msg = "callback fail"
		return
	}
	ctx := context.Background()
	rl, err := rlock.NewLockForRedis(ctx, db.Redis().GetRedis(), time.Duration(10), time.Duration(8))
	if err != nil {
		log.Error("rlock err:%v", err)
		return
	}

	url, err := rl.Lock(ctx, fmt.Sprintf("payout_order:%s", b.CallbackResp.OrderID))
	if err != nil {
		log.Error("call rl.Lock func err: %v, order_no: %s", err, b.CallbackResp.OrderID)
		return
	}
	defer url()
	orderID := b.CallbackResp.OrderID
	or := &common.OrderWithdraw{OrderID: orderID}
	db.Mysql().Get(or)
	if or.ID == 0 {
		log.Error("order:%v not exist", orderID)
		return
	}
	if or.PayChannel != int(b.Channel) {
		log.Error("order:%v,pay channel:%v,this channel:%v", orderID, or.PayChannel, b.Channel)
		return
	}
	if or.MOrderID == "" {
		or.MOrderID = b.CallbackResp.APIOrderID
	}
	success := b.CallbackResp.Success
	or.Message = b.CallbackResp.FailMessage
	or.Utr = b.CallbackResp.Utr
	originStatus := or.Status
	if success {
		or.Status = common.StatusROrderFinish
	} else {
		or.Status = common.StatusROrderFail
	}
	if err := call.WithdrawCallback(or, originStatus); err != nil {
		b.CallbackResp.Msg = "callback fail"
		return
	}
}

func (b *Base) Callback(c *gin.Context) (str string, err error) {
	if b.HttpType == HttpTypeJson {
		str, err = b.CallbackJson(c)
	} else if b.HttpType == HttpTypeForm {
		str, err = b.CallbackForm(c)
	} else if b.HttpType == HttpTypeGet {
		str, err = b.CallbackGet(c)
	}
	return
}

func (b *Base) CallbackJson(c *gin.Context) (str string, err error) {
	var tmp []byte
	tmp, err = io.ReadAll(c.Request.Body)
	defer func() {
		c.Request.Body.Close()
	}()
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	str = string(tmp)
	log.Debug("Callback body:%v", str)
	if b.SubCallbackDecode != nil {
		return b.SubCallbackDecode(tmp)
	}
	err = json.Unmarshal(tmp, b.CallbackReq)
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	return
}

func (b *Base) CallbackForm(c *gin.Context) (str string, err error) {
	req := b.CallbackReq
	typ := reflect.TypeOf(req).Elem()
	val := reflect.ValueOf(req).Elem()
	for i := 0; i < typ.NumField(); i++ {
		var fieldVal string
		if b.PostformType == PostformTypeQuery {
			fieldVal = c.Query(typ.Field(i).Tag.Get("json"))
		} else {
			fieldVal = c.PostForm(typ.Field(i).Tag.Get("json"))
		}
		kind := val.Field(i).Kind()
		if kind == reflect.Int64 || kind == reflect.Int {
			tmp, _ := strconv.ParseInt(fieldVal, 10, 64)
			val.Field(i).SetInt(tmp)
		} else {
			val.Field(i).SetString(fieldVal)
		}
	}
	str = values.GetSignStrForm(c, b.SignPassStr...)
	return
}

func (b *Base) CallbackGet(c *gin.Context) (str string, err error) {
	raw, err := url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		log.Error("err:%v", err)
	}
	b.GetResp = raw
	req := b.CallbackReq
	typ := reflect.TypeOf(req).Elem()
	val := reflect.ValueOf(req).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldVal := b.GetResp.Get(typ.Field(i).Tag.Get("json"))
		kind := val.Field(i).Kind()
		switch kind {
		case reflect.String:
			val.Field(i).SetString(fieldVal)
		case reflect.Int, reflect.Int64:
			tmp, _ := strconv.ParseInt(fieldVal, 10, 64)
			val.Field(i).SetInt(tmp)
		case reflect.Float64:
			tmp, _ := strconv.ParseFloat(fieldVal, 64)
			val.Field(i).SetFloat(tmp)
		}
	}
	str = values.GetSignStrGet(b.GetResp, "sign")
	return
}
