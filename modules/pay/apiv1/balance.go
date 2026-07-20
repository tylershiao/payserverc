package apiv1

import (
	"server/common"
	"server/db"
	"server/modules/pay/app"

	"github.com/gin-gonic/gin"
)

type BalanceReq struct {
	MerchantID int    `json:"merchantId" binding:"required"`
	Signature  string `json:"signature" binding:"required"`
}

type BalanceResp struct {
	MerchantID     int   `json:"merchantId"`     // 商户号
	AvailableFunds int64 `json:"availableFunds"` // 可用余额
	PendingFunds   int64 `json:"pendingFunds"`   // 待结算余额
	LockedFunds    int64 `json:"lockedFunds"`    // 冻结金额
}

// Balance 函数用于获取商户余额信息
// 该函数处理商户余额查询请求，返回商户的可用余额、结算余额和冻结余额
func Balance(c *gin.Context) {
	a := app.NewApp(c)
	defer func() {
		a.Response()
	}()
	req := new(BalanceReq)
	if !a.S(req) {
		return
	}
	mer := &common.Merchats{MID: req.MerchantID}
	err := db.Mysql().Get(mer)
	if err != nil {
		return
	}
	if mer.ID > 0 {
		a.Merchat = mer
	}
	if !a.CheckMerchat(req, req.MerchantID, req.Signature, 0) {
		return
	}
	a.Data = &BalanceResp{
		MerchantID:     req.MerchantID,
		AvailableFunds: int64(mer.Balance * 100),
		PendingFunds:   int64(mer.BalanceSettle * 100),
		LockedFunds:    int64(mer.BalanceFreeze * 100),
	}
}
