package cfg

import (
	"fmt"
	"strconv"
	"time"
)

func Encode(v interface{}) []byte {
	switch v := v.(type) {
	case int:
		return EncodeInt(v)
	case int64:
		return EncodeInt64(v)
	case uint:
		return EncodeUint(v)
	case uint64:
		return EncodeUint64(v)
	case float64:
		return EncodeFloat64(v)
	case bool:
		return EncodeBool(v)
	case string:
		return EncodeString(v)
	case time.Duration:
		return EncodeDuration(v)
	case []byte:
		return v
	default:
		return EncodeString(fmt.Sprintf("%v", v))
	}
}

// DecodeInt is a decoder for the int type.
func DecodeInt(b []byte) (int, error) {
	return strconv.Atoi(string(b))
}

// EncodeInt is an encoder for the int type.
func EncodeInt(v int) []byte {
	return []byte(strconv.Itoa(v))
}

// DecodeUint is a decoder for the uint type.
func DecodeUint(b []byte) (uint, error) {
	v, err := strconv.ParseUint(string(b), 0, 8)
	return uint(v), err
}

// EncodeUint is an encoder for the uint type.
func EncodeUint(v uint) []byte {
	return []byte(strconv.FormatUint(uint64(v), 10))
}

// DecodeInt64 is a decoder for the int64 type.
func DecodeInt64(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}

// EncodeInt64 is an encoder for the int64 type.
func EncodeInt64(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}

// DecodeUint64 is a decoder for the uint64 type.
func DecodeUint64(b []byte) (uint64, error) {
	return strconv.ParseUint(string(b), 10, 64)
}

// EncodeUint64 is an encoder for the uint64 type.
func EncodeUint64(v uint64) []byte {
	return []byte(strconv.FormatUint(v, 10))
}

// DecodeFloat64 is a decoder for the float64 type.
func DecodeFloat64(b []byte) (float64, error) {
	return strconv.ParseFloat(string(b), 64)
}

// EncodeFloat64 is an encoder for the float64 type.
func EncodeFloat64(v float64) []byte {
	return []byte(strconv.FormatFloat(v, 'f', -1, 64))
}

// DecodeBool is a decoder for the bool type.
func DecodeBool(b []byte) (bool, error) {
	return strconv.ParseBool(string(b))
}

// EncodeBool is an encoder for the bool type.
func EncodeBool(v bool) []byte {
	return []byte(strconv.FormatBool(v))
}

// DecodeString is a decoder for the string type.
func DecodeString(b []byte) (string, error) {
	return string(b), nil
}

// EncodeString is an encoder for the string type.
func EncodeString(v string) []byte {
	return []byte(v)
}

// DecodeDuration is a decoder for the time.Duration type.
func DecodeDuration(b []byte) (time.Duration, error) {
	return time.ParseDuration(string(b))
}

// EncodeDuration is an encoder for the time.Duration type.
func EncodeDuration(v time.Duration) []byte {
	return []byte(v.String())
}
