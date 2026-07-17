package routers

import (
	"fmt"
	"server/common"
	"server/db"

	"github.com/gin-gonic/gin"
)

func pageRouter(e *gin.RouterGroup) {
	e.GET("pay", Pay)
}

var (
	basePage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>payment</title>
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">

  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }

    html, body {
      width: 100%%;
      height: 100%%;
      background-color: #f9f9f9;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", sans-serif;
    }

    .iframe-wrapper {
      width: 100%%;
      height: 100vh;
      overflow: hidden;
      position: relative;
    }

    iframe {
      width: 100%%;
      height: 100%%;
      border: none;
    }

    /* 可选：加载动画 */
    .loading {
      position: absolute;
      top: 50%%;
      left: 50%%;
      transform: translate(-50%%, -50%%);
      font-size: 16px;
      color: #888;
    }
  </style>
</head>
<body>

<div class="iframe-wrapper">
  <div class="loading" id="loading">Loading...</div>
  <iframe src="%v"  id="myIframe" onload="document.getElementById('loading').style.display='none'"></iframe>
</div>

</body>
</html>`
)

func Pay(c *gin.Context) {
	id := c.Query("order")
	order := &common.OrderPay{OrderID: id}
	db.Mysql().Get(order)
	if order.PayUrl == "" {
		return
	}
	page := fmt.Sprintf(basePage, order.PayUrl)
	// c.String(200, page)
	c.Data(200, "text/html; charset=utf-8", []byte(page))
}
