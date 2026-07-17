package middleware

import (
	"fmt"
	"server/config"
	"server/modules/pay/app"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
	"golang.org/x/time/rate"
)

func InitApiRateLimits() {
	orderRate := config.GetConfig().Pay.OrderRate
	if orderRate == 0 {
		orderRate = 10
	}
	apiRateLimits["/api/pi/create"] = OneRate{Rate: rate.Limit(orderRate), BurstSize: orderRate}
	apiRateLimits["/api/po/create"] = OneRate{Rate: rate.Limit(orderRate), BurstSize: orderRate}
	queryRate := config.GetConfig().Pay.QueryRate
	if queryRate == 0 {
		queryRate = 10
	}
	apiRateLimits["/api/pi/query"] = OneRate{Rate: rate.Limit(queryRate), BurstSize: queryRate}
	apiRateLimits["/api/po/query"] = OneRate{Rate: rate.Limit(queryRate), BurstSize: queryRate}
	log.Debug("InitApiRateLimits,orderRate:%d,queryRate:%d", orderRate, queryRate)
}

type OneRate struct {
	Rate      rate.Limit
	BurstSize int
}

// 每个 API 路径的限流配置（可根据实际调整）
var apiRateLimits = map[string]OneRate{
	"/api/payIn/create":  {Rate: 10, BurstSize: 10}, // 每秒 10 次
	"/api/payIn/query":   {Rate: 20, BurstSize: 20}, // 每秒 20 次
	"/api/payOut/create": {Rate: 10, BurstSize: 10}, // 每秒 10 次
	"/api/payOut/query":  {Rate: 20, BurstSize: 20}, // 每秒 20 次
	"/api/balance":       {Rate: 10, BurstSize: 10}, // 每秒 10 次
}

// 使用 sync.Map 替代 map + mutex
var (
	limiters    sync.Map // key: ip|api, value: *rate.Limiter
	defaultRate = OneRate{Rate: 10, BurstSize: 10}
)

// 获取或创建限流器
func getLimiter(ip, api string) *rate.Limiter {
	log.Debug("api:%v", api)
	key := fmt.Sprintf("%s|%s", ip, api)

	// 快速读取
	if v, ok := limiters.Load(key); ok {
		return v.(*rate.Limiter)
	}

	// 创建新 limiter，加载或存储
	cfg, ok := apiRateLimits[api]
	if !ok {
		cfg = defaultRate // 默认限流策略
	}
	limiter := rate.NewLimiter(cfg.Rate, cfg.BurstSize)

	// 双重检查，防止并发创建
	actual, _ := limiters.LoadOrStore(key, limiter)
	return actual.(*rate.Limiter)
}

// Gin 限流中间件（基于 IP + API 路径）
func RateLimitByIPAndAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := app.NewApp(c)
		ip := a.GetRemoteIP()
		api := c.FullPath() // 获取定义路由，如 "/api/pay"
		if api == "" {
			// 匹配不到路由，跳过限流
			c.Next()
			return
		}

		api = normalizePath(api)

		if api != "/cb/pi/*action" && api != "/cb/po/*action" { // 回调不做限制
			limiter := getLimiter(ip, api)
			log.Debug("limiter:%v", limiter)
			if !limiter.Allow() {
				a.Code = app.CodeTooManyRequests
				a.Response()
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// 版本号前缀列表，后续可拓展
var versionPrefixes = []string{"/v1/", "/v2/", "/v3/"}

// 归一化路径，去除版本号部分
func normalizePath(path string) string {
	for _, v := range versionPrefixes {
		// 举例："/api/v1/payIn/create" 去掉 "/v1"
		// 变成 "/api/payIn/create"
		if strings.Contains(path, v) {
			path = strings.ReplaceAll(path, v, "/")
		}
	}
	return path
}
