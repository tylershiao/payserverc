package apiv1

import (
	"server/common"
	"server/db"
	"server/modules/pay/app"

	"github.com/gin-gonic/gin"
)

type BalanceReq struct {
	MID  int    `json:"MID" binding:"required"`
	Sign string `json:"Sign" binding:"required"`
}

type BalanceResp struct {
	MID               int   // 商户号
	AvailableBalance  int64 // 可用余额
	SettlementBalance int64 // 待结算余额
	FrozenBalance     int64 // 冻结金额
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
	mer := &common.Merchats{MID: req.MID}
	err := db.Mysql().Get(mer)
	if err != nil {
		return
	}
	if mer.ID > 0 {
		a.Merchat = mer
	}
	if !a.CheckMerchat(req, req.MID, req.Sign, 0) {
		return
	}
	a.Data = &BalanceResp{
		MID:               req.MID,
		AvailableBalance:  int64(mer.Balance * 100),
		SettlementBalance: int64(mer.BalanceSettle * 100),
		FrozenBalance:     int64(mer.BalanceFreeze * 100),
	}
}
