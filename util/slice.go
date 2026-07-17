package util

import (
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/liangdas/mqant/log"
)

func ArrToString(res []int) string {
	var str string
	for i := 0; i < len(res); i++ {
		if i == len(res)-1 {
			str += strconv.Itoa(res[i])
		} else {
			str += strconv.Itoa(res[i]) + ","
		}
	}
	return str
}

func StringToArr(str string) []int {
	var res []int

	if len(str) == 0 {
		return res
	}

	arr := strings.Split(str, ",")

	for i := 0; i < len(arr); i++ {
		number, err := strconv.Atoi(arr[i])
		if err != nil {
			log.Debug(err.Error())
		}
		res = append(res, number)
	}
	return res
}

// SliceInt2Int32 数组转换
func SliceInt2Int32(c []int) []uint32 {
	ret := []uint32{}
	for _, v := range c {
		ret = append(ret, uint32(v))
	}
	return ret
}

// SliceInt322Int 数组转换
func SliceInt322Int(c []uint32) []int {
	ret := []int{}
	for _, v := range c {
		ret = append(ret, int(v))
	}
	return ret
}

// SliceInt2Int64 数组转换
func SliceInt2Int64(c []int) []int64 {
	ret := []int64{}
	for _, v := range c {
		ret = append(ret, int64(v))
	}
	return ret
}

// RemoveSlice 根据b中的元素，从a中剔除，相当于两个slice做差
func RemoveSlice(a, b []int) []int {
	ret := []int{}
	dm := map[int]int{}
	for _, v := range b {
		dm[v]++
	}
	for _, v := range a {
		if num, ok := dm[v]; ok {
			num--
			if num == 0 {
				delete(dm, v)
			} else {
				dm[v] = num
			}
			continue
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

// RemoveOne 从a中剔除b
func RemoveOne(a []int, b int) []int {
	ret := []int{}
	count := 0
	for _, v := range a {
		if v == b && count == 0 {
			count++
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

// SliceContain a中是否包含b
func SliceContainString(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// SliceContainIP a中是否包含b对应的IP
func SliceContainIP(a []string, b string) bool {
	normalized := NormalizeIP(b)
	if normalized == "" {
		return false
	}
	for _, v := range a {
		if NormalizeIP(v) == normalized {
			return true
		}
	}
	return false
}

// SliceContain a中是否包含b
func SliceContain(a []int, b int) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// SliceContainInt64 a中是否包含b
func SliceContainInt64(a []int64, b int64) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

func IsEqualSlice(a, b []int) bool {
	dm := map[int]int{}
	for _, v := range a {
		dm[v]++
	}
	for _, v := range b {
		dm[v]--
	}
	for _, v := range dm {
		if v != 0 {
			return false
		}
	}
	return true
}

// 根据权重,从数组中随机一个元素
// name 权重字段名(该字段值必须为整型)
func RandomOneEleFromSlice(slice interface{}, name string) int {
	ref := reflect.ValueOf(slice)
	var total int64
	for i := 0; i < ref.Len(); i++ {
		tmp := ref.Index(i)
		if tmp.Kind() == reflect.Pointer {
			tmp = tmp.Elem()
		}
		total += tmp.FieldByName(name).Int()
	}
	if total <= 0 {
		return -1
	}
	ran := rand.Int63n(total)
	var rans int64
	for i := 0; i < ref.Len(); i++ {
		tmp := ref.Index(i)
		if tmp.Kind() == reflect.Pointer {
			tmp = tmp.Elem()
		}
		rans += tmp.FieldByName(name).Int()
		if ran < rans {
			return i
		}
	}
	return 0
}

// 根据权重,从整数数组中随机一个元素
func RandomSliceWeight(slice []int) int {
	// ref := reflect.ValueOf(slice)
	var total int
	for _, v := range slice {
		total += v
	}
	ran := rand.Intn(total)
	rans := 0
	for i := 0; i < len(slice); i++ {
		rans += slice[i]
		if ran < rans {
			return i
		}
	}
	// for i := 0; i < ref.Len(); i++ {
	// 	tmp := ref.Index(i)
	// 	total += tmp.Int()
	// }
	// ran := rand.Int63n(total)
	// var rans int64
	// for i := 0; i < ref.Len(); i++ {
	// 	tmp := ref.Index(i)
	// 	rans += tmp.Int()
	// 	if ran < rans {
	// 		return i
	// 	}
	// }
	return 0
}

// 数组中每个元素扩大100倍
func Multi100(tmp []int) (ret []int) {
	for _, v := range tmp {
		ret = append(ret, v*100)
	}
	return
}
