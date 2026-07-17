package natsClient

const (
	TopicNewPlayer        = "TopicNewPlayer"
	TopicClientDisconnect = "TopicClientDisconnect"
	TopicReloadConfig     = "TopicReloadConfig" // 重置配置
	TopicBroadcast        = "TopicBroadcast"    // 广播
	TopicBroadcastReq     = "TopicBroadcastReq" // 广播
	TopicCancelMatch      = "TopicCancelMatch"
	TopicInnerRefreshGold = "TopicInnerRefreshGold" // 通知游戏服刷新金币
	// TopicInnerStatistics     = "TopicInnerStatistics"     // 通知common服写入玩家统计数据
	// TopicInnerOnline = "TopicInnerOnline" // 通知common服写入玩家在线统计数据
	// TopicInnerPushRed        = "TopicInnerPushRed"        // 通知common服推送玩家红点
	// TopicInnerUpdateRed = "TopicInnerUpdateRed" // 通知common服更新并推送玩家红点
	// TopicInnerESPlayerData   = "TopicInnerESPlayerData"   // 通知common服更新玩家数据
	// TopicInnerWarn           = "TopicInnerWarn"           // 通知common服预警
	TopicInnerUpdateWarn     = "TopicInnerUpdateWarn"     // 通知common服更新预警参数
	TopicInnerOptPlayer      = "TopicInnerOptPlayer"      // 通知全服对玩家进行一些操作
	TopicInnerPlayerWithdraw = "TopicInnerPlayerWithdraw" // 通知游戏服某个玩家提现回调了
	TopicInnerAfterSettle    = "TopicInnerAfterSettle"    // 玩家游戏场结算后的一些处理逻辑
	TopicBroadcastAll        = "TopicBroadcastAll"        // 任意协议进行广播操作

	TopicInnerESBulk = "TopicInnerESBulk" // es批量写入
)

const (
	TopicBackRefreshMail = "TopicBackRefreshMail" // 后台通知游戏服刷新玩家邮件
)
