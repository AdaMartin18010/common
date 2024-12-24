package common

import (
	"bytes"
	"math"
)

// 比较精度 CompareMin
func FloatIsEqual(f1, f2, CompareMin float64) bool {
	if f1 > f2 {
		return math.Dim(f1, f2) < CompareMin
	} else {
		return math.Dim(f2, f1) < CompareMin
	}
}

// 比较精度
const intCompareMin = 1

func IntIsEqual(i1, i2 int) bool {
	if i1 > i2 {
		return i1-i2 < intCompareMin
	} else {
		return i2-i1 < intCompareMin
	}
}

func EraseControlChar(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("\b"), []byte(""))
	data = bytes.ReplaceAll(data, []byte("\f"), []byte(""))
	data = bytes.ReplaceAll(data, []byte("\t"), []byte(""))
	data = bytes.ReplaceAll(data, []byte("\n"), []byte(""))
	data = bytes.ReplaceAll(data, []byte("\r"), []byte(""))
	return data
}
