package cfg

import (
	"strconv"
	"time"
)

// DecodeInt is a decoder for the int type.
func DecodeInt(b []byte) (int, error) {
	return strconv.Atoi(string(b))
}

// DecodeUint is a decoder for the uint type.
func DecodeUint(b []byte) (uint, error) {
	v, err := strconv.ParseUint(string(b), 0, 8)
	return uint(v), err
}

// DecodeInt64 is a decoder for the int64 type.
func DecodeInt64(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}

// DecodeUint64 is a decoder for the uint64 type.
func DecodeUint64(b []byte) (uint64, error) {
	return strconv.ParseUint(string(b), 10, 64)
}

// DecodeFloat64 is a decoder for the float64 type.
func DecodeFloat64(b []byte) (float64, error) {
	return strconv.ParseFloat(string(b), 64)
}

// DecodeBool is a decoder for the bool type.
func DecodeBool(b []byte) (bool, error) {
	return strconv.ParseBool(string(b))
}

// DecodeString is a decoder for the string type.
func DecodeString(b []byte) (string, error) {
	return string(b), nil
}

// DecodeDuration is a decoder for the time.Duration type.
func DecodeDuration(b []byte) (time.Duration, error) {
	return time.ParseDuration(string(b))
}
