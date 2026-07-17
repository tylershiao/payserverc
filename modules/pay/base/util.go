package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
)

// GetSignStringFromJSON 保持原始 JSON 文本表示，生成签名字符串
func GetSignStringFromJSON(jsonStr string, passStr ...string) (string, error) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
	decoder.UseNumber() // 保留数字原始文本

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}

	// 删除 sign 字段
	for _, v := range passStr {
		delete(data, v)
	}

	// 按 key 升序排序
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接 k=value
	var parts []string
	for _, k := range keys {
		v := data[k]
		switch val := v.(type) {
		case string:
			parts = append(parts, fmt.Sprintf("%s=%s", k, val))
		case json.Number:
			parts = append(parts, fmt.Sprintf("%s=%s", k, val.String())) // 保持原始数字文本
		default:
			// 对于布尔或 null 等情况也保持原始格式
			raw, err := json.Marshal(val)
			if err != nil {
				return "", err
			}
			parts = append(parts, fmt.Sprintf("%s=%s", k, string(raw)))
		}
	}

	signString := ""
	for i, p := range parts {
		if i == 0 {
			signString = p
		} else {
			signString += "&" + p
		}
	}

	return signString, nil
}

func ReplaceDomain(rawURL, newDomain string) string {
	// 1. 解析旧 URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// 2. 替换域名（Host）
	u.Host = newDomain

	// 强制使用 https
	u.Scheme = "https"

	// 3. 返回完整 URL
	return u.String()
}
