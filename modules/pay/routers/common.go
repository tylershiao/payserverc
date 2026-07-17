package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
)

func commonRouter(e *gin.RouterGroup) {
	e.GET("front/payCallback", frontPayCallback)
	e.GET("front/payCallback/fail", frontPayCallbackFail)
}

func frontPayCallback(c *gin.Context) {
	log.Debug("frontcallback:%v", c.Request.RequestURI)
	c.Redirect(http.StatusMovedPermanently, "zypayment://result=true")
}

func frontPayCallbackFail(c *gin.Context) {
	log.Debug("frontcallback:%v", c.Request.RequestURI)
	c.Redirect(http.StatusMovedPermanently, "zypayment://result=false")
}
