package tool

import (
	"github.com/spf13/cast"
	"strings"
)

// Exclude 排除特定的字符串
func Exclude(list []string, excludeElems []string) []string {
	if len(excludeElems) == 0 {
		return list
	}
	var resp []string
	for _, s := range list {
		var excludeFlag bool
		for _, elem := range excludeElems {
			if s == elem {
				excludeFlag = true
				break
			}
		}
		if !excludeFlag {
			resp = append(resp, s)
		}
	}
	return resp
}

// ToUniqueIntSlice 整型数组去重
func ToUniqueIntSlice(i []int) []int {
	var data []int

	var aMap = make(map[int]struct{}, len(i))

	for _, v := range i {
		if _, ok := aMap[v]; ok {
			continue
		}
		aMap[v] = struct{}{}
		data = append(data, v)
	}

	return data
}

// ToUniqueStringSlice 字符串数组去重
func ToUniqueStringSlice(i []string) []string {
	var data []string

	var aMap = make(map[string]struct{}, len(i))

	for _, v := range i {
		if _, ok := aMap[v]; ok {
			continue
		}
		aMap[v] = struct{}{}
		data = append(data, v)
	}

	return data
}

// IsInIntSlice 检查元素是否在整型切片中
func IsInIntSlice(s1 []int, ele int) bool {
	for _, v := range s1 {
		if v == ele {
			return true
		}
	}
	return false
}

// IsInStringSlice 检查元素是否在字符串切片中
func IsInStringSlice(s1 []string, ele string) bool {
	for _, v := range s1 {
		if v == ele {
			return true
		}
	}
	return false
}

// IntersectIntSlice 整型切片求交集
func IntersectIntSlice(a []int, b []int) []int {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	var m = make(map[int]struct{}, len(a))
	for _, i := range a {
		m[i] = struct{}{}
	}

	var inter []int
	for _, i := range b {
		if _, ok := m[i]; ok {
			inter = append(inter, i)
		}
	}
	return inter
}

// IntersectStringSlice 字符串切片求交集
func IntersectStringSlice(a []string, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	var m = make(map[string]struct{}, len(a))
	for _, i := range a {
		m[i] = struct{}{}
	}

	var inter []string
	for _, i := range b {
		if _, ok := m[i]; ok {
			inter = append(inter, i)
		}
	}
	return inter
}

// UnionIntSlice 整型切片求并集
func UnionIntSlice(a []int, b []int) []int {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	var unionMap = make(map[int]struct{})
	var union []int

	for _, i := range a {
		if _, ok := unionMap[i]; ok {
			continue
		}
		unionMap[i] = struct{}{}
		union = append(union, i)
	}

	for _, i := range b {
		if _, ok := unionMap[i]; ok {
			continue
		}
		unionMap[i] = struct{}{}
		union = append(union, i)
	}

	return union
}

// UnionStringSlice 字符串切片求并集
func UnionStringSlice(a []string, b []string) []string {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	var union []string
	var unionMap = make(map[string]struct{})

	for _, i := range a {
		if _, ok := unionMap[i]; ok {
			continue
		}
		unionMap[i] = struct{}{}
		union = append(union, i)
	}

	for _, i := range b {
		if _, ok := unionMap[i]; ok {
			continue
		}
		unionMap[i] = struct{}{}
		union = append(union, i)
	}
	return union
}

// DifferenceIntSlice 整型切片求差集
func DifferenceIntSlice(a []int, b []int) (inA []int, inB []int) {
	if len(a) == 0 || len(b) == 0 {
		return a, b
	}

	inA = make([]int, 0)
	inB = make([]int, 0)

	var aMap = make(map[int]struct{}, len(a))
	for _, i := range a {
		aMap[i] = struct{}{}
	}

	for _, i := range b {
		if _, ok := aMap[i]; ok {
			delete(aMap, i)
			continue
		}
		inB = append(inB, i)
	}

	if len(aMap) > 0 {
		for _, i := range a {
			if _, ok := aMap[i]; ok {
				inA = append(inA, i)
			}
		}
	}

	return inA, inB
}

// DifferenceStringSlice 字符串切片求差集
func DifferenceStringSlice(a []string, b []string) (inA []string, inB []string) {
	if len(a) == 0 || len(b) == 0 {
		return a, b
	}

	inA = make([]string, 0)
	inB = make([]string, 0)

	var aMap = make(map[string]struct{}, len(a))
	for _, i := range a {
		aMap[i] = struct{}{}
	}

	for _, i := range b {
		if _, ok := aMap[i]; ok {
			delete(aMap, i)
			continue
		}
		inB = append(inB, i)
	}

	if len(aMap) > 0 {
		for _, i := range a {
			if _, ok := aMap[i]; ok {
				inA = append(inA, i)
			}
		}
	}

	return inA, inB
}

// JoinSliceToString 切片转字符串
func JoinSliceToString(elems interface{}, sep string) (string, error) {
	s, err := ToStringSliceE(elems)
	if err != nil {
		return "", err
	}

	return strings.Join(s, sep), nil
}

// GetChunkDataWithIntArr 获取整型切片分组数据
func GetChunkDataWithIntArr(targets []int, sep int) [][]int {
	num := len(targets)
	var quantity int
	if num%sep == 0 {
		quantity = cast.ToInt(num / sep)
	} else {
		quantity = (num / sep) + 1
	}

	var start, end, i int
	var segments = make([][]int, 0)
	for i = 1; i <= quantity; i++ {
		end = i * sep
		if i != quantity {
			segments = append(segments, targets[start:end])
		} else {
			segments = append(segments, targets[start:])
		}
		start = i * sep
	}

	return segments
}

// GetChunkDataWithStringArr 获取字符串切片分组数据
func GetChunkDataWithStringArr(targets []string, sep int) [][]string {
	num := len(targets)
	var quantity int
	if num%sep == 0 {
		quantity = cast.ToInt(num / sep)
	} else {
		quantity = (num / sep) + 1
	}

	var start, end, i int
	var segments = make([][]string, 0)
	for i = 1; i <= quantity; i++ {
		end = i * sep
		if i != quantity {
			segments = append(segments, targets[start:end])
		} else {
			segments = append(segments, targets[start:])
		}
		start = i * sep
	}

	return segments
}
