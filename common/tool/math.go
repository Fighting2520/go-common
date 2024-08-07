package tool

import (
	"math"
	"strconv"
)

// Round 保留小数点后位数
// value float64 浮点数
// precise int 需保留小数点后的位数
func Round(value float64, precise int) float64 {
	value, _ = strconv.ParseFloat(strconv.FormatFloat(value, 'f', precise, 64), 64)
	return value
}

func Radian2Angle(r float64) float64 {
	return r * 180 / math.Pi
}

func Angle2Radian(a float64) float64 {
	return a * math.Pi / 180
}

// IntMax 获取最大值
func IntMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// IntMin 获取最小值
func IntMin(a, b int) int {
	if a > b {
		return b
	}
	return a
}
