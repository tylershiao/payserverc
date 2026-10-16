package masterpay

import (
	"encoding/json"
	"server/common"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/sign"
	"strings"
	"testing"
)

func TestGetRespCopiesCardInfoLink(t *testing.T) {
	b := &base.Base{
		Opt:     base.OPTPay,
		Channel: values.MasterPay,
		Resp: &PayResp{
			RetCode:    "SUCCESS",
			PayOrderID: "master-order-1",
			PayURL:     "https://pay.example.com/order/1",
			CardInfo: &CardInfo{
				Link: "upi://pay?pa=test@example.com",
			},
		},
	}

	if err := (&Sub{Base: b}).GetResp(); err != nil {
		t.Fatalf("GetResp() error = %v", err)
	}
	if b.PayResp == nil {
		t.Fatal("GetResp() did not set PayResp")
	}
	if got, want := b.PayResp.Link, "upi://pay?pa=test@example.com"; got != want {
		t.Fatalf("PayResp.Link = %q, want %q", got, want)
	}
}

func TestGetRespAllowsMissingCardInfo(t *testing.T) {
	b := &base.Base{
		Opt:     base.OPTPay,
		Channel: values.MasterPay,
		Resp: &PayResp{
			RetCode:    "SUCCESS",
			PayOrderID: "master-order-2",
			PayURL:     "https://pay.example.com/order/2",
		},
	}

	if err := (&Sub{Base: b}).GetResp(); err != nil {
		t.Fatalf("GetResp() error = %v", err)
	}
	if b.PayResp == nil {
		t.Fatal("GetResp() did not set PayResp")
	}
	if b.PayResp.Link != "" {
		t.Fatalf("PayResp.Link = %q, want empty string", b.PayResp.Link)
	}
}

func TestPackPayReqRequestsCardInfo(t *testing.T) {
	b := &base.Base{
		Opt:     base.OPTPay,
		Channel: values.MasterPay,
		PayReq: &common.PayReq{
			OrderID: "order-1",
			Amount:  100,
			IP:      "127.0.0.1",
		},
	}
	NewSub(b)
	b.SignKey = "test-secret"
	req := b.Sub.PackReq().(*PayReq)
	encoded, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["requestCardInfo"] != true {
		t.Fatalf("requestCardInfo = %v, want true", payload["requestCardInfo"])
	}
	if signStr := sign.GetSignStr(req, "sign"); !strings.Contains(signStr, "requestCardInfo=true") {
		t.Fatalf("requestCardInfo is missing from signature payload")
	}
	if req.Sign == "" || req.Sign != b.SignMD5(req) {
		t.Fatal("request signature does not cover the completed request")
	}
}
