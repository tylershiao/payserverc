package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/liangdas/mqant/log"
)

const (
	key = "HXSQnFRhs1v1cEty"
)

// UnPackMD5 统一的解析方法
func UnPackMD5(body []byte) string {
	pkg := rawPack{}
	if err := json.Unmarshal(body, &pkg); err != nil {
		log.Error("umarshal msg fail %v", err)
		return ""
	}
	sign := pkg.Sign
	data := pkg.Data
	if CalculateMD5(data) != sign {
		log.Error("check sign fail")
		return ""
	}
	return data
}

type rawPack struct {
	Sign string
	Data string
}

// CalculateMD5 calculate md5
func CalculateMD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func CalculateSHA256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)

	// Convert the hash bytes to a hexadecimal string
	// hashString := hex.EncodeToString(hashBytes)
	hashString := fmt.Sprintf("%x", hashBytes)

	return hashString
}

// HttpUseProxy 使用代理请求
func HttpUseProxy(proxyURL, target string, ret interface{}) error {
	uri, err := url.Parse(proxyURL)
	if err != nil {
		log.Error("parse url error: ", err)
		return err
	}
	client := http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(uri),
		},
		Timeout: 5 * time.Second,
	}
	log.Debug("reqURL:%v", target)
	resp, err := client.Get(target)
	if err != nil {
		log.Error("get err:%v", err)
		return err
	}
	// log.Debug("ret:%v", resp)
	if resp.StatusCode != http.StatusOK {
		log.Error("req fail err code:%v", resp.StatusCode)
		return errors.New("http req fail")
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read err %v", err)
		return err
	}
	if err := json.Unmarshal(data, ret); err != nil {
		return err
	}
	return nil
}

// HttpGet Get请求
func HttpGet(target string, ret interface{}) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(target)
	if err != nil {
		log.Error("get err:%v", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("req fail err code:%v", resp.StatusCode)
		return errors.New("http req fail")
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read err %v", err)
		return err
	}
	if err := json.Unmarshal(data, ret); err != nil {
		log.Error("Unmarshal err %v", err)
		return err
	}
	return nil
}

// HttpPost post请求
func HttpPost(target string, send, ret interface{}, header map[string]string) error {
	reqStr, _ := json.Marshal(send)
	requestID := NewRequestID()
	log.Debug("[%s] Post to:%v,req:%v", requestID, target, string(reqStr))
	req, err := http.NewRequest("POST", target, bytes.NewBuffer(reqStr))
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	// req.Header.Set("Accept", "application/problem+json")
	// req.Header.Set("Authorization", "api-key "+config.GetConfig().Web.OrangePay.APISecret)
	// req.Header.Set("Authorization", "api-key "+OrangeAPISecret)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	log.Debug("[%s] req:%+v", requestID, req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("[%s] http post call err:%v", requestID, err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("[%s] req fail err code:%v", requestID, resp.StatusCode)
		return errors.New("http req fail")
	}
	// if resp.StatusCode != http.StatusOK {
	// 	log.Error("req fail err code:%v", resp.StatusCode)
	// 	return errors.New("http req fail")
	// }
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Debug("[%s] response Body:%v", requestID, string(body))
	if err := json.Unmarshal(body, ret); err != nil {
		log.Error("unmarshal fail err:%v", err)
		return err
	}
	return nil
}

