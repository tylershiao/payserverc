package deepay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math"
	"net/http"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/util"
	"strings"
	"time"

	"github.com/liangdas/mqant/log"
)

func NewSub(b *base.Base) {
	sub := &Sub{Base: b}
	b.HttpType = base.HttpTypeJson
	b.CallbackResp.Msg = "OK"
	b.UseOriginURL = true
	switch b.Opt {
	case base.OPTPay:
		b.Resp = new(PayResp)
		b.ReqURL = baseURL + payURL
	case base.OPTWithdraw:
		b.Resp = new(WithdrawResp)
		b.ReqURL = baseURL + withdrawURL
	case base.OPTPayCB:
		b.CallbackReq = new(PayCallbackReq)
	case base.OPTWithdrawCB:
		b.CallbackReq = new(WithdrawCallbackReq)
	case base.OPTQueryWithdraw:
		b.Resp = new(QueryWithdrawResp)
		b.ReqURL = baseURL + queryWithdrawURL
	case base.OPTQueryPay:
		b.Resp = new(QueryPayResp)
		b.ReqURL = baseURL + queryPayURL
	}
	b.Sub = sub
}

type Sub struct {
	Base *base.Base
}

func (s *Sub) PackHeader(header http.Header) {
	payload, err := json.Marshal(s.Base.PostBody)
	if err != nil {
		log.Error("deepay marshal sign payload err:%v", err)
		return
	}
	sign, err := createSign(payload, privateKey)
	if err != nil {
		log.Error("deepay sign err:%v", err)
		return
	}
	header.Set("X-Sign", sign)
}

func (s *Sub) PackReq() interface{} {
	switch s.Base.Opt {
	case base.OPTPay:
		return s.PackPayReq()
	case base.OPTWithdraw:
		return s.PackWithdrawReq()
	case base.OPTQueryWithdraw:
		return s.PackQueryWithdrawReq()
	case base.OPTQueryPay:
		return s.PackQueryPayReq()
	default:
		return nil
	}
}

