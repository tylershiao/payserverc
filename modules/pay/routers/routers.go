package routers

import (
	"reflect"
	"runtime"
	"server/common"
	"server/config"
	"server/modules/pay/allpay"
	"server/modules/pay/base"
	"server/modules/pay/middleware"
	"server/modules/pay/values"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
)

func SetUpRouter() *gin.Engine {
	gin.EnableJsonDecoderUseNumber()
	if config.GetBase().Release {
		gin.SetMode(gin.ReleaseMode)
		// 禁用控制台颜色
		gin.DisableConsoleColor()
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	// 跨域处理
	r.Use(middleware.CrosHandler())
	r.Use(middleware.RateLimitByIPAndAPI())
	r.Use(gin.RecoveryWithWriter(log.LogBeego(), func(c *gin.Context, err interface{}) {
		buf := make([]byte, 2048)
		runtime.Stack(buf, false)
		log.Error("panic(%+v), stack:\n%s", err, string(buf))
	}))
	r.GET("/", common.HealthCheck)
	com := r.Group("common")
	commonRouter(com)
	page := r.Group("page")
	pageRouter(page)
	// sys := r.Group("sys")
	// sysRouter(sys)
	api := r.Group("api")
	apiv1 := api.Group("v1")
	apiV1Router(apiv1)
	// 统一回调
	cb := r.Group("cb")
	Callback(cb)

	return r
}

func Callback(e *gin.RouterGroup) {
	e.POST("pi/*action", PayCallback)
	e.GET("pi/*action", PayCallback)
	e.POST("po/*action", WithdrawCallback)
	e.GET("po/*action", WithdrawCallback)
}

func PayCallback(c *gin.Context) {
	log.Debug("PayCallback:%+v", c.Request)
	path := c.Request.URL.Path
	path = strings.ReplaceAll(path, "/cb/pi/", "")
	index, err := strconv.Atoi(path)
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	b := base.NewCallbackBase(base.OPTPayCB)
	b.Channel = values.PayWay(index)
	allpay.NewSub(b, index)
	if b.Sub == nil {
		return
	}
	b.PayCallback(c)
}

func WithdrawCallback(c *gin.Context) {
	log.Debug("WithdrawCallback:%+v", c.Request)
	path := c.Request.URL.Path
	path = strings.ReplaceAll(path, "/cb/po/", "")
	index, err := strconv.Atoi(path)
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	ref := reflect.ValueOf(allpay.All).Elem().Field(index)
	if !ref.IsValid() {
		log.Error("invalid index:%v", index)
		return
	}
	b := base.NewCallbackBase(base.OPTWithdrawCB)
	b.Channel = values.PayWay(index)
	ref.Call([]reflect.Value{reflect.ValueOf(b)})
	b.WithdrawCallback(c)
}

func SetUpSysRouter() *gin.Engine {
	gin.EnableJsonDecoderUseNumber()
	r := gin.New()

	// 跨域处理
	r.Use(middleware.CrosHandler())
	// r.Use(middleware.RateLimitByIPAndAPI())
	r.Use(gin.RecoveryWithWriter(log.LogBeego(), func(c *gin.Context, err interface{}) {
		buf := make([]byte, 2048)
		runtime.Stack(buf, false)
		log.Error("panic(%+v), stack:\n%s", err, string(buf))
	}))
	r.GET("/", common.HealthCheck)
	// com := r.Group("common")
	// commonRouter(com)
	sys := r.Group("sys")
	sysRouter(sys)
	// api := r.Group("api")
	// apiv1 := api.Group("v1")
	// apiV1Router(apiv1)
	// // 统一回调
	// cb := r.Group("callback")
	// Callback(cb)

	return r
}
