package common

import (
	"fmt"
	"time"
)

const (
	RedisExpireLock  = 10 * time.Second
	RedisExpireToken = 7 * 24 * time.Hour
)

const (
	RedisTimeoutTx = 10 * time.Second
)

const (
	RedisKeyUser                  = "user"
	RedisKeyPlayerPay             = "playerPay"             // 玩家10分钟内付费拉起的渠道记录
	RedisKeyPayStatus             = "payStatus"             // 支付渠道策略
	RedisKeyWithdrawRiskCardCount = "withdrawRiskCardCount" // 代付卡号日请求计数
	RedisKeyWithdrawRiskCardPass  = "withdrawRiskCardPass"  // 代付卡号当日审核放行
)

const (
	RedisLockKeyIP = "lockKeyIP" // 玩家操作频率限制
)

func GetRedisLockKeyIP(ip string) string {
	return fmt.Sprintf("%v:%v", RedisLockKeyIP, ip)
}

func GetRedisKeyUser(uid string) string {
	return fmt.Sprintf("%s:%s", RedisKeyUser, uid)
}

func GetRedisKeyPlayerPay(uid string) string {
	return fmt.Sprintf("%s:%s", RedisKeyPlayerPay, uid)
}

func GetRedisKeyWithdrawRiskCardCount(date, account string) string {
	return fmt.Sprintf("%s:%s:%s", RedisKeyWithdrawRiskCardCount, date, account)
}

func GetRedisKeyWithdrawRiskCardPass(date, account string) string {
	return fmt.Sprintf("%s:%s:%s", RedisKeyWithdrawRiskCardPass, date, account)
}
