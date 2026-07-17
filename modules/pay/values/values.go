package values

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"server/call"
	"server/common"
	"server/config"
	"server/db"
	"server/util"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liangdas/mqant/log"
	"gorm.io/gorm"
)

// 支付方式
type PayWay int

const (
	InvalidPay PayWay = iota
	TKPay
	ProPay
	WandaPay
	HongxinPay
	MQR
	RichPay
	WorthPay
	DeepLink
	VirgoPay
	VirgoPay2
	Toplemonpay
	PassPay
	MvPay
	RedPay
	WkPay
	MasterPay
	KingPay
	TransafePay
	HdPay
	NPay
	DecentPay
	ZoomPay
	SimplyPay
	AIPay
	SwayPay
	EPay
	GogoPay
	HexPay
	PrimePay
	EPayGolbalPay
	WheatPay
	LuluPay
	Velopay
	Sypay
	CrownPay
	EastPay
	Sypay2
	Penguinpay
	Mlpay
	AIPay2
	EPay178
	Deepay
	DKPay
	GogoPayhx
	Megipay
	PayAll
)

// 错误码
const (
	CodeOk = iota
)

var (
	PayWeightLock      = new(sync.RWMutex)
	WithdrawWeightLock = new(sync.RWMutex)
)

// SelectPayWay 特定条件为玩家优先选择渠道
func SelectPayWay(uid string, amount float64, banList []int) PayWay {
	if !config.GetConfig().Pay.SelectPayWay {
		return -1
	}
	zlist := db.Redis().ZRange(common.GetRedisKeyPlayerPay(uid))
	for _, v := range zlist {
		banList = append(banList, util.GetInt(v.Member))
	}
	lastSuccess := &common.OrderPay{UID: uid, Status: common.StatusROrderFinish}
	err := db.Mysql().GetLast(lastSuccess)
	if err == nil { // 直接按上次成功的渠道作为首选推荐渠道
		if !util.SliceContain(banList, lastSuccess.PayChannel) && IsTopChannel(lastSuccess.PayChannel, 3, amount) {
			return PayWay(lastSuccess.PayChannel)
		}
	}
	banAll := GetValidPayChannelsWithBanList(amount, banList, uid)
	if PayStatus == PayStatusChannel {
		for _, v := range banAll {
			return PayWay(v.ChannelID)
		}
	}
	// 尝试当前权重最高的
	all := GetValidPayChannels(amount)
	if len(all) > 0 {
		if !util.SliceContain(banList, all[0].ChannelID) {
			return PayWay(all[0].ChannelID)
		}
	}
	return -1
}

// 手动选择模式获取支付渠道
func GetPayWayChannel(index int, all []*common.ConfigPayChannels) (w PayWay) {
	if len(all) == 0 {
		w = -1
		return
	}
	if index > len(all)-1 || index < 0 {
		index = 0
	}
	w = PayWay(all[index].ChannelID)
	return
}

// 权重模式获取支付渠道
func GetPayWayWeight(all []*common.ConfigPayChannels) PayWay {
	if len(all) == 0 {
		return -1
	}
	index := util.RandomOneEleFromSlice(all, "Weight")
	if index >= 0 {
		return PayWay(all[index].ChannelID)
	}
	return -1
}

// ChoosePayWay 选择充值渠道
func ChoosePayWay(uid string, index int, amount float64, trys []int) PayWay {
	pw := SelectPayWay(uid, amount, trys)
	if pw >= 0 {
		return pw
	}
	banAll := GetValidPayChannelsWithBanList(amount, trys, uid)
	if PayStatus == PayStatusChannel {
		return GetPayWayChannel(index, banAll)
	}
	return GetPayWayWeight(banAll)
}

func PayFail(w PayWay) {
	failWeight := config.GetConfig().Pay.PayFailWeight
	if failWeight <= 0 {
		return
	}
	PayWeightLock.Lock()
	for _, v := range call.ConfigPayChannels {
		if v.ChannelID == int(w) {
			// 这里需要防止出现负值（出现负值可能出现代理没有渠道可用出现panic错误）
			if v.Weight > failWeight {
				v.Weight -= failWeight
			}
		}
	}
	sort.Slice(call.ConfigPayChannels, func(i, j int) bool {
		return call.ConfigPayChannels[i].Weight > call.ConfigPayChannels[j].Weight
	})
	PayWeightLock.Unlock()
	AddPayChannel(int(w), -failWeight)
}

