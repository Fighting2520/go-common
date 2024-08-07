package tool

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"reflect"
	"testing"
)

func TestExclude(t *testing.T) {
	type caseItem struct {
		in1       []string
		in2       []string
		expectOut []string
	}
	var cases = []caseItem{
		{
			in1:       []string{"a", "b", "c"},
			in2:       []string{"a", "d"},
			expectOut: []string{"b", "c"},
		},
		{
			in1:       []string{"a", "b", "c"},
			in2:       []string{},
			expectOut: []string{"a", "b", "c"},
		},
		{
			in1:       []string{""},
			in2:       []string{"a", "d"},
			expectOut: []string{""},
		},
	}
	for _, item := range cases {
		result := Exclude(item.in1, item.in2)
		assert.Equal(t, result, item.expectOut, fmt.Sprintf("返回值不匹配"))
	}
}

func TestToUniqueIntSlice(t *testing.T) {
	type args struct {
		i []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "int01", args: args{i: []int{1, 2, 3, 4, 5}}, want: []int{1, 2, 3, 4, 5}},
		{name: "int02", args: args{i: []int{1, 2, 3, 1, 2, 6}}, want: []int{1, 2, 3, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUniqueIntSlice(tt.args.i)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUniqueIntSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUniqueStringSlice(t *testing.T) {
	type args struct {
		i []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "int01", args: args{i: []string{"a", "c", "a", "f", "e", "g"}}, want: []string{"a", "c", "f", "e", "g"}},
		{name: "int02", args: args{i: []string{"a", "c", "d", "f", "f", "g"}}, want: []string{"a", "c", "d", "f", "g"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUniqueStringSlice(tt.args.i)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUniqueStringSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInIntSlice(t *testing.T) {
	type args struct {
		s1  []int
		ele int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "int1", args: args{s1: []int{1, 2, 3}, ele: int(1)}, want: true},
		{name: "int2", args: args{s1: []int{1, 2, 3}, ele: int(4)}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInIntSlice(tt.args.s1, tt.args.ele); got != tt.want {
				t.Errorf("IsInIntSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInStringSlice(t *testing.T) {
	type args struct {
		s1  []string
		ele string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "string1", args: args{s1: []string{"a", "b", "c"}, ele: "a"}, want: true},
		{name: "string2", args: args{s1: []string{"a", "b", "c"}, ele: "d"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInStringSlice(tt.args.s1, tt.args.ele); got != tt.want {
				t.Errorf("IsInStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersectIntSlice(t *testing.T) {
	type args struct {
		s1 []int
		s2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "int", args: args{s1: []int{1, 2, 3}, s2: []int{1, 2, 4}}, want: []int{1, 2}},
		{name: "int1", args: args{s1: []int{1, 2, 3}, s2: []int{1, 5, 4}}, want: []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntersectIntSlice(tt.args.s1, tt.args.s2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntersectIntSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersectStringSlice(t *testing.T) {
	type args struct {
		s1 []string
		s2 []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "int", args: args{s1: []string{"a", "b", "c"}, s2: []string{"a", "c", "d"}}, want: []string{"a", "c"}},
		{name: "int1", args: args{s1: []string{"a", "b", "c"}, s2: []string{"a", "c", "d"}}, want: []string{"a", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntersectStringSlice(tt.args.s1, tt.args.s2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntersectIntSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnionIntSlice(t *testing.T) {
	type args struct {
		a []int
		b []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "int1", args: args{a: []int{1, 2, 3}, b: []int{4, 3, 2}}, want: []int{1, 2, 3, 4}},
		{name: "int1", args: args{a: []int{1, 2, 3}, b: []int{4, 4, 5}}, want: []int{1, 2, 3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnionIntSlice(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnionIntSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnionStringSlice(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "string1", args: args{a: []string{"a", "b", "c"}, b: []string{"d", "e", "f"}}, want: []string{"a", "b", "c", "d", "e", "f"}},
		{name: "string1", args: args{a: []string{"a", "b", "c"}, b: []string{"c", "e", "e"}}, want: []string{"a", "b", "c", "e"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnionStringSlice(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnionStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifferenceIntSlice(t *testing.T) {
	type args struct {
		a []int
		b []int
	}
	tests := []struct {
		name    string
		args    args
		wantInA []int
		wantInB []int
	}{
		{name: "int1", args: args{a: []int{1, 2, 3}, b: []int{2, 3, 4}}, wantInA: []int{1}, wantInB: []int{4}},
		{name: "int1", args: args{a: []int{1, 3}, b: []int{2, 3, 5}}, wantInA: []int{1}, wantInB: []int{2, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInA, gotInB := DifferenceIntSlice(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(gotInA, tt.wantInA) {
				t.Errorf("DifferenceIntSlice() gotInA = %v, want %v", gotInA, tt.wantInA)
			}
			if !reflect.DeepEqual(gotInB, tt.wantInB) {
				t.Errorf("DifferenceIntSlice() gotInB = %v, want %v", gotInB, tt.wantInB)
			}
		})
	}
}

func TestDifferenceStringSlice(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name    string
		args    args
		wantInA []string
		wantInB []string
	}{
		{name: "string1", args: args{a: []string{"a", "b", "c"}, b: []string{"d", "e", "f"}}, wantInA: []string{"a", "b", "c"}, wantInB: []string{"d", "e", "f"}},
		{name: "string2", args: args{a: []string{"d", "e", "f"}, b: []string{"d", "e", "f"}}, wantInA: make([]string, 0), wantInB: make([]string, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInA, gotInB := DifferenceStringSlice(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(gotInA, tt.wantInA) {
				t.Errorf("DifferenceIntSlice() gotInA = %v, want %v", gotInA, tt.wantInA)
			}
			if !reflect.DeepEqual(gotInB, tt.wantInB) {
				t.Errorf("DifferenceIntSlice() gotInB = %v, want %v", gotInB, tt.wantInB)
			}
		})
	}
}

func TestJoinSliceToString(t *testing.T) {
	type args struct {
		elems interface{}
		sep   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "int", args: args{elems: []int{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "int8", args: args{elems: []int8{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "int16", args: args{elems: []int16{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "int32", args: args{elems: []int32{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "int64", args: args{elems: []int64{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "uint", args: args{elems: []uint{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "uint8", args: args{elems: []uint8{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "uint16", args: args{elems: []uint16{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "uint32", args: args{elems: []uint32{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "uint64", args: args{elems: []uint64{1, 2, 3, 4}, sep: ","}, want: "1,2,3,4", wantErr: false},
		{name: "float32", args: args{elems: []float32{1.1, 2.2, 3.3, 4.4}, sep: ","}, want: "1.1,2.2,3.3,4.4", wantErr: false},
		{name: "float32", args: args{elems: []float32{1.1, 2.2, 3.3, 4.4}, sep: ","}, want: "1.1,2.2,3.3,4.4", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JoinSliceToString(tt.args.elems, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinSliceToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JoinSliceToString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChunkDataWithInt(t *testing.T) {
	type args struct {
		targets []int
		sep     int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "test_001",
			args: args{
				targets: []int{1, 2, 3, 4},
				sep:     2,
			},
			want: [][]int{
				{1, 2},
				{3, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChunkDataWithIntArr(tt.args.targets, tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChunkDataWithIntArr() = %v, want %v", got, tt.want)
			}
		})
	}
}
