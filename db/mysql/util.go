package mysql

import (
	"fmt"
	"server/util"
)

func PackChannels(channel ...*int) string {
	if len(channel) == 0 {
		return ""
	}
	sql := ""
	for i, v := range channel {
		if v == nil {
			continue
		}
		if i == 0 {
			sql += " and ("
		}
		sql += fmt.Sprintf(" channel_id = %v", *v)
		if i != len(channel)-1 {
			sql += " or"
		} else {
			sql += ")"
		}
	}
	return sql
}

// 获取百分比
func GetPer(a, b int64) string {
	if b == 0 {
		return fmt.Sprintf("%d", a) + "%"
	}
	return util.FormatFloat(float64(a*100)/float64(b), 2) + "%"
}

// 获取比例小数点
func GetPoint(a, b int64) string {
	if b == 0 {
		return "0"
	}
	return util.FormatFloat(float64(a)/float64(b), 2)
}
