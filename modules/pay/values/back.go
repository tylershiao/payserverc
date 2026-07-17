package values

// 报错后通知后台机器人发送错误到群里
type BackReq struct {
	Opt      int
	Channel  int
	ErrorMsg string
}
