package base

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"server/common"
	"server/config"
	"server/modules/pay/values"
	"server/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
)

type HttpType int
type OPTType int
type CallbackType int

const (
	HttpTypeInvalid HttpType = iota
	HttpTypeJson
	HttpTypeForm
	HttpTypeMultiPartForm
	HttpTypeGet       // get请求
	HttpTypePostGet   // 带参数的get请求
	HttpTypePostParam // 拼接参数的post请求
	HttpTypeAll
)

const (
	OPTPay OPTType = iota + 1
	OPTWithdraw
	OPTPayCB
	OPTWithdrawCB
	OPTQueryWithdraw
	OPTQueryPay
)

const (
	PostformTypeBasic = iota // c.postform的模式
	PostformTypeQuery        // c.query的模式
)

const (
	CallbackTypeString CallbackType = iota
	CallbackTypeJson
)

const (
	CallbackStatusNormal = iota // 正常处理后续操作
	CallbackStatusPass          // 不需要处理
)

type Base struct {
	Channel             values.PayWay
	PayReq              *common.PayReq
	WithdrawReq         *common.WithdrawReq
	SignKey             string
	Sign                string // 签名
	ShouldSignUpper     bool
	HttpType            HttpType
	PostformType        int // 0是postform 1是query
	ReqURL              string
	Opt                 OPTType // 1支付 2退出 3支付回调 4退出回调 5退出查询 6支付查询
	Status              int     // 返回的状态
	GetResp             url.Values
	Resp                interface{}
	CallbackReq         interface{}
	CallbackResp        CallbackResp
	ShouldQueryWithdraw bool // 是否开启自动代付查询
	QueryWithdrawReq    *QueryReq
	QueryWithdrawResp   QueryResp
	QueryPayReq         *QueryReq
	QueryPayResp        QueryResp
	PayResp             *common.PayResp
	WithdrawResp        *common.WithdrawResp
	CallbackStatus      int // 回调状态判断,如果是支付中/打款中的状态不需要处理
	Sub
	SubReq            func() error                 // 特殊处理一些非常规请求
	SubCallbackDecode func([]byte) (string, error) // 特殊处理回调解析的方法
	SignPassStr       []string                     // 不参与签名的字段
	PassKeyName       bool                         // 是否签名的时候跳过key字段名(有的通道在签名后直接拼接key值,不需要&key=)
	KeyName           string                       // 有时候有些渠道签名字段不一样
	C                 *gin.Context
	ReqData           interface{}
	CallbackType      CallbackType
	CallbackRespJson  interface{}
	Body              string
	PostBody          interface{} //携带了一个请求体，因为有的渠道需要签名的时候签名在头部，但是签名需要同时拿到请求体和请求头
	Route             func() bool // 路由到其他服务器
	ESS1              string      // 额外存储字段1
	ESS2              string      // 额外存储字段2
	UseOriginURL      bool        // 如果为true,返回三方的原始url
}

type CallbackResp struct {
	Msg         string
	OrderID     string
	APIOrderID  string
	FailMessage string
	Utr         string
	Success     bool
}

const (
	QueryStatusFail         = iota + 1 // 查询到订单失败了或者提交失败
	QueryStatusSuccess                 // 查询到订单提交成功
	QueryStatusUnknown                 // 查询失败,状态未知,需要重新查询
	QueryStatusOrderSuccess            // 查询到订单支付成功
)

type QueryReq struct {
	Channel    int
	OrderID    string
	APIOrderID string
}

type QueryResp struct {
	Msg        string
	OrderID    string
	APIOrderID string
	Status     int
}

type Sub interface {
	PackHeader(header http.Header) // 生成请求头
	PackReq() interface{}          // 生成支付请求体
	GetResp() error                // 获取返回
	CheckSign(str string) bool     // 验签同时给callbackresp赋值
}

func NewRechargeBase(req *common.PayReq) *Base {
	values.PackPay(req)
	base := &Base{
		Channel:     values.PayWay(req.Channel),
		PayReq:      req,
		Opt:         OPTPay,
		SignPassStr: []string{"sign"},
	}
	return base
}

func NewWithdrawBase(req *common.WithdrawReq) *Base {
	values.PackWithdraw(req)
	base := &Base{
		Channel:     values.PayWay(req.Channel),
		WithdrawReq: req,
		Opt:         OPTWithdraw,
	}
	return base
}

func NewQueryWithdrawBase(req *QueryReq) *Base {
	base := &Base{
		Channel:          values.PayWay(req.Channel),
		QueryWithdrawReq: req,
		Opt:              OPTQueryWithdraw,
	}
	return base
}

func NewQueryPayBase(req *QueryReq) *Base {
	base := &Base{
		Channel:     values.PayWay(req.Channel),
		QueryPayReq: req,
		Opt:         OPTQueryPay,
	}
	return base
}

