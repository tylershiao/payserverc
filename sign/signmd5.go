package sign

import (
	"fmt"
	"server/util"
	"sort"
	"strings"

	"github.com/liangdas/mqant/log"
)

func SignMD5(send interface{}, key string, passStr ...string) string {
	signStr := GetSignStr(send, passStr...)
	signStr += "&key=" + key
	log.Debug("final signStr:%s", signStr)
	return strings.ToUpper(util.CalculateMD5(signStr))
}

func SignMD5Empty(send interface{}, key string, passStr ...string) string {
	signStr := GetSignStrEmpty(send, passStr...)
	signStr += "&key=" + key
	log.Debug("final signStr:%s", signStr)
	return strings.ToUpper(util.CalculateMD5(signStr))
}

func APISignMD5(send interface{}, key string, passStr ...string) string {
	signStr := GetAPISignStr(send, passStr...)
	signStr += "&key=" + key
	log.Debug("final signStr:%s", signStr)
	return strings.ToUpper(util.CalculateMD5(signStr))
}

func GetAPISignStr(send interface{}, pass ...string) string {
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
		signStr += fmt.Sprintf("%s=%s", v, FormatSignValue(m[v]))
		if i != len(str)-1 {
			signStr += "&"
		}
	}
	log.Debug("signStr:%s", signStr)
	return signStr
}
