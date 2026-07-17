package util

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func ToFloat32(v interface{}) float32 {
	switch v.(type) {
	case string:
		f, _ := strconv.ParseFloat(v.(string), 32)
		return float32(f)
	case []byte:
		var f float32
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &f)
		return f
	case int:
		return float32(v.(int))
	case int8:
		return float32(v.(int8))
	case int16:
		return float32(v.(int16))
	case int32:
		return float32(v.(int32))
	case int64:
		return float32(v.(int64))
	case uint8:
		return float32(v.(uint8))
	case uint16:
		return float32(v.(uint16))
	case uint32:
		return float32(v.(uint32))
	case uint64:
		return float32(v.(uint64))
	case float32:
		return v.(float32)
	case float64:
		return float32(v.(float64))
	default:
		return 0
	}
}

func ToFloat64(v interface{}) float64 {
	switch v.(type) {
	case string:
		f, _ := strconv.ParseFloat(v.(string), 64)
		return f
	case []byte:
		var f float64
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &f)
		return f
	case int:
		return float64(v.(int))
	case int8:
		return float64(v.(int8))
	case int16:
		return float64(v.(int16))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	case uint:
		return float64(v.(uint))
	case uint8:
		return float64(v.(uint8))
	case uint16:
		return float64(v.(uint16))
	case uint32:
		return float64(v.(uint32))
	case uint64:
		return float64(v.(uint64))
	case float32:
		return float64(v.(float32))
	case float64:
		return v.(float64)
	default:
		return 0
	}
}

func ToInt(v interface{}) int {
	switch v.(type) {
	case string:
		i, _ := strconv.Atoi(v.(string))
		return i
	case []byte:
		var i int
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return v.(int)
	case int8:
		return int(v.(int8))
	case int16:
		return int(v.(int16))
	case int32:
		return int(v.(int32))
	case int64:
		return int(v.(int64))
	case uint8:
		return int(v.(uint8))
	case uint16:
		return int(v.(uint16))
	case uint32:
		return int(v.(uint32))
	case uint64:
		return int(v.(uint64))
	case float32:
		return int(v.(float32))
	case float64:
		return int(v.(float64))
	default:
		return 0
	}
}

func ToUint(v interface{}) uint {
	switch v.(type) {
	case string:
		i, _ := strconv.Atoi(v.(string))
		return uint(i)
	case []byte:
		var i int
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint(i)
	case int:
		return uint(v.(int))
	case int8:
		return uint(v.(int8))
	case int16:
		return uint(v.(int16))
	case int32:
		return uint(v.(int32))
	case int64:
		return uint(v.(int64))
	case uint:
		return v.(uint)
	case uint8:
		return uint(v.(uint8))
	case uint16:
		return uint(v.(uint16))
	case uint32:
		return uint(v.(uint32))
	case uint64:
		return uint(v.(uint64))
	case float32:
		return uint(v.(float32))
	case float64:
		return uint(v.(float64))
	default:
		return 0
	}
}

func ToInt8(v interface{}) int8 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 8)
		return int8(i)
	case []byte:
		var i int8
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int8(v.(int))
	case int8:
		return v.(int8)
	case int16:
		return int8(v.(int16))
	case int32:
		return int8(v.(int32))
	case int64:
		return int8(v.(int64))
	case uint:
		return int8(v.(uint))
	case uint8:
		return int8(v.(uint8))
	case uint16:
		return int8(v.(uint16))
	case uint32:
		return int8(v.(uint32))
	case uint64:
		return int8(v.(uint32))
	case float32:
		return int8(v.(float32))
	case float64:
		return int8(v.(float64))
	default:
		return 0
	}
}

func ToUint8(v interface{}) uint8 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 8)
		return uint8(i)
	case []byte:
		var i int8
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint8(i)
	case int:
		return uint8(v.(int))
	case int8:
		return uint8(v.(int8))
	case int16:
		return uint8(v.(int16))
	case int32:
		return uint8(v.(int32))
	case int64:
		return uint8(v.(int64))
	case uint:
		return uint8(v.(uint))
	case uint8:
		return v.(uint8)
	case uint16:
		return uint8(v.(uint16))
	case uint32:
		return uint8(v.(uint32))
	case uint64:
		return uint8(v.(uint64))
	case float32:
		return uint8(v.(float32))
	case float64:
		return uint8(v.(float64))
	default:
		return 0
	}
}

