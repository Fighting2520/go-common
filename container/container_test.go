package container

import (
	"reflect"
	"testing"
)

func TestCreateContainerFactory(t *testing.T) {
	tests := []struct {
		name string
		want *containers
	}{
		{
			name: "test001",
			want: &containers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateContainerFactory()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateContainerFactory() = %v, want %v", got, tt.want)
			} else {
				_ = got.Set("key1", "value1")
				_ = got.Set("key2", "value1")
				_ = got.Set("test", "value1")
				if reflect.DeepEqual(got.Get("key1"), "value1") == false {
					t.Errorf("Set() = %v, want %v", got.Get("key1"), "value1")
				}
				_, check := got.KeyIsExists("key1")
				if check != true {
					t.Errorf("KeyIsExists() = %v, want %v", check, true)
				}
				got.FuzzyDelete("test")
				if reflect.DeepEqual(got.Get("test"), nil) == false {
					t.Errorf("FuzzyDelete() = %v, want %v", got.Get("test"), nil)
				}
				got.Delete("key1")
				if reflect.DeepEqual(got.Get("key1"), nil) == false {
					t.Errorf("Delete() = %v, want %v", got.Get("key1"), nil)
				}
			}
		})
	}
}
