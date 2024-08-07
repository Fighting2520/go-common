package tool

import (
	"fmt"
	"reflect"
)

// IsInSliceE 检测特定内容
func IsInSliceE(s1 interface{}, ele interface{}) (bool, error) {
	v1 := reflect.ValueOf(s1)
	if v1.Kind() != reflect.Slice {
		return false, fmt.Errorf("s1 type %T is not slice. type is file", reflect.ValueOf(s1))
	}

	for i := 0; i < v1.Len(); i++ {
		if v1.Index(i).Interface() == ele {
			return true, nil
		}
	}

	return false, nil
}
