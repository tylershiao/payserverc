package routers

import (
	"server/modules/pay/apiv1"

	"github.com/gin-gonic/gin"
)

func apiV1Router(e *gin.RouterGroup) {
	e.POST("pi/create", apiv1.PayInCreate)
	e.POST("pi/query", apiv1.PayInQuery)
	e.POST("po/create", apiv1.PayOutCreate)
	e.POST("po/query", apiv1.PayOutQuery)
	e.POST("mer/balance", apiv1.Balance)
}
