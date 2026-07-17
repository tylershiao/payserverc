package sign

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/liangdas/mqant/log"
)

func GetSignStr(send interface{}, pass ...string) string {
	m := StructToMapJson(send)
	str := []string{}
	for i := range m {
		if i == "sign" {
			continue
		}
		shouldPass := false
		for _, v := range pass {
			if v == i {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		str = append(str, i)
	}
	sort.Strings(str)
	signStr := ""
	for i, v := range str {
		signStr += fmt.Sprintf("%v=%v", v, m[v])
		if i != len(str)-1 {
			signStr += "&"
		}
	}
	log.Debug("signStr:%s", signStr)
	return signStr
}

func GetSignStrl(send interface{}, pass ...string) string {
	m := StructToMapJson(send)
	str := []string{}
	for i := range m {
		if i == "sign" {
			continue
		}
		shouldPass := false
		for _, v := range pass {
			if v == i {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		str = append(str, i)
	}
	sort.Strings(str)
	signStr := ""
	for i, v := range str {
		signStr += fmt.Sprintf("%v", m[v])
		if i != len(str)-1 {
			signStr += "|"
		}
	}
	log.Debug("signStr:%s", signStr)
	return signStr
}

// 包含空字段
func GetSignStrEmpty(send interface{}, pass ...string) string {
	m := StructToMapJsonEmpty(send)
	str := []string{}
	for i := range m {
		if i == "sign" {
			continue
		}
		shouldPass := false
		for _, v := range pass {
			if v == i {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		str = append(str, i)
	}
	sort.Strings(str)
	signStr := ""
	for i, v := range str {
		signStr += fmt.Sprintf("%v=%v", v, m[v])
		if i != len(str)-1 {
			signStr += "&"
		}
	}
	log.Debug("signStr:%s", signStr)
	return signStr
}

// 连接参数时不加&符号，且不加key值，直接value拼接
func GetSignStr2(send interface{}, pass ...string) string {
	m := StructToMapJson(send)
	str := []string{}
	for i := range m {
		if i == "sign" {
			continue
		}
		shouldPass := false
		for _, v := range pass {
			if v == i {
				shouldPass = true
				break
			}
		}
		if shouldPass {
			continue
		}
		str = append(str, i)
	}
	sort.Strings(str)
	signStr := ""
	for _, v := range str {
		signStr += fmt.Sprintf("%v", m[v])
	}
	log.Debug("signStr:%s", signStr)
	return signStr
}

func StructToMapJson(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	var result = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}
		tagName := t.Field(i).Tag.Get("json")
		if tagName != "" && tagName != "-" && (v.Field(i).Kind() == reflect.Bool || !v.Field(i).IsZero()) {
			data := v.Field(i).Interface()
			if v.Field(i).Kind() == reflect.Struct {
				data = GetSignStr(data)
			} else if v.Field(i).Kind() == reflect.Slice {
				tmp, err := json.Marshal(data)
				if err != nil {
					log.Error("err:%v", err)
				}
				data = string(tmp)
			}
			name := strings.Split(tagName, ",")
			result[name[0]] = data
		}
	}

	return result
}

func StructToMapJsonEmpty(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	var result = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}
		tagName := t.Field(i).Tag.Get("json")
		if tagName != "" && tagName != "-" {
			data := v.Field(i).Interface()
			if v.Field(i).Kind() == reflect.Struct {
				data = GetSignStr(data)
			}
			name := strings.Split(tagName, ",")
			result[name[0]] = data
		}
	}

	return result
}

// 根据body里的字段直接拼接出签名字符串(包含null)
func GetSignStrWithNull(str string, pass ...string) string {
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

// 根据body里的字段直接拼接出签名字符串(去除null包含空值)
func GetSignWithStrEmpty(str string, pass ...string) string {
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
		// if string(m[v]) == "null" {
		// 	continue
		// }
		signStr += fmt.Sprintf("%v=%v", v, string(m[v]))
		signStr += "&"
	}
	signStr = signStr[:len(signStr)-1]
	log.Debug("signStr:%v", signStr)
	return signStr
}

func SHA512(data string) string {
	hash := sha512.Sum512([]byte(data))
	return fmt.Sprintf("%x", hash)
}

var gmailRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// SafeRandomGmail 返回 8~20 位、只含小写+数字、不含连续点的 Gmail
func SafeRandomGmail() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	const maxLocal = 20
	n := 8 + gmailRand.Intn(13) // 8~20
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = alphabet[gmailRand.Intn(len(alphabet))]
	}
	// 随机插一个点，更像真人
	if n > 9 && gmailRand.Intn(2) == 0 {
		pos := 3 + gmailRand.Intn(n-6)
		b = append(b[:pos], append([]byte{'.'}, b[pos:]...)...)
	}
	return string(b) + "@gmail.com"
}

func DecodeUnicode(str string) string {
	var result strings.Builder
	for i := 0; i < len(str); i++ {
		if i+5 < len(str) && str[i] == '\\' && str[i+1] == 'u' {
			hexVal := str[i+2 : i+6]
			var r rune
			fmt.Sscanf(hexVal, "%04x", &r)
			result.WriteRune(r)
			i += 5
		} else {
			result.WriteByte(str[i])
		}
	}
	return result.String()
}

func FormatSignValue(v interface{}) string {
	switch val := v.(type) {
	case json.Number:
		return val.String()
	case float64:
		// 判断是不是整数
		if val == float64(int64(val)) {
			return strconv.FormatInt(int64(val), 10)
		}
		// 小数，去掉科学计数法
		return strconv.FormatFloat(val, 'f', -1, 64)

	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)

	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)

	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)

	case string:
		return val

	default:
		return fmt.Sprintf("%v", val)
	}
}
