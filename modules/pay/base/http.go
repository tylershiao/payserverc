package base

import (
	"encoding/json"
	"io"
	"net/http"
	"server/util"
	"time"

	"github.com/liangdas/mqant/log"
)

// status=1时，属于不知道第三方是什么情况的报错，不能当成业务中断
// status=0时，可以根据返回判断是否是业务错误
func (b *Base) Request(req *http.Request, ret interface{}) int {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req.Close = true
	requestID := util.NewRequestID()
	log.Debug("[%v]req:%+v", requestID, req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("[%v]http post call err:%v", requestID, err)
		return 1
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	b.Body = string(body)
	log.Debug("[%v]response Body:%v", requestID, b.Body)
	if err := json.Unmarshal(body, ret); err != nil {
		log.Error("[%v]unmarshal fail err:%v", requestID, err)
		return 1
	}
	return 0
}