func ToInt16(v interface{}) int16 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 16)
		return int16(i)
	case []byte:
		var i int16
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int16(v.(int))
	case int8:
		return int16(v.(int8))
	case int16:
		return v.(int16)
	case int32:
		return int16(v.(int32))
	case int64:
		return int16(v.(int64))
	case uint:
		return int16(v.(uint))
	case uint8:
		return int16(v.(uint8))
	case uint16:
		return int16(v.(uint16))
	case uint32:
		return int16(v.(uint32))
	case uint64:
		return int16(v.(uint64))
	case float32:
		return int16(v.(float32))
	case float64:
		return int16(v.(float64))
	default:
		return 0
	}
}

func ToUint16(v interface{}) uint16 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 16)
		return uint16(i)
	case []byte:
		var i int16
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint16(i)
	case int:
		return uint16(v.(int))
	case int8:
		return uint16(v.(int8))
	case int16:
		return uint16(v.(int16))
	case int32:
		return uint16(v.(int32))
	case int64:
		return uint16(v.(int64))
	case uint:
		return uint16(v.(uint))
	case uint8:
		return uint16(v.(uint8))
	case uint16:
		return v.(uint16)
	case uint32:
		return uint16(v.(uint32))
	case uint64:
		return uint16(v.(uint64))
	case float32:
		return uint16(v.(float32))
	case float64:
		return uint16(v.(float64))
	default:
		return 0
	}
}

func ToInt32(v interface{}) int32 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 32)
		return int32(i)
	case []byte:
		var i int32
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int32(v.(int))
	case int8:
		return int32(v.(int8))
	case int16:
		return int32(v.(int16))
	case int32:
		return v.(int32)
	case int64:
		return int32(v.(int64))
	case uint:
		return int32(v.(uint))
	case uint8:
		return int32(v.(uint8))
	case uint16:
		return int32(v.(uint16))
	case uint32:
		return int32(v.(uint32))
	case uint64:
		return int32(v.(uint64))
	case float32:
		return int32(v.(float32))
	case float64:
		return int32(v.(float64))
	default:
		return 0
	}
}

func ToUint32(v interface{}) uint32 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 32)
		return uint32(i)
	case []byte:
		var i int32
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint32(i)
	case int:
		return uint32(v.(int))
	case int8:
		return uint32(v.(int8))
	case int16:
		return uint32(v.(int16))
	case int32:
		return uint32(v.(int32))
	case int64:
		return uint32(v.(int64))
	case uint:
		return uint32(v.(uint))
	case uint8:
		return uint32(v.(uint8))
	case uint16:
		return uint32(v.(uint16))
	case uint32:
		return v.(uint32)
	case uint64:
		return uint32(v.(uint64))
	case float32:
		return uint32(v.(float32))
	case float64:
		return uint32(v.(float64))
	default:
		return 0
	}
}

func ToInt64(v interface{}) int64 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 64)
		return i
	case []byte:
		var i int64
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int64(v.(int))
	case int8:
		return int64(v.(int8))
	case int16:
		return int64(v.(int16))
	case int32:
		return int64(v.(int32))
	case int64:
		return v.(int64)
	case uint:
		return int64(v.(uint))
	case uint8:
		return int64(v.(uint8))
	case uint16:
		return int64(v.(uint16))
	case uint32:
		return int64(v.(uint32))
	case uint64:
		return int64(v.(uint64))
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	default:
		return 0
	}
}

func ToUint64(v interface{}) uint64 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 64)
		return uint64(i)
	case []byte:
		var i int64
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint64(i)
	case int:
		return uint64(v.(int))
	case int8:
		return uint64(v.(int8))
	case int16:
		return uint64(v.(int16))
	case int32:
		return uint64(v.(int32))
	case int64:
		return uint64(v.(int64))
	case uint:
		return uint64(v.(uint))
	case uint8:
		return uint64(v.(uint8))
	case uint16:
		return uint64(v.(uint16))
	case uint32:
		return uint64(v.(uint32))
	case uint64:
		return v.(uint64)
	case float32:
		return uint64(v.(float32))
	case float64:
		return uint64(v.(float64))
	default:
		return 0
	}
}