func NewCallbackBase(opt OPTType) *Base {
	base := &Base{
		Opt: opt,
	}
	return base
}

func (b *Base) PackHeader(header http.Header) {
	switch b.HttpType {
	case HttpTypeJson:
		header.Set("Content-Type", "application/json")
	case HttpTypeForm:
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	b.Sub.PackHeader(header)
}

func (b *Base) Req() error {
	if b.SubReq != nil {
		return b.SubReq()
	}
	send := b.Sub.PackReq()
	if send == nil {
		return errors.New("req fail")
	}
	var req *http.Request
	var err error
	if b.HttpType == HttpTypeGet {
		req, _ = http.NewRequest("GET", b.ReqURL, nil)
	} else if b.HttpType == HttpTypePostGet {
		ref := reflect.ValueOf(send)
		reft := reflect.TypeOf(send)
		if ref.Kind() == reflect.Ptr {
			ref = ref.Elem()
			reft = reft.Elem()
		}
		data := url.Values{}
		for i := 0; i < ref.NumField(); i++ {
			if ref.Field(i).IsZero() {
				continue
			}
			data.Set(reft.Field(i).Tag.Get("json"), fmt.Sprintf("%v", ref.Field(i).Interface()))
		}
		str := data.Encode()
		req, err = http.NewRequest("GET", b.ReqURL+"?"+str, nil)
		if err != nil {
			log.Error("err:%v", err)
			return errors.New("req fail")
		}
	} else if b.HttpType == HttpTypePostParam {
		ref := reflect.ValueOf(send)
		reft := reflect.TypeOf(send)
		if ref.Kind() == reflect.Ptr {
			ref = ref.Elem()
			reft = reft.Elem()
		}
		data := url.Values{}
		for i := 0; i < ref.NumField(); i++ {
			if ref.Field(i).IsZero() {
				continue
			}
			data.Set(reft.Field(i).Tag.Get("json"), fmt.Sprintf("%v", ref.Field(i).Interface()))
		}
		str := data.Encode()
		req, err = http.NewRequest("POST", b.ReqURL+"?"+str, nil)
		if err != nil {
			log.Error("err:%v", err)
			return errors.New("req fail")
		}
	} else {
		reqData := b.PackReq(send)
		req, err = http.NewRequest("POST", b.ReqURL, reqData)
		if err != nil {
			log.Error("err:%v", err)
			return errors.New("req fail")
		}
	}
	b.PostBody = send
	b.PackHeader(req.Header)
	b.Status = b.Request(req, b.Resp)
	err = b.Sub.GetResp()
	if err != nil {
		backReq := &values.BackReq{
			Opt:      int(b.Opt),
			Channel:  int(b.Channel),
			ErrorMsg: string(b.Body),
		}
		util.Go(func() {
			util.HttpPostWithoutRet(config.GetConfig().Pay.BackSysAddr+"/errNotify", backReq, nil)
		})
		log.Error("err:%v", err)
		return err
	}
	return nil
}

func (b *Base) PackReq(send interface{}) io.Reader {
	ref := reflect.ValueOf(send)
	reft := reflect.TypeOf(send)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
		reft = reft.Elem()
	}
	switch b.HttpType {
	case HttpTypeJson:
		reqData, _ := json.Marshal(send)
		log.Debug("post to:%v,req:%v", b.ReqURL, string(reqData))
		return bytes.NewBuffer(reqData)
	case HttpTypeForm:
		data := url.Values{}
		for i := 0; i < ref.NumField(); i++ {
			if ref.Field(i).IsZero() {
				continue
			}
			data.Set(reft.Field(i).Tag.Get("json"), fmt.Sprintf("%v", ref.Field(i).Interface()))
		}
		str := data.Encode()
		log.Debug("post to:%v,req:%v", b.ReqURL, str)
		return strings.NewReader(str)
	case HttpTypeMultiPartForm:
		data := &bytes.Buffer{}
		writer := multipart.NewWriter(data)
		for i := 0; i < ref.NumField(); i++ {
			if ref.Field(i).IsZero() {
				continue
			}
			writer.WriteField(reft.Field(i).Tag.Get("json"), fmt.Sprintf("%v", ref.Field(i).Interface()))
		}
		log.Debug("post to:%v,req:%v", b.ReqURL, data.String())
		return data
	}
	return nil
}

func (b *Base) GetPayCallbackURL() string {
	return values.GetPayCallback(b.Channel)
}

func (b *Base) GetWithdrawCallbackURL() string {
	return values.GetWithdrawCallback(b.Channel)
}

// func (b *Base) Proxy() {
// 	if !b.ShouldProxy || b.Opt != OPTPay {
// 		return
// 	}
// 	b.PayResp.URL = ReplaceDomain(b.PayResp.URL, fmt.Sprintf("payment%d.lh-pays.com", b.Channel))
// }
