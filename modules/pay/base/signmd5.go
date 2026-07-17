package base

import (
	"server/sign"
	"server/util"
	"strings"

	"github.com/liangdas/mqant/log"
)

func (b *Base) SignMD5(send interface{}) string {
	signStr := sign.GetSignStr(send, b.SignPassStr...)
	if b.PassKeyName {
		signStr += b.SignKey
	} else if b.KeyName == "" {
		signStr += "&key=" + b.SignKey
	} else {
		signStr += "&" + b.KeyName + "=" + b.SignKey
	}
	log.Debug("final signStr:%s", signStr)
	ret := util.CalculateMD5(signStr)
	if b.ShouldSignUpper {
		ret = strings.ToUpper(ret)
	}
	return ret
}

func (b *Base) SignMD5Empty(send interface{}) string {
	signStr := sign.GetSignStrEmpty(send, b.SignPassStr...)
	if b.PassKeyName {
		signStr += b.SignKey
	} else if b.KeyName == "" {
		signStr += "&key=" + b.SignKey
	} else {
		signStr += "&" + b.KeyName + "=" + b.SignKey
	}
	log.Debug("final signStr:%s", signStr)
	ret := util.CalculateMD5(signStr)
	if b.ShouldSignUpper {
		ret = strings.ToUpper(ret)
	}
	return ret
}

func (b *Base) SignMD5WithStr(str string) string {
	signStr := sign.GetSignStrNull(str, b.SignPassStr...)
	if b.PassKeyName {
		signStr += b.SignKey
	} else if b.KeyName == "" {
		signStr += "&key=" + b.SignKey
	} else {
		signStr += "&" + b.KeyName + "=" + b.SignKey
	}
	log.Debug("final signStr:%s", signStr)
	ret := util.CalculateMD5(signStr)
	if b.ShouldSignUpper {
		ret = strings.ToUpper(ret)
	}
	return ret
}

func (b *Base) SignMD5WithStrEmpty(str string) string {
	signStr := sign.GetSignWithStrEmpty(str, b.SignPassStr...)
	if b.PassKeyName {
		signStr += b.SignKey
	} else if b.KeyName == "" {
		signStr += "&key=" + b.SignKey
	} else {
		signStr += "&" + b.KeyName + "=" + b.SignKey
	}
	log.Debug("final signStr:%s", signStr)
	ret := util.CalculateMD5(signStr)
	if b.ShouldSignUpper {
		ret = strings.ToUpper(ret)
	}
	return ret
}

func (b *Base) SignSHA256(send interface{}) string {
	signStr := sign.GetSignStr(send, b.SignPassStr...)
	if b.PassKeyName {
		signStr += b.SignKey
	} else if b.KeyName == "" {
		signStr += "&key=" + b.SignKey
	} else {
		signStr += "&" + b.KeyName + "=" + b.SignKey
	}
	log.Debug("final signStr:%s", signStr)
	ret := util.CalculateSHA256(signStr)
	if b.ShouldSignUpper {
		ret = strings.ToUpper(ret)
	}
	return ret
}

func (b *Base) SignHmacSHA256Base64(send interface{}) string {
	signStr := sign.GetSignStr(send, b.SignPassStr...)
	ret := SignWithHmacSHA256Base64(signStr, b.SignKey)
	return ret
}

func (b *Base) SignHmacSHA1Base64(send interface{}) string {
	signStr := sign.GetSignStr(send, b.SignPassStr...)
	ret := SignWithHmacSHA1Base64(signStr, b.SignKey)
	return ret
}