func PayCallback(w PayWay) {
	successWeight := config.GetConfig().Pay.PaySuccessWeight
	if successWeight <= 0 {
		return
	}
	PayWeightLock.Lock()
	for _, v := range call.ConfigPayChannels {
		if v.ChannelID == int(w) {
			v.Weight += successWeight
			break
		}
	}
	sort.Slice(call.ConfigPayChannels, func(i, j int) bool {
		return call.ConfigPayChannels[i].Weight > call.ConfigPayChannels[j].Weight
	})
	PayWeightLock.Unlock()
	AddPayChannel(int(w), successWeight)
}

func WithdrawFails(w int) {
	call.Go(func() {
		db.Mysql().UpdateW(&common.ConfigWithdrawChannels{}, map[string]interface{}{"fails": gorm.Expr("fails + ?", 1)}, fmt.Sprintf("channel_id = %d", w))
	})
}

// ChooseWithdrawWay 选择代付渠道
// 参数:
//   - req: 提现请求信息
//   - trys: 已经尝试过的渠道ID列表
//   - mid: 商户ID
//
// 返回值:
//   - int: 选择的代付渠道ID，如果无可用渠道则返回-1
func ChooseWithdrawWay(req *common.WithdrawReq, trys []int, mid int) int {
	// 如果不是发布环境且配置了测试提现渠道，则优先使用测试渠道
	if !config.GetBase().Release && config.GetConfig().Pay.TestWithdraw >= 0 {
		// 检查测试渠道是否未在已尝试列表中
		if !util.SliceContain(trys, config.GetConfig().Pay.TestWithdraw) {
			return config.GetConfig().Pay.TestWithdraw
		}
		// 如果测试渠道已在尝试列表中，则返回-1表示不使用测试渠道
		return -1
	}

	// 获取商户配置信息
	var all []*common.ConfigWithdrawChannels
	mer := call.GetConfigMerchatsByID(mid)

	// 加读锁以安全访问全局渠道配置
	WithdrawWeightLock.RLock()
	if mer != nil && len(mer.PayoutChannelsSub) > 0 {
		// 如果商户有子渠道配置，则获取对应的渠道信息
		all = call.GetConfigWithdrawChannelsByIDs(mer.PayoutChannelsSub)
	} else {
		// 否则使用全局默认渠道配置
		all = call.ConfigWithdrawChannels
	}
	// 释放读锁
	WithdrawWeightLock.RUnlock()
	for _, v := range all {
		log.Info("获取代付渠道信息: %+v", *v)
	}

	// 遍历所有可用渠道，寻找合适的提现渠道
	for _, v := range all {
		// 如果该渠道已经尝试过，则跳过
		if util.SliceContain(trys, v.ChannelID) {
			continue
		}

		// 根据提现类型检查渠道是否支持
		if req.PayType == common.WithdrawTypeUPI && !v.UPI {
			continue
		}
		if req.PayType == common.WithdrawTypeBank && !v.Bank {
			continue
		}

		// 新增金额区间过滤
		// 如果权重小于等于0，或者金额不在渠道支持的范围内，则跳过
		if v.Weight <= 0 ||
			req.Amount < v.PayDown ||
			(v.PayUp > 0 && req.Amount > v.PayUp) {
			continue
		}

		// 记录选择的渠道信息
		log.Info("选择渠道:%d, name:%s, amount:%.2f, range:[%.2f-%.2f]",
			v.ChannelID, v.Name, req.Amount, v.PayDown, v.PayUp)

		// 返回选中的渠道ID
		return v.ChannelID
	}

	// 如果没有找到合适的渠道，则记录日志并返回-1
	log.Info("没有可用的代付渠道")
	return -1
}

// 获取支付回调路径
func GetPayCallback(way PayWay) string {
	return config.GetConfig().Pay.CallbackURL + "/cb/pi/" + fmt.Sprintf("%d", way)
}

// 获取退出回调路径
func GetWithdrawCallback(way PayWay) string {
	return config.GetConfig().Pay.CallbackURL + "/cb/po/" + fmt.Sprintf("%d", way)
}

func GetFrontCallback() string {
	return config.GetConfig().Pay.CallbackURL + "/common/front/payCallback"
}

func GetFrontCallbackFail() string {
	return config.GetConfig().Pay.CallbackURL + "/common/front/payCallback/fail"
}

