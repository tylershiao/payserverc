package call

import (
	"time"

	timewheel "github.com/liangdas/mqant/module/modules/timer"
)

func InitTimeWheel(closeSig chan bool) {
	timewheel.GetTimeWheel().Stop()
	timewheel.SetTimeWheel(timewheel.New(100*time.Millisecond, 36))
	Go(func() {
		timewheel.GetTimeWheel().Start(closeSig)
	})
}
