package util

import (
	"bufio"
	"container/list"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/liangdas/mqant/log"
)

var privateIPRanges = []netip.Prefix{
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
}

// FormatFloat 取小数点后n位小数，四舍五入
func FormatFloat(num float64, decimal int) string {
	if math.Trunc(num) == num || decimal == 0 {
		return fmt.Sprintf("%.f", math.Trunc(num))
	}
	format := "%." + strconv.Itoa(decimal) + "f"
	return fmt.Sprintf(format, num)
}

// RoundFloat 取小数点后n位非零小数
func RoundFloat(num float64, decimal int) string {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}

// Decimal 保留x位小数，返回float64
func Decimal(value float64, deci int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(deci)+"f", value), 64)
	return value
}

func GetInt(data interface{}) int {
	switch data := data.(type) {
	case int:
		return data
	case int64:
		return int(data)
	case float32:
		return int(data)
	case float64:
		return int(data)
	case string:
		ret, _ := strconv.Atoi(data)
		return ret
	default:
		return 0
	}
}

func GetInt64(data interface{}) int64 {
	switch data := data.(type) {
	case int:
		return int64(data)
	case int64:
		return data
	case float32:
		return int64(data)
	case float64:
		return int64(data)
	case string:
		ret, _ := strconv.ParseInt(data, 10, 64)
		return ret
	default:
		return 0
	}
}

// CheckFileIsExist 判断文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// RandomLetters 随机26位不重复的字符串
func RandomLetters() string {
	a := rand.Perm(26)
	s := []byte{}
	for _, v := range a {
		s = append(s, byte(97+v))
	}
	return string(s)
}

var (
	characters52 = [52]byte{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}
)

func GenerateRandomString(size int) string {
	raw := make([]byte, size)
	for i := 0; i < size; i++ {
		raw[i] = characters52[rand.Intn(52)]
	}
	return string(raw)
}

// StrFirstToUpper 首字母大写
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []byte(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}

