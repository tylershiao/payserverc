package values

import (
	"time"
)

const (
	PayTimerInterval = 5 * time.Minute
)

var (
	PayStatus = PayStatusChannel // 当前支付选择方式 1是按渠道号选择 2按渠道权重动态分配
	// CheckCount  = 0                // 按照动态权重配置的检查次数
	// CheckReset  = 3                // 检查了预定次数还没有找到合适渠道，重置渠道
	RootChannel = []int{} // 常驻渠道,在渠道重置时，默认加上
)

// 1是按渠道号选择 2按渠道权重动态分配
const (
	PayStatusChannel = iota + 1
	PayStatusAuto
	PayStatusForceAuto // 强制自动模式,此模式下不再切换
)
