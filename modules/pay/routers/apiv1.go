package routers

import (
	"server/modules/pay/apiv1"

	"github.com/gin-gonic/gin"
)

func apiInV1Router(e *gin.RouterGroup) {
	e.POST("ds/ld", apiv1.PayInCreate)
	e.POST("ds/cd", apiv1.PayInQuery)
	e.POST("df/ld", apiv1.PayOutCreate)
	e.POST("df/cd", apiv1.PayOutQuery)
	e.POST("sh/ye", apiv1.Balance)
}
