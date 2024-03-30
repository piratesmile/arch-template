package utils

import (
	"unsafe"
)

func BytesToString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return unsafe.String(&data[0], len(data))
}

func StringToBytes(v string) []byte {
	return unsafe.Slice(unsafe.StringData(v), len(v))
}