// 根据body里的字段直接拼接出签名字符串
func GetSignStr(str string, pass ...string) string {
	m := map[string]json.RawMessage{}
	sortStrs := []string{}
	json.Unmarshal([]byte(str), &m)
	for i := range m {
		sortStrs = append(sortStrs, i)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		if len(m[v]) > 1 && m[v][0] == 34 {
			m[v] = m[v][1 : len(m[v])-1]
		}
		if len(m[v]) == 0 {
			continue
		}
		signStr += fmt.Sprintf("%v=%v", v, string(m[v]))
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

// 根据body里的字段直接拼接出签名字符串(去除null)
func GetSignStrNull(str string, pass ...string) string {
	m := map[string]json.RawMessage{}
	sortStrs := []string{}
	json.Unmarshal([]byte(str), &m)
	for i := range m {
		sortStrs = append(sortStrs, i)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		if len(m[v]) > 1 && m[v][0] == 34 {
			m[v] = m[v][1 : len(m[v])-1]
		}
		if len(m[v]) == 0 || string(m[v]) == "null" {
			continue
		}
		signStr += fmt.Sprintf("%v=%v", v, string(m[v]))
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

// 根据body里的字段直接拼接出签名字符串(去除null，但空参与验签)
func GetSignStrNullWithEmpty(str string, pass ...string) string {
	m := map[string]json.RawMessage{}
	sortStrs := []string{}
	json.Unmarshal([]byte(str), &m)
	for i := range m {
		sortStrs = append(sortStrs, i)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		if len(m[v]) > 1 && m[v][0] == 34 {
			m[v] = m[v][1 : len(m[v])-1]
		}
		if string(m[v]) == "null" {
			continue
		}
		signStr += fmt.Sprintf("%v=%v", v, string(m[v]))
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

// 根据body里的字段直接拼接出签名字符串（包括空字符串）
func GetSignStrAll(str string, pass ...string) string {
	m := map[string]json.RawMessage{}
	sortStrs := []string{}
	json.Unmarshal([]byte(str), &m)
	for i := range m {
		sortStrs = append(sortStrs, i)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		if len(m[v]) > 1 && m[v][0] == 34 {
			m[v] = m[v][1 : len(m[v])-1]
		}
		signStr += fmt.Sprintf("%v=%v", v, string(m[v]))
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

// 根据body里的字段直接拼接出签名字符串（包括空字符串,unicode转码）
func GetSignStrAllUnicode(str string, pass ...string) string {
	m := map[string]json.RawMessage{}
	sortStrs := []string{}
	json.Unmarshal([]byte(str), &m)
	for i := range m {
		sortStrs = append(sortStrs, i)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		str := string(m[v])
		if len(str) > 1 && str[0] == 34 {
			str = str[1 : len(str)-1]
		}
		tmp, err := util.ZhToUnicode(str)
		if err == nil {
			str = tmp
		}
		str = strings.ReplaceAll(str, `\`, "")
		signStr += fmt.Sprintf("%v=%v", v, str)
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

// 根据body里的字段直接拼接出签名字符串
func GetSignStrURLEncode(str string, pass ...string) string {
	sortStrs := []string{}
	all, err := url.ParseQuery(str)
	if err != nil {
		log.Error("err:%e", err)
		return ""
	}
	for k := range all {
		sortStrs = append(sortStrs, k)
	}

	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		tmp := all.Get(v)
		if len(tmp) > 1 && tmp[0] == 34 {
			tmp = tmp[1 : len(tmp)-1]
		}
		if len(tmp) == 0 {
			continue
		}
		signStr += fmt.Sprintf("%v=%v", v, string(tmp))
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

func GetSignStrGet(vals url.Values, pass ...string) string {
	sortStrs := []string{}
	for k := range vals {
		sortStrs = append(sortStrs, k)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		if len(vals.Get(v)) == 0 {
			continue
		}
		signStr += fmt.Sprintf("%v=%v", v, vals.Get(v))
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

func GetSignStrForm(c *gin.Context, pass ...string) string {
	sortStrs := []string{}
	for k := range c.Request.PostForm {
		sortStrs = append(sortStrs, k)
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		if len(c.PostForm(v)) == 0 {
			continue
		}
		signStr += fmt.Sprintf("%v=%v", v, c.PostForm(v))
		signStr += "&"
	}
	if len(signStr) > 0 {
		signStr = signStr[:len(signStr)-1]
		log.Debug("signStr:%v", signStr)
	}
	return signStr
}

func GetSignStrFormURLEncode(c *gin.Context, model interface{}, pass ...string) string {
	sortStrs := []string{}
	for k := range c.Request.URL.Query() {
		sortStrs = append(sortStrs, k)
	}
	// ref := reflect.ValueOf(model)
	reft := reflect.TypeOf(model)
	if reft.Kind() == reflect.Ptr {
		// ref = ref.Elem()
		reft = reft.Elem()
	}
	signStr := ""
	sort.Strings(sortStrs)
	for _, v := range sortStrs {
		shouldPass := false
		for _, s := range pass {
			if v == s {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		param := c.Query(v)
		if len(param) == 0 {
			continue
		}
		first := v[0]
		tmpName := strings.ToUpper(string(first)) + v[1:]
		field, ok := reft.FieldByName(tmpName)
		if ok && field.Tag.Get("encode") == "1" { // 需要urlencode
			param, _ = url.QueryUnescape(param)
		}

		signStr += fmt.Sprintf("%v=%v", v, param)
		signStr += "&"
	}
	fmt.Println(signStr)
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

// status=1时，属于不知道第三方是什么情况的报错，不能当成业务中断
// status=0时，可以根据返回判断是否是业务错误
func Post(req *http.Request, ret interface{}) int {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req.Close = true
	log.Debug("req:%+v", req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("http post call err:%v", err)
		return 1
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Debug("response Body:%v", string(body))
	if err := json.Unmarshal(body, ret); err != nil {
		log.Error("unmarshal fail err:%v", err)
		return 1
	}
	return 0
}

func AddPayChannel(p int, amount int) {
	if err := db.Mysql().C().Model(&common.ConfigPayChannels{}).Where("channel_id = ?", p).Updates(map[string]interface{}{"weight": gorm.Expr("weight + ?", amount)}).Error; err != nil {
		log.Error("err:%v", err)
	}
}

func SetPayChannel(p int, amount int) {
	if err := db.Mysql().C().Model(&common.ConfigPayChannels{}).Where("channel_id = ?", p).Updates(map[string]interface{}{"weight": amount}).Error; err != nil {
		log.Error("err:%v", err)
	}
}

func GetValidPayChannels(amount float64) []*common.ConfigPayChannels {
	PayWeightLock.RLock()
	defer PayWeightLock.RUnlock()
	all := []*common.ConfigPayChannels{}
	for _, v := range call.ConfigPayChannels {
		if v.Weight > 0 && amount >= v.PayDown && (amount <= v.PayUp || v.PayUp <= 0) {
			all = append(all, v)
		}
	}
	return all
}

/**
 * 获取有效的支付渠道列表，排除禁用列表中的渠道
 * @param amount 交易金额
 * @param banList 禁用的支付渠道ID列表
 * @param uid 用户ID
 * @return 返回有效的支付渠道列表
 */
func GetValidPayChannelsWithBanList(amount float64, banList []int, uid string) []*common.ConfigPayChannels {
	// 初始化一个空的支付渠道切片
	all := []*common.ConfigPayChannels{}
	// 获取支付权重锁的读锁，确保在读取过程中数据不会被修改
	PayWeightLock.RLock()
	// 在函数返回前释放读锁
	defer PayWeightLock.RUnlock()
	// 遍历所有配置的支付渠道
	for _, v := range call.ConfigPayChannels {
		// 检查渠道权重是否大于0，且交易金额在渠道的支付范围内
		if v.Weight > 0 && amount >= v.PayDown && (amount <= v.PayUp || v.PayUp <= 0) {
			// 检查渠道ID是否不在禁用列表中
			if !util.SliceContain(banList, v.ChannelID) {
				// 将符合条件的渠道添加到结果列表中
				all = append(all, v)
			}
		}
	}
	// 返回有效的支付渠道列表
	return all
}

func IsPayChannelValid(cid int, amount float64) bool {
	PayWeightLock.RLock()
	defer PayWeightLock.RUnlock()
	for _, v := range call.ConfigPayChannels {
		if cid == v.ChannelID {
			if v.Weight > 0 && amount >= v.PayDown && (amount <= v.PayUp || v.PayUp <= 0) {
				return true
			}
			break
		}
	}
	return false
}

// IsTopChannel 判断渠道是否是排名前几的可用渠道
func IsTopChannel(cid, top int, amount float64) bool {
	PayWeightLock.RLock()
	defer PayWeightLock.RUnlock()
	for i := 0; i < top; i++ {
		v := call.ConfigPayChannels[i]
		if cid == v.ChannelID {
			if v.Weight > 0 && amount >= v.PayDown && (amount <= v.PayUp || v.PayUp <= 0) {
				return true
			}
		}
	}
	return false
}
