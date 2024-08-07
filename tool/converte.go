package tool

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

// ToStringSliceE casts an interface to a []string type.
func ToStringSliceE(i interface{}) ([]string, error) {
	var a []string

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []string:
		return v, nil
	case []int:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int8:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int16:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint8:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint16:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case string:
		return strings.Fields(v), nil
	case []error:
		for _, err := range i.([]error) {
			a = append(a, err.Error())
		}
		return a, nil
	case interface{}:
		str, err := cast.ToStringE(v)
		if err != nil {
			return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		}
		return []string{str}, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	}
}
