package util

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

func StructToMap(obj interface{}, tag string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	var result = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get(tag)
		if tagName != "" && tagName != "-" {
			result[tagName] = v.Field(i).Interface()
		}
	}

	return result
}

func GetStructFieldByJsonTag(obj interface{}, tag string) (field reflect.StructField, ok bool) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)
		tagName := field.Tag.Get("json")
		if tagName == tag {
			ok = true
			return
		}
	}
	return
}

func StructToMapJson(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	var result = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}
		tagName := t.Field(i).Tag.Get("json")
		if tagName != "" && tagName != "-" && !v.Field(i).IsZero() {
			result[tagName] = v.Field(i).Interface()
		}
	}

	return result
}

// StructToMapJsonAll 非空字段也包括
func StructToMapJsonAll(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	var result = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}
		tagName := t.Field(i).Tag.Get("json")
		if tagName != "" && tagName != "-" {
			result[tagName] = v.Field(i).Interface()
		}
	}

	return result
}

func Map2Struct(obj map[string]interface{}, result interface{}) error {
	return mapstructure.Decode(obj, result)
}

// DecodeRedisMap 解析map到目标t中,目前只实现了部分方法，如果需要解析更复杂结构需实现
func DecodeRedisMap(m map[string]string, t interface{}) (err error) {
	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)
		fi := val.Field(i)
		name := f.Name
		if s, ok := m[name]; ok {
			switch f.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ret, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}
				fi.SetInt(ret)
			case reflect.Float32, reflect.Float64:
				ret, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}
				fi.SetFloat(ret)
			case reflect.String:
				fi.SetString(s)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ret, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return err
				}
				fi.SetUint(ret)
			case reflect.Bool:
				ret, err := strconv.ParseBool(s)
				if err != nil {
					return err
				}
				fi.SetBool(ret)
			default:
				fmt.Println(t)
				continue
			}
		}
	}
	return nil
}

// CopySameStruct 将b中的非空字段复制到相同的结构体a中
func CopySameStruct(a, b interface{}) bool {
	typa := reflect.TypeOf(a)
	typb := reflect.TypeOf(b)
	if typa != typb {
		return false
	}
	if typa.Kind() != reflect.Ptr {
		return false
	}
	vala := reflect.ValueOf(a).Elem()
	valb := reflect.ValueOf(b).Elem()
	// typaE := typa.Elem()
	// typbE := typb.Elem()
	for i := 0; i < vala.NumField(); i++ {
		if valb.Field(i).IsZero() {
			continue
		}
		vala.Field(i).Set(valb.Field(i))
	}
	return true
}
