package config

type FacebookParams struct {
	Prefix    string
	AppKey    string
	AppScrect string
	AuthURL   string
}

type GooglePlayParams struct {
	Prefix  string
	AuthURL string
}

type Configure struct {
	Pay struct {
		Addr                string
		SysAddr             string // 内部系统调用接口
		BackSysAddr         string // 后台内部系统调用接口
		TLS                 bool
		CertFile            string
		KeyFile             string
		CallbackURL         string
		CheckLimit          int // 触发充值权重检测的阀值
		SelectPayWay        bool
		PayCheckTime        int
		PayFailWeight       int // 失败改变的权重值
		PaySuccessWeight    int // 成功改变的权重值
		BaseSuccess         int // 基准成功率
		WithdrawCheckTime   int
		WithdrawScanTime    int
		QueryPayTime        int
		QueryWithdrawTime   int
		WithdrawRiskAudit   bool
		TestWithdraw        int
		BaseWithdrawSuccess int64 // 退出基准成功率
		PayCheckCount       int   // 订单数低于该值不进行重置
		OrderRate           int   // 下单请求限制
		QueryRate           int   // 查询请求限制
		VirgoPay            struct {
			MID int
			URL string
			Key string
		}
		VirgoPay2 struct {
			MID int
			URL string
			Key string
		}
		EpayGlobal struct {
			PayCode         int
			WithdrawPayCode int
		}
		PaymentURL string
	}
	WorkID int
}

func GetConfig() *Configure {
	return &cfg
}

var cfg Configure

func LoadConfig(cfgFile string) error {
	return LoadToml(cfgFile, &cfg)
}