/*
判断文件或文件夹是否存在
如果返回的错误为nil,说明文件或文件夹存在
如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
如果返回的错误为其它类型,则不确定是否在存在
*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CheckPath 判断目录是否存在，不存在自动创建
func CheckPath(path string) error {
	b, err := PathExists(path)
	if err != nil {
		return err
	}
	if b {
		return fmt.Errorf("path %s 存在", path)
	}
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// SavePic 下载并保存图片
func SavePic(imgPath, imgUrl, fileName string) error {
	res, err := http.Get(imgUrl)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	path := imgPath + string(os.PathSeparator) + fileName + ".png"
	if ok, _ := PathExists(path); ok {
		if err := os.Remove(path); err != nil {
			log.Error("err:%v", err)
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	defer file.Close()
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	if _, err := io.Copy(writer, reader); err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// GetSimpleRandomString 获取纯英文大写的随机n位字符串
func GetSimpleRandomString(n int) string {
	ret := []byte{}
	for i := 0; i < n; i++ {
		tmp := rand.Intn(26)
		ret = append(ret, byte(tmp+65))
	}
	return string(ret)
}

func TrimHiddenCharacter(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 {
			continue
		}
		if c == 127 {
			continue
		}
		dstRunes = append(dstRunes, c)
	}
	return string(dstRunes)
}

// GetRandSection 获取区间随机值
func GetRandSection(a, b int) int {
	if a >= b {
		return a
	}
	return a + rand.Intn(b-a)
}

func GetRandSectionRange(a []int64) (ret int64) {
	if len(a) != 2 {
		return
	}
	diff := a[1] - a[0]
	if diff == 0 {
		return a[0]
	}
	if diff < 0 {
		diff = -diff
	}
	return a[0] + rand.Int63n(diff)
}

// 字符串去重
func RemoveDuplication(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

// Abs 取绝对值
func Abs(a int64) int64 {
	if a < 0 {
		a = -a
	}
	return a
}

func CleanList(l *list.List) {
	e := l.Front()
	if e == nil {
		return
	}
	for {
		next := e.Next()
		l.Remove(e)
		if next == nil {
			break
		}
		e = next
	}
}

func CheckPhone(phone string) string {
	if len(phone) != 10 || (phone[0] != '9' && phone[0] != '8' && phone[0] != '7' && phone[0] != '6') {
		ra := rand.Intn(90000000) + 9110000000
		phone = strconv.Itoa(ra)
	}
	return phone
}

// 获取下一个整点5分钟时间戳
func GetNext5MinUnix() int64 {
	now := time.Now()
	zero := GetZeroTime(now)
	h, m, _ := now.Clock()
	m -= m % 5
	last := zero.Add(time.Duration(h) * time.Hour).Add(time.Duration(m) * time.Minute)
	return last.Add(5 * time.Minute).Unix()
}

// 获取下一个整点时间戳
func GetNextHour() int64 {
	now := time.Now()
	zero := GetZeroTime(now)
	h, _, _ := now.Clock()
	return zero.Add(time.Duration(h+1) * time.Hour).Unix()
}

// stringsToInt 数字字符串转数字数组
func StringToInt(s string) []int {
	b := []byte(s)
	ret := []int{}
	for _, v := range b {
		if v < 48 || v > 57 {
			return nil
		}
		ret = append(ret, int(v)-48)
	}
	return ret
}

func ZhToUnicode(raw string) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}

// IsSingle 判断数字是不是单数
func IsSingle(n int) bool {
	return n/2*2 != n
}

func ValidateEmail(email string) bool {
	// 正则表达式来匹配电子邮箱地址
	// 基本电子邮箱格式是: 用户名@域名
	// 这里的正则表达式要求用户名可以包含字母、数字、点、下划线和减号，域名则要求包含至少一个点
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func GetNumberFromStr(str string) int {

	// 正则表达式模式
	re := regexp.MustCompile(`\d+`)

	// 查找所有匹配的数字
	numbers := re.FindAllString(str, -1)

	if len(numbers) > 0 {
		ret, err := strconv.Atoi(numbers[0])
		if err != nil {
			log.Error("err:%v", err)
			return 0
		}
		return ret
	}
	return 0
}

// Round100 按100取整
func Round100(tmp int64) int64 {
	tmp2 := tmp / 100 * 100
	if tmp2 < tmp {
		return tmp2 + 100
	}
	return tmp2
}

// Go 用协程处理f
func Go(f func()) {
	go func() {
		defer Recover()
		f()
	}()
}

func IsValidIndianLocalMobile(mobile string) bool {
	// 精确匹配 10 位，以 6-9 开头
	pattern := `^[6-9]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(mobile)
}

func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func NormalizeIP(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(raw); err == nil {
		raw = host
	}
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	parsed := net.ParseIP(raw)
	if parsed == nil {
		return ""
	}
	return parsed.String()
}

func IsPrivateIP(ip string) bool {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	for _, prefix := range privateIPRanges {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

func RequestIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	remoteIP := NormalizeIP(r.RemoteAddr)
	if remoteIP == "" {
		return ""
	}
	if !IsPrivateIP(remoteIP) {
		return remoteIP
	}
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		if ip := firstForwardedPublicIP(r.Header.Get(header)); ip != "" {
			return ip
		}
	}
	return remoteIP
}

func firstForwardedPublicIP(raw string) string {
	for _, part := range strings.Split(raw, ",") {
		ip := NormalizeIP(part)
		if ip != "" && !IsPrivateIP(ip) {
			return ip
		}
	}
	return ""
}

func IsAllLetters(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) || r > 'z' {
			return false
		}
	}
	return true
}

func IsValidHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

// Android ID 风格：16位小写 hex（常见的 android_id 是 64-bit 的 hex）
func AndroidID() string {
	b := make([]byte, 8) // 8 bytes = 16 hex chars
	if _, err := crand.Read(b); err != nil {
		log.Error("err:%v", err)
		return GenerateRandomString(16)
	}
	return hex.EncodeToString(b)
}

// AnyToString 将 any 类型安全转换为 string
func AnyToString(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case json.Number:
		return val.String()
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.Itoa(int(val))
	case int16:
		return strconv.Itoa(int(val))
	case int32:
		return strconv.Itoa(int(val))
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		// 对其他类型尝试用 fmt.Sprintf
		return fmt.Sprintf("%v", val)
	}
}

func NewRequestID() string {
	return fmt.Sprintf("%d%d",
		time.Now().UnixNano(),
		rand.Uint32(),
	)
}
func IsZeroIntOrString(v any) bool {
	switch n := v.(type) {
	case string:
		return n == "0"
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(n).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(n).Uint() == 0
	case float32, float64:
		return reflect.ValueOf(n).Float() == 0
	default:
		return false
	}
}
