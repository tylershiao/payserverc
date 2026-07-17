package values

import (
	"fmt"
	"math/rand"
	"server/common"
	"server/util"
)

func PackPay(r *common.PayReq) {
	if len(r.Name) == 0 {
		r.Name = randName()
		// r.Name = util.GenerateRandomString(5)
	}
	r.Mobile = util.CheckPhone(r.Mobile)
	if !util.ValidateEmail(r.Email) {
		r.Email = util.GetSimpleRandomString(5) + "@" + "gmail.com"
	}
	if r.IP == "" {
		r.IP = fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}
	if r.UID == "" {
		r.UID = util.GenerateRandomString(6)
	}
	if r.DeviceID == "" {
		r.DeviceID = util.AndroidID()
	}
}

func PackWithdraw(r *common.WithdrawReq) {
	if len(r.Name) == 0 {
		r.Name = randName()
	}
	// r.Phone = util.CheckPhone(r.Phone)
	if !util.ValidateEmail(r.Email) {
		r.Email = util.GetSimpleRandomString(5) + "@" + "gmail.com"
	}
	if r.UID == "" {
		r.UID = util.GenerateRandomString(6)
	}
	if len(r.Address) == 0 {
		r.Address = "Mumbai"
	}
	if r.DeviceID == "" {
		r.DeviceID = util.AndroidID()
	}
}

func randName() string {
	return util.GenerateRandomString(5)
}
