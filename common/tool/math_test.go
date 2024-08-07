package tool

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestDecimal(t *testing.T) {
	var data = []struct {
		name  string
		value float64
		pecr  int
		want  float64
	}{
		{name: "ceshi01", value: 121.7892615040, pecr: 9, want: 121.789261504},
		{name: "ceshi02", value: 121.7892615049, pecr: 9, want: 121.789261505},
		{name: "ceshi03", value: 121.7892615042, pecr: 9, want: 121.789261504},
		{name: "ceshi04", value: 121.7892615047, pecr: 9, want: 121.789261505},
	}

	// 单元测试
	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			val := Round(test.value, test.pecr)
			assert.Equal(t, val == test.want, true)
		})
	}
}
