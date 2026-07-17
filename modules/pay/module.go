package pay

import (
	"context"
	"net/http"
	"server/call"
	"server/config"
	"server/db"
	mdb "server/db/mysql"
	rdb "server/db/redis"
	"server/modules/pay/handler"
	"server/modules/pay/middleware"
	"server/modules/pay/routers"
	"strconv"
	"sync"
	"time"

	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	basemodule "github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/server"
)

// Pay 支付模块结构体
type Pay struct {
	basemodule.BaseModule
	httpSvr *http.Server
	sysSvr  *http.Server
}

// Module 创建支付模块实例的工厂函数
var Module = func() module.Module {
	this := new(Pay)
	return this
}

// GetType 获取模块类型
func (p *Pay) GetType() string {
	return "pay"
}

// Version 获取模块版本号
func (p *Pay) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}

// OnInit 模块初始化函数
func (p *Pay) OnInit(app module.App, settings *conf.ModuleSettings) {
	metadata := make(map[string]string)
	metadata["workID"] = strconv.Itoa(config.GetConfig().WorkID)
	p.BaseModule.OnInit(p, app, settings,
		server.RegisterInterval(5*time.Second),
		server.RegisterTTL(10*time.Second),
		server.Metadata(metadata),
	)
	call.NewCaller(p)
	call.NewSnowflake(int64(config.GetConfig().WorkID))
	db.InitDB(&mdb.MysqlClient{}, &rdb.RedisClient{})
	log.Info("[%v]module init finish, config:%+v", p.GetType(), config.GetConfig())
	log.Info("[%v]module init finish, base:%+v", p.GetType(), config.GetBase())

	if err := loadConfig(); err != nil {
		panic(err)
	}
	// call.InitReload(p.App.Transport())
	handler.InitTimer()
}

// Run 模块运行函数
func (p *Pay) Run(closeSig chan bool) {
	log.Info("[%v]module running", p.GetType())
	// call.InitTimeWheel(closeSig)
	middleware.InitApiRateLimits()
	p.startHttpServer()
	<-closeSig
	log.Info("[%v]module stop", p.GetType())
}

// startHttpServer 启动HTTP服务器
func (p *Pay) startHttpServer() {
	cfg := config.GetConfig().Pay
	router := routers.SetUpRouter()
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	go func() {
		if config.GetConfig().Pay.TLS {
			if err := srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil {
				log.Error("ListenAndServeTLS fail error:%v", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil {
				log.Error("ListenAndServe fail error:%v", err)
			}
		}
	}()
	// returning reference so caller can call Shutdown()
	p.httpSvr = srv

	if cfg.SysAddr != "" {
		router := routers.SetUpSysRouter()
		sysSrv := &http.Server{
			Addr:    cfg.SysAddr,
			Handler: router,
		}
		go func() {
			if err := sysSrv.ListenAndServe(); err != nil {
				log.Error("ListenAndServe fail error:%v", err)
			}
		}()
		p.sysSvr = sysSrv
	}
}

// OnDestroy 模块销毁函数
func (p *Pay) OnDestroy() {
	log.Info("模块开始销毁...")

	// 并行关闭服务
	var wg sync.WaitGroup
	wg.Add(2)

	// HTTP 服务 shutdown
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := p.httpSvr.Shutdown(ctx); err != nil {
			log.Error("HTTP Shutdown error: %v", err)
		} else {
			log.Info("HTTP 服务已关闭")
		}
	}()

	// 系统服务 shutdown
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := p.sysSvr.Shutdown(ctx); err != nil {
			log.Error("SysSvr Shutdown error: %v", err)
		} else {
			log.Info("SysSvr 已关闭")
		}
	}()

	// 等待服务 shutdown 完成
	wg.Wait()
	handler.WithdrawGroup.Wait()

	p.BaseModule.OnDestroy()

	log.Info("模块已销毁")
}
