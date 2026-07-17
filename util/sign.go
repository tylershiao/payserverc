package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// sha1加密
func HMACSHA1(keyStr, value string) string {

	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	//进行base64编码
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return res
}

// HMACSHA1X 先转16进制
func HMACSHA1X(keyStr, value string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	//进行base64编码
	tmp := fmt.Sprintf("%x", mac.Sum(nil))
	res := base64.StdEncoding.EncodeToString([]byte(tmp))
	return res
}