func (s *Sub) GetResp() error {
	switch s.Base.Opt {
	case base.OPTPay:
		resp := s.Base.Resp.(*PayResp)
		if resp.Status != http.StatusOK || resp.Data == nil {
			return errors.New(resp.Msg)
		}
		payURL := resp.Data.Link
		if payURL == "" {
			payURL = resp.Data.DeepLink
		}
		if payURL == "" {
			return errors.New("deepay pay url empty")
		}
		s.Base.PayResp = &common.PayResp{APIOrderID: resp.Data.PlatOrderSn, URL: payURL, Message: resp.Msg, Channel: int(s.Base.Channel)}
	case base.OPTQueryPay:
		resp := s.Base.Resp.(*QueryPayResp)
		s.Base.QueryPayResp.Msg = resp.Msg
		if resp.Data != nil {
			s.Base.QueryPayResp.OrderID = resp.Data.TradeSn
			s.Base.QueryPayResp.APIOrderID = resp.Data.PlatOrderSn
			if resp.Status == http.StatusOK && resp.Data.OrderStatus == "SUCCESS" && amountsEqual(resp.Data.Amount, resp.Data.ActualAmount) {
				s.Base.QueryPayResp.Status = base.QueryStatusOrderSuccess
			} else if resp.Status == http.StatusOK && (resp.Data.OrderStatus == "FAILED" || resp.Data.OrderStatus == "REFUND") {
				s.Base.QueryPayResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryPayResp.Status == 0 {
			s.Base.QueryPayResp.Status = base.QueryStatusUnknown
		}
	case base.OPTWithdraw:
		resp := s.Base.Resp.(*WithdrawResp)
		if resp.Status != http.StatusOK || resp.Data == nil {
			return errors.New(resp.Msg)
		}
		s.Base.WithdrawResp = &common.WithdrawResp{APIOrderID: resp.Data.PlatOrderSn, Message: resp.Msg, Channel: int(s.Base.Channel)}
	case base.OPTQueryWithdraw:
		resp := s.Base.Resp.(*QueryWithdrawResp)
		s.Base.QueryWithdrawResp.Msg = resp.Msg
		if resp.Data != nil {
			s.Base.QueryWithdrawResp.OrderID = resp.Data.TradeSn
			s.Base.QueryWithdrawResp.APIOrderID = resp.Data.PlatOrderSn
			if resp.Status == http.StatusOK {
				switch resp.Data.OrderStatus {
				case "SUCCESS":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusOrderSuccess
				case "PENDING":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusSuccess
				case "FAILED", "REFUND":
					s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
				}
			} else if resp.Data.OrderStatus == "FAILED" || resp.Data.OrderStatus == "REFUND" {
				s.Base.QueryWithdrawResp.Status = base.QueryStatusFail
			}
		}
		if s.Base.QueryWithdrawResp.Status == 0 {
			s.Base.QueryWithdrawResp.Status = base.QueryStatusUnknown
		}
	default:
		return errors.New("unknown opt")
	}
	return nil
}

func (s *Sub) PackPayReq() interface{} {
	r := s.Base.PayReq
	return &PayReq{
		TradeSn:   r.OrderID,
		Currency:  currency,
		Amount:    r.Amount,
		Name:      normalizeName(r.Name),
		Phone:     util.CheckPhone(r.Mobile),
		Email:     r.Email,
		Timestamp: time.Now().Unix(),
		Mid:       mid,
		Remarks:   remarks,
		NotifyURL: s.Base.GetPayCallbackURL(),
		ReturnURL: values.GetFrontCallback(),
		TradeCode: payTradeCode,
	}
}

func (s *Sub) PackWithdrawReq() interface{} {
	r := s.Base.WithdrawReq
	send := &WithdrawReq{
		TradeSn:   r.OrderID,
		Currency:  currency,
		Amount:    r.Amount,
		Account:   r.Account,
		NotifyURL: s.Base.GetWithdrawCallbackURL(),
		Name:      normalizeName(r.Name),
		Phone:     util.CheckPhone(r.Mobile),
		Email:     r.Email,
		Timestamp: time.Now().Unix(),
		Mid:       mid,
		Remarks:   remarks,
	}
	if r.PayType == common.WithdrawTypeBank {
		send.TradeCode = bankTradeCode
		send.IFSC = r.IFSC
		if send.IFSC == "" {
			return nil
		}
		return send
	}
	if r.PayType == common.WithdrawTypeUPI {
		send.TradeCode = upiTradeCode
		return send
	}
	return nil
}

func (s *Sub) PackQueryWithdrawReq() interface{} {
	r := s.Base.QueryWithdrawReq
	return &QueryWithdrawReq{
		Mid:       mid,
		TradeSn:   r.OrderID,
		Timestamp: time.Now().Unix(),
	}
}

func (s *Sub) PackQueryPayReq() interface{} {
	r := s.Base.QueryPayReq
	return &QueryPayReq{
		Mid:       mid,
		TradeSn:   r.OrderID,
		Timestamp: time.Now().Unix(),
	}
}

func (s *Sub) CheckSign(str string) bool {
	checkSign := s.Base.C.GetHeader("X-Sign")
	if checkSign == "" {
		log.Error("deepay callback sign missing")
		return false
	}
	if err := base.RsaSignVer([]byte(str), checkSign, publicKey); err != nil {
		log.Error("deepay verify sign err:%v", err)
		return false
	}
	switch s.Base.Opt {
	case base.OPTPayCB:
		req := s.Base.CallbackReq.(*PayCallbackReq)
		log.Debug("deepay pay callback:%+v", *req)
		s.Base.CallbackResp.OrderID = req.TradeSn
		s.Base.CallbackResp.APIOrderID = req.PlatOrderSn
		s.Base.CallbackResp.FailMessage = req.ErrMsg
		if req.OrderStatus == "PENDING" {
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		s.Base.CallbackResp.Success = req.OrderStatus == "SUCCESS" && amountsEqual(req.Amount, req.ActualAmount)
		if req.OrderStatus == "SUCCESS" && !s.Base.CallbackResp.Success && s.Base.CallbackResp.FailMessage == "" {
			s.Base.CallbackResp.FailMessage = "deepay actual amount mismatch"
		}
	case base.OPTWithdrawCB:
		req := s.Base.CallbackReq.(*WithdrawCallbackReq)
		log.Debug("deepay withdraw callback:%+v", *req)
		s.Base.CallbackResp.OrderID = req.TradeSn
		s.Base.CallbackResp.APIOrderID = req.PlatOrderSn
		s.Base.CallbackResp.FailMessage = req.ErrMsg
		s.Base.CallbackResp.Utr = req.BankRefNo
		if req.OrderStatus != "SUCCESS" && req.OrderStatus != "FAILED" && req.OrderStatus != "REFUND" { // 未知状态
			s.Base.CallbackStatus = base.CallbackStatusPass
			return false
		}
		s.Base.CallbackResp.Success = req.OrderStatus == "SUCCESS"
	default:
		return false
	}
	return true
}

func createSign(data []byte, key string) (string, error) {
	privateKey, err := parsePrivateKey(key)
	if err != nil {
		return "", err
	}
	hash := sha1.Sum(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func parsePrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("invalid private key")
	}
	if pkcs1Key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return pkcs1Key, nil
	}
	pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := pkcs8Key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid rsa private key")
	}
	return privateKey, nil
}

func normalizeName(name string) string {
	clean := strings.TrimSpace(name)
	clean = strings.Map(func(r rune) rune {
		if r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return r
		}
		return -1
	}, clean)
	clean = strings.TrimSpace(clean)
	if len(clean) > 5 {
		return clean
	}
	if clean == "" {
		return "Pay User"
	}
	return clean + " User"
}

func amountsEqual(left, right float64) bool {
	return math.Abs(left-right) < 0.000001
}
