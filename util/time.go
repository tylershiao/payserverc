package util

import "time"

func IsSameDay(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func IsSameDayTimeStamp(tm1 int64, tm2 int64) bool {
	t1 := time.Unix(tm1, 0)
	t2 := time.Unix(tm2, 0)

	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// GetFirstDateOfMonth 获取本月第一天零点
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// GetLastDateOfMonth 获取本月最后一天零点
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetWeekZeroTime 获取本周第一天的0点时间
func GetWeekZeroTime(d time.Time) time.Time {
	wd := int(d.Weekday())
	if wd == 0 {
		wd = 7
	}
	todayZero := GetZeroTime(d)
	return todayZero.AddDate(0, 0, 1-wd)
}

// GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func IsSameWeek(a, b int64) bool {
	// 首先判断是否相差7天内
	dif := a - b
	if dif < 0 {
		dif = -dif
	}
	if dif > 7*24*6060 {
		return false
	}
	ta := time.Unix(a, 0)
	tb := time.Unix(b, 0)
	taw := ta.Weekday()
	tbw := tb.Weekday()
	// if taw == time.Sunday {
	// 	taw = 7
	// }
	// if tbw == time.Sunday {
	// 	tbw = 7
	// }
	if ta.After(tb) {
		return taw >= tbw
	}
	return tbw >= taw
}

func IsSameMonth(a, b int64) bool {
	ta := time.Unix(a, 0)
	tb := time.Unix(b, 0)
	return GetFirstDateOfMonth(ta) == GetFirstDateOfMonth(tb)
}

func GetMonthDay(t time.Time) int {
	y, m, _ := t.Date()
	if m == 2 {
		if ((y%4) == 0 && (y%100) != 0) || (y%400) == 0 {
			return 29
		} else {
			return 28
		}
	}
	if m == 4 || m == 6 || m == 9 || m == 11 {
		return 30
	}
	return 31
}
