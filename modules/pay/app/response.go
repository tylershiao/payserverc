package app

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"server/call"
	"server/common"
	"server/config"
	"server/db"
	"server/sign"
	"server/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
	"gorm.io/gorm"
)

type Gin struct {
	R
	Context *gin.Context
	Merchat *common.Merchats
}

type R struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
}

func NewApp(c *gin.Context) *Gin {
	g := &Gin{Context: c, R: R{Code: CodeOK}}
	return g
}

// Response setting gin.JSON
func (g *Gin) Response() {
	if g.Msg == "" {
		switch g.Code {
		case CodeOK:
			g.Msg = "Success"
		case CodeInnerError:
			g.Msg = "Internal service error"
		case CodeParam:
			if g.Msg == "" {
				g.Msg = "Invalid request parameters"
			}
		case CodeServer:
			g.Msg = "Service temporarily unavailable"
		case CodeSign:
			g.Msg = "Signature verification failed"
		case CodeOrderNotExist:
			g.Msg = "Order not found"
		case CodeDuplicateOrder:
			g.Msg = "Merchant order already exists"
		case CodeRisk:
			g.Msg = "Transaction rejected by risk control"
		case CodeInvalidMerchat:
			g.Msg = "Merchant not found"
		case CodeMerchatClose:
			g.Msg = "Merchant service unavailable"
		case CodeInsufficientBalance:
			g.Msg = "Merchant balance is insufficient"
		case CodeIP:
			g.Msg = "Source IP is not allowed"
		case CodeNoChannel:
			g.Msg = "No payment channel is currently available"
		case CodeTooManyRequests:
			g.Msg = "Request rate limit exceeded"
		}
	}
	g.Context.JSON(http.StatusOK, g.R)
}

// S 解析请求
func (g *Gin) S(one interface{}) (pass bool) {
	// 读取 body
	bodyBytes, _ := io.ReadAll(g.Context.Request.Body)
	raw := string(bodyBytes)
	// 恢复 body 让 ShouldBind 能继续用
	g.Context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	err := g.Context.ShouldBind(one)
	if err != nil {
		g.Code = CodeParam
		log.Error("bind %v err:%v,reqbody:%s", reflect.TypeOf(one), err, raw)
		return
	}
	pass = true
	return
}

// 从数据库中读取数据
func (g *Gin) MGet(tar interface{}) (pass bool) {
	err := db.Mysql().Get(tar)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("err:%v", err)
		g.R.Code = CodeInnerError
		g.R.Msg = err.Error()
		return
	}
	pass = true
	return
}

// MCommit 提交事务
func (g *Gin) MCommit(tx *gorm.DB) {
	if g.Code == CodeOK {
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return
		}
	} else {
		tx.Rollback()
		return
	}
}

func (g *Gin) GetRemoteIP() string {
	return util.RequestIP(g.Context.Request)
}

func (g *Gin) CheckSign(send interface{}, key, signStr string) (pass bool) {
	str := sign.APISignMD5(send, key, "Sign", "signature")
	pass = strings.EqualFold(str, signStr)
	if !pass {
		log.Debug("check sign fail:%v", send)
		g.Code = CodeSign
	}
	return
}

// 风控
func (g *Gin) Risk(data *common.CommonPayData) (pass bool) {
	pass = true
	if !pass {
		g.Code = CodeRisk
	}
	return
}

// 检查商户是否合法
func (g *Gin) CheckMerchat(req interface{}, mid int, sign string, opt int) (pass bool) {
	merchat := g.Merchat
	if merchat == nil {
		merchat = call.GetConfigMerchatsByID(mid)
		g.Merchat = merchat
	}
	if merchat == nil {
		g.Code = CodeInvalidMerchat
		return
	}
	if !merchat.Open {
		g.Code = CodeMerchatClose
		return
	}
	switch opt {
	case 1:
		if config.GetBase().Release {
			ip := g.GetRemoteIP()
			if merchat.WhiteIPSub == nil || !util.SliceContainIP(merchat.WhiteIPSub, ip) {
				g.Code = CodeIP
				log.Debug("ip:%s invalid,merchat.WhiteIPSub:%v", ip, merchat.WhiteIPSub)
				return
			}
		}
		if !merchat.PaySwitch {
			g.Code = CodeMerchatClose
			return
		}
	case 3:
		if config.GetBase().Release {
			ip := g.GetRemoteIP()
			if merchat.WhiteIPSub == nil || !util.SliceContainIP(merchat.WhiteIPSub, ip) {
				g.Code = CodeIP
				log.Debug("ip:%s invalid,merchat.WhiteIPSub:%v", ip, merchat.WhiteIPSub)
				return
			}
		}
		if !merchat.WithdrawSwitch {
			g.Code = CodeMerchatClose
			return
		}
	}
	if !g.CheckSign(req, merchat.Key, sign) {
		return
	}
	pass = true
	return
}
