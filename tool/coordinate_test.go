package tool

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/icholy/utm"
)

func TestConvert2UTM(t *testing.T) {
	data := []LonLat{
		{120.47662402058401, 31.61325787596599},
	}

	for _, val := range data {
		east, northing, zone := utm.ToUTM(val.Lat, val.Lon)
		fmt.Println(east, northing, zone)
	}

	log.Println("**************************************")
}

func TestUTMConvertLonLat(t *testing.T) {
	lonlat, err := UTMConvertLonLat(UTM{East: 260626.41996078, North: 3500333.8899216, Zone: Zone{Number: 51, Letter: 82, North: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lonlat)
}

func TestUTMConvertLonLat2(t *testing.T) {
	lonlat, err := UTMConvertLonLat2(UTM{East: 334325.7125, North: 3434097.8, Zone: Zone{Number: 51, Letter: 82, North: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lonlat)
}

func TestGetDistanceByLonLat(t *testing.T) {
	lat := GetDistanceByLonLat(LonLat{
		Lon: 116.644969,
		Lat: 39.881979,
	}, LonLat{
		116.644969, 39.881979,
	})
	t.Log(lat)
}

func TestGetDistanceByUtm(t *testing.T) {
	utm1 := LonLatConvertUTM(LonLat{Lon: 116.644969, Lat: 39.881979})
	utm2 := LonLatConvertUTM(LonLat{
		116.651063, 39.751321,
	})
	t.Log(utm1)
	t.Log(utm2)
	d := GetDistanceByUtm(utm1, utm2)
	t.Log(d)
}

func TestUTMConvertLonLatByZoneString(t *testing.T) {
	type args struct {
		east    float64
		north   float64
		zoneStr string
	}
	tests := []struct {
		name    string
		args    args
		want    LonLat
		wantErr bool
	}{
		{
			name: "test-001",
			args: args{
				east:    116.644969,
				north:   39.881979,
				zoneStr: "51R",
			},
			want: LonLat{
				Lon: 118.51230112540965,
				Lat: 0.0003608445577133357,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UTMConvertLonLatByZoneString(tt.args.east, tt.args.north, tt.args.zoneStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("UTMConvertLonLatByZoneString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UTMConvertLonLatByZoneString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLongLatToMercator(t *testing.T) {
	data := []LonLat{
		{121.8304121594476, 30.87016029869196},
		{121.81840555985546, 30.872337107243748},
	}

	for _, val := range data {
		t.Log(LongLatToMercator(val))
	}
}

func TestLonLatConvertUTM2(t *testing.T) {
	type args struct {
		lonLat LonLat
	}
	tests := []struct {
		name string
		args args
		want UTM
	}{
		{name: "test01", args: args{lonLat: LonLat{Lon: 116.644969, Lat: 39.881979}}, want: UTM{East: 469642.5160692344, North: 4414718.48497586, Zone: Zone{Number: 50, Letter: 83, North: true}}},
		{name: "test02", args: args{lonLat: LonLat{Lon: 121.26417573319446, Lat: 31.02889215973292}}, want: UTM{East: 334325.7124938039, North: 3434097.800059111, Zone: Zone{Number: 51, Letter: 82, North: true}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LonLatConvertUTM2(tt.args.lonLat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LonLatConvertUTM2() = %v, want %v", got, tt.want)
			}
		})
	}
}