// HttpPostWithoutRet post请求
func HttpPostWithoutRet(target string, send interface{}, header map[string]string) error {
	reqStr, _ := json.Marshal(send)
	requestID := NewRequestID()
	log.Debug("[%s] Post to:%v,req:%v", requestID, target, string(reqStr))
	req, err := http.NewRequest("POST", target, bytes.NewBuffer(reqStr))
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	// req.Header.Set("Accept", "application/problem+json")
	// req.Header.Set("Authorization", "api-key "+config.GetConfig().Web.OrangePay.APISecret)
	// req.Header.Set("Authorization", "api-key "+OrangeAPISecret)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}
	log.Debug("[%s] req:%+v", requestID, req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("[%s] http post call err:%v", requestID, err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("[%s] req fail err code:%v", requestID, resp.StatusCode)
		return errors.New("http req fail")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Debug("[%s] response Body:%v", requestID, string(body))
	return nil
}

// HttpPostForm post请求
func HttpPostForm(target string, send, ret interface{}, header map[string]string) error {
	requestID := NewRequestID()
	log.Debug("[%s] Post to:%v,req:%v", requestID, target, send)
	uv := url.Values{}
	typ := reflect.TypeOf(send).Elem()
	val := reflect.ValueOf(send).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if val.Field(i).IsZero() {
			continue
		}
		uv.Add(typ.Field(i).Tag.Get("json"), fmt.Sprintf("%v", val.Field(i).Interface()))
		// ul += fmt.Sprintf("%v=%v&", typ.Field(i).Tag.Get("json"), val.Field(i).Interface())
	}
	log.Debug("[%s] uv:%v", requestID, uv)
	// ul += fmt.Sprintf("%v=%v", "sign", util.CalculateMD5(oid+strconv.Itoa(apiReq.AdvChannel)+apiReq.ClientDate+apiReq.Money+values.ZYAppKey))

	req, err := http.NewRequest("POST", target, strings.NewReader(uv.Encode()))
	if err != nil {
		log.Error("[%s] err:%v", requestID, err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	// req.Header.Set("Accept", "application/problem+json")
	// req.Header.Set("Authorization", "api-key "+config.GetConfig().Web.OrangePay.APISecret)
	// req.Header.Set("Authorization", "api-key "+OrangeAPISecret)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	log.Debug("[%s] req:%+v", requestID, req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("[%s] http post call err:%v", requestID, err)
		return err
	}
	// if resp.StatusCode != http.StatusOK {
	// 	log.Error("req fail err code:%v", resp.StatusCode)
	// 	return errors.New("http req fail")
	// }
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Debug("[%s] response Body:%v", requestID, string(body))
	if err := json.Unmarshal(body, ret); err != nil {
		log.Error("[%s]unmarshal fail err:%v", requestID, err)
		return err
	}
	return nil
}

// 将http的请求头替换为https
func ReplaceHTTPWithHTTPS(text string) string {
	// 创建正则表达式，用于匹配 HTTP 链接
	httpPattern := regexp.MustCompile(`http://[^\s]+`)

	// 将 HTTP 替换为 HTTPS
	httpsText := httpPattern.ReplaceAllStringFunc(text, func(match string) string {
		httpsURL := "https" + strings.TrimPrefix(match, "http")
		return httpsURL
	})

	return httpsText
}

func ParseFormReq(va url.Values, req interface{}) {
	typ := reflect.TypeOf(req).Elem()
	val := reflect.ValueOf(req).Elem()
	for i := 0; i < typ.NumField(); i++ {
		one := val.Field(i)
		value := va.Get(typ.Field(i).Tag.Get("json"))
		if value == "" {
			continue
		}
		if one.Kind() == reflect.Int || one.Kind() == reflect.Int64 {
			this, err := strconv.Atoi(value)
			if err != nil {
				log.Error("err:%v", err)
			}
			one.SetInt(int64(this))
		} else if one.Kind() == reflect.Float32 || one.Kind() == reflect.Float64 {
			this, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Error("err:%v", err)
			}
			one.SetFloat(this)
		} else if one.Kind() == reflect.Bool {
			this, err := strconv.ParseBool(value)
			if err != nil {
				log.Error("err:%v", err)
			}
			one.SetBool(this)
		} else {
			one.SetString(value)
		}
	}
}

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

// HttpPostReturnBody post请求返回body
func HttpPostReturnBody(target string, send interface{}, header map[string]string) (string, error) {
	reqStr, _ := json.Marshal(send)
	log.Debug("Post to:%v,req:%v", target, string(reqStr))
	req, err := http.NewRequest("POST", target, bytes.NewBuffer(reqStr))
	if err != nil {
		log.Error("err:%v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	// req.Header.Set("Accept", "application/problem+json")
	// req.Header.Set("Authorization", "api-key "+config.GetConfig().Web.OrangePay.APISecret)
	// req.Header.Set("Authorization", "api-key "+OrangeAPISecret)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	log.Debug("req:%+v", req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("http post call err:%v", err)
		return "", err
	}
	// if resp.StatusCode != http.StatusOK {
	// 	log.Error("req fail err code:%v", resp.StatusCode)
	// 	return errors.New("http req fail")
	// }
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	ret := string(body)
	log.Debug("response Body:%v", ret)
	return ret, nil
}
