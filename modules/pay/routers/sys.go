package routers

import (
	"server/modules/pay/sys"

	"github.com/gin-gonic/gin"
)

func sysRouter(e *gin.RouterGroup) {
	e.POST("reload", sys.Reload)
	e.POST("testPay", sys.TestPay)
}
