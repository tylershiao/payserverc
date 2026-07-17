package allpay

import (
	"reflect"
	"server/modules/pay/base"
	"server/modules/pay/paychannels/aipay"
	"server/modules/pay/paychannels/aipay2"
	"server/modules/pay/paychannels/crownpay"
	"server/modules/pay/paychannels/decentpay"
	"server/modules/pay/paychannels/deepay"
	"server/modules/pay/paychannels/deeplink"
	"server/modules/pay/paychannels/dkpay"
	"server/modules/pay/paychannels/eastpay"
	"server/modules/pay/paychannels/epay"
	"server/modules/pay/paychannels/epay178"
	"server/modules/pay/paychannels/epayglobal"
	"server/modules/pay/paychannels/gogopay"
	"server/modules/pay/paychannels/gogopayhx"
	"server/modules/pay/paychannels/hdpay"
	"server/modules/pay/paychannels/hexpay"
	"server/modules/pay/paychannels/hongxinpay"
	"server/modules/pay/paychannels/kingpay"
	"server/modules/pay/paychannels/lulupay"
	"server/modules/pay/paychannels/masterpay"
	"server/modules/pay/paychannels/megipay"
	"server/modules/pay/paychannels/mlpay"
	"server/modules/pay/paychannels/mqr"
	"server/modules/pay/paychannels/mvpay"
	"server/modules/pay/paychannels/npay"
	"server/modules/pay/paychannels/passpay"
	"server/modules/pay/paychannels/penguinpay"
	"server/modules/pay/paychannels/primepay"
	"server/modules/pay/paychannels/propay"
	"server/modules/pay/paychannels/redpay"
	"server/modules/pay/paychannels/richpay"
	"server/modules/pay/paychannels/simplypay"
	"server/modules/pay/paychannels/swaypay"
	"server/modules/pay/paychannels/sypay"
	"server/modules/pay/paychannels/sypay2"
	"server/modules/pay/paychannels/tkpay"
	"server/modules/pay/paychannels/toplemonpay"
	transafepay "server/modules/pay/paychannels/transafepay"
	"server/modules/pay/paychannels/velopay"
	"server/modules/pay/paychannels/virgopay"
	"server/modules/pay/paychannels/virgopay2"
	"server/modules/pay/paychannels/wandapay"
	"server/modules/pay/paychannels/wheatpay"
	"server/modules/pay/paychannels/wkpay"
	"server/modules/pay/paychannels/worth_pay"
	"server/modules/pay/paychannels/zoompay"
	"server/modules/pay/values"

	"github.com/liangdas/mqant/log"
)

type AllPay struct {
	Invalid       func(b *base.Base)
	TKPay         func(b *base.Base)
	ProPay        func(b *base.Base)
	WandaPay      func(b *base.Base)
	HongxinPay    func(b *base.Base)
	MQR           func(b *base.Base)
	RichPay       func(b *base.Base)
	WorthPay      func(b *base.Base)
	DeepLink      func(b *base.Base)
	VirgoPay      func(b *base.Base)
	VirgoPay2     func(b *base.Base)
	Toplemonpay   func(b *base.Base)
	Passpay       func(b *base.Base)
	MvPay         func(b *base.Base)
	RedPay        func(b *base.Base)
	WkPay         func(b *base.Base)
	MasterPay     func(b *base.Base)
	KingPay       func(b *base.Base)
	TransafePay   func(b *base.Base)
	HdPay         func(b *base.Base)
	NPay          func(b *base.Base)
	DecentPay     func(b *base.Base)
	ZoomPay       func(b *base.Base)
	SimplyPay     func(b *base.Base)
	AIPay         func(b *base.Base)
	SwayPay       func(b *base.Base)
	EPay          func(b *base.Base)
	GogoPay       func(b *base.Base)
	HexPay        func(b *base.Base)
	PrimePay      func(b *base.Base)
	EpayGlobalPay func(b *base.Base)
	Wheatpay      func(b *base.Base)
	Lulupay       func(b *base.Base)
	Velopay       func(b *base.Base)
	SyPay         func(b *base.Base)
	CrownPay      func(b *base.Base)
	EastPay       func(b *base.Base)
	SyPay2        func(b *base.Base)
	Penguinpay    func(b *base.Base)
	Mlpay         func(b *base.Base)
	AIPay2        func(b *base.Base)
	EPay178       func(b *base.Base)
	Deepay        func(b *base.Base)
	DKPay         func(b *base.Base)
	GogoPayhx     func(b *base.Base)
	Megipay       func(b *base.Base)
}

var All = &AllPay{}

func init() {
	All.TKPay = tkpay.NewSub
	All.ProPay = propay.NewSub
	All.WandaPay = wandapay.NewSub
	All.HongxinPay = hongxinpay.NewSub
	All.MQR = mqr.NewSub
	All.RichPay = richpay.NewSub
	All.WorthPay = worth_pay.NewSub
	All.DeepLink = deeplink.NewSub
	All.VirgoPay = virgopay.NewSub
	All.VirgoPay2 = virgopay2.NewSub
	All.Toplemonpay = toplemonpay.NewSub
	All.Passpay = passpay.NewSub
	All.MvPay = mvpay.NewSub
	All.RedPay = redpay.NewSub
	All.WkPay = wkpay.NewSub
	All.MasterPay = masterpay.NewSub
	All.KingPay = kingpay.NewSub
	All.TransafePay = transafepay.NewSub
	All.HdPay = hdpay.NewSub
	All.NPay = npay.NewSub
	All.DecentPay = decentpay.NewSub
	All.ZoomPay = zoompay.NewSub
	All.SimplyPay = simplypay.NewSub
	All.AIPay = aipay.NewSub
	All.SwayPay = swaypay.NewSub
	All.EPay = epay.NewSub
	All.GogoPay = gogopay.NewSub
	All.HexPay = hexpay.NewSub
	All.PrimePay = primepay.NewSub
	All.EpayGlobalPay = epayglobal.NewSub
	All.Wheatpay = wheatpay.NewSub
	All.Lulupay = lulupay.NewSub
	All.Velopay = velopay.NewSub
	All.SyPay = sypay.NewSub
	All.CrownPay = crownpay.NewSub
	All.EastPay = eastpay.NewSub
	All.SyPay2 = sypay2.NewSub
	All.Penguinpay = penguinpay.NewSub
	All.Mlpay = mlpay.NewSub
	All.AIPay2 = aipay2.NewSub
	All.EPay178 = epay178.NewSub
	All.Deepay = deepay.NewSub
	All.DKPay = dkpay.NewSub
	All.GogoPayhx = gogopayhx.NewSub
	All.Megipay = megipay.NewSub
}

func NewSub(b *base.Base, index int) {
	if index <= int(values.InvalidPay) {
		return
	}
	if index >= int(values.PayAll) {
		log.Error("invalid index:%v", index)
		return
	}
	ele := reflect.ValueOf(All).Elem()
	if index >= ele.NumField() {
		log.Error("invalid index:%v", index)
		return
	}
	ref := ele.Field(index)
	if !ref.IsValid() || ref.IsZero() {
		log.Error("invalid index:%v", index)
		return
	}
	ref.Call([]reflect.Value{reflect.ValueOf(b)})
}
