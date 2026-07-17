package middleware

import (
	"fmt"
	"net/http"
	"server/modules/pay/app"

	"github.com/gin-gonic/gin"
)

// 跨域访问：cross  origin resource share
func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Printf("%+v\n", *context.Request)
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			context.Writer.Header().Set("Access-Control-Allow-Origin", "*")

			// 设置允许访问所有域
			context.Header("Access-Control-Allow-Origin", "*")

			//服务器支持的所有跨域请求的方法
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")

			//允许跨域设置可以返回其他子段，可以自定义字段
			context.Header("Access-Control-Allow-Headers", "Channel,uuid,version,Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")

			// 允许浏览器（客户端）可以解析的头部 （重要）
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")

			//设置缓存时间
			context.Header("Access-Control-Max-Age", "172800")

			//允许客户端传递校验信息比如 cookie (重要)
			context.Header("Access-Control-Allow-Credentials", "true")

			// 设置返回格式是json
			context.Set("content-type", "application/json")
		}

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, app.R{Code: http.StatusOK, Data: "Options Request!"})
		}

		//处理请求
		context.Next()
	}
}
