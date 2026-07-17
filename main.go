package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"server/modules/pay"
	"time"

	"github.com/liangdas/mqant/selector"

	"github.com/liangdas/mqant"
	"github.com/nats-io/nats.go"

	"server/config"

	_ "net/http/pprof"

	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/registry"
	"github.com/liangdas/mqant/registry/consul"
)

func main() {
	log.Info("启动服务器 ...")
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	modulefile := flag.String("module", "module.json", "module configure file")
	logdir := flag.String("log", "logs", "log file directory")
	wddir := flag.String("wd", ".", "Server work directory")
	processID := flag.String("pid", "development", "Server ProcessID?")
	bidir := flag.String("bi", "bi", "bi file directory?")

	baseConf := flag.String("baseConf", "../baseConf.toml", "base configure file")

	conf := flag.String("conf", "conf.toml", "server configure file")

	debug := flag.Bool("debug", true, "debug mode, output to screen?")

	flag.Parse()
	wd1, _ := os.Getwd()
	log.Info(wd1)
	if err := config.LoadBaseConfig(*baseConf); err != nil {
		panic(fmt.Errorf("load base config file error %+v", err))
	}

	if err := config.LoadConfig(*conf); err != nil {
		panic(fmt.Errorf("load config file error %+v", err))
	}

	log.Info("-------------- configure -------------------")
	log.Info("%+v", config.GetBase())
	log.Info("%+v", config.GetConfig())
	log.Info("---------------------------------------------")

	// if config.GetConfig().WorkID > 0 {
	// 	go http.ListenAndServe("localhost:"+strconv.Itoa(config.GetConfig().WorkID), nil)
	// }

	rs := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{config.GetRegistry()}
	})

	nc, err := nats.Connect(config.GetNats(), nats.MaxReconnects(10000))
	if err != nil {
		panic(fmt.Errorf("nats error %v", err))
	}

	log.Info("================  flags ================")
	log.Info("module %s", *modulefile)
	log.Info("log dir %s", *logdir)
	log.Info("wd dir %s", *wddir)
	log.Info("pid %s", *processID)
	log.Info("bi dir %s", *bidir)
	log.Info("debug %v", *debug)
	wd, _ := os.Getwd()
	log.Info("work directory %s", wd)
	log.Info("========================================")

	app := mqant.CreateApp(
		module.Parse(false),
		module.Configure(*modulefile),
		module.LogDir(*logdir),
		module.WorkDir(*wddir),
		module.ProcessID(*processID),
		module.BILogDir(*bidir),
		module.KillWaitTTL(1*time.Minute),
		module.Debug(*debug), // 只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		module.Nats(nc),      // 指定nats rpc
		module.Registry(rs),  // 指定服务发现
		module.RegisterTTL(20*time.Second),
		module.RegisterInterval(10*time.Second),
		module.Selector(selector.NewSelector(selector.Registry(rs))),
	)

	err = app.Run(
		pay.Module(),
	)

	if err != nil {
		panic(fmt.Errorf("app run failed %v", err))
	}
}
