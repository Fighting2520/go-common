package tool

import (
	"reflect"
	"testing"
	"time"
)

func TestToFormatTime(t *testing.T) {
	type args struct {
		tm       string
		tmFormat string
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{name: "ceshi001", args: args{tm: "2022-08-18 17:05:21", tmFormat: "2006-01-02 15:04:05"}, want: time.Date(2022, 8, 18, 17, 5, 21, 0, loc)},
		{name: "ceshi002", args: args{tm: "2022-09-18 17:05:21", tmFormat: "2006-01-02 15:04:05"}, want: time.Date(2022, 9, 18, 17, 5, 21, 0, loc)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFormatTime(tt.args.tm, tt.args.tmFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFormatTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFormatTime() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
