package tool

import (
	"errors"
	"fmt"
	"math"

	"github.com/icholy/utm"
	"github.com/spf13/cast"
)

var (
	gridSize = 100000.0
	wgs84A   = 6378137.0
	wgs84B   = 6356752.31424518
	wgs84F   = 0.0033528107
	wgs84E   = 0.0818191908
	wgs84Ep  = 0.0820944379

	utmK0  = 0.9996
	utmFe  = 500000.0
	utmFnN = 0.0
	utmFnS = 10000000.0

	utmE2    = wgs84E * wgs84E
	utmE4    = utmE2 * utmE2
	utmE6    = utmE4 * utmE2
	utmEp2   = utmE2 / (1 - utmE2)
	raoTODeg = 180.0 / math.Pi
	degToRad = math.Pi / 180.0
)

type LonLat struct {
	Lon float64 // 经度
	Lat float64 // 纬度
}

type Zone struct {
	Number int  // Zone number 1 to 60
	Letter rune // Zone letter C to X (omitting O, I)
	North  bool // Zone hemisphere
}

type UTM struct {
	East  float64
	North float64
	Zone  Zone
}

type Mercator struct {
	Lon float64
	Lat float64
}

func LonLatConvertUTM(lonLat LonLat) UTM {
	east, north, zone := utm.ToUTM(lonLat.Lat, lonLat.Lon)
	return UTM{
		East:  east,
		North: north,
		Zone: Zone{
			Number: zone.Number,
			Letter: zone.Letter,
			North:  zone.North,
		},
	}
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// LonLatConvertUTM2 公司内部经纬度转utm 【算法直译】
func LonLatConvertUTM2(lonLat LonLat) UTM {
	a := wgs84A
	eccSquared := utmE2
	k0 := utmK0

	// Make sure the longitude is between -180.00 .. 179.9
	longTemp := (lonLat.Lon + 180) - cast.ToFloat64(cast.ToInt((lonLat.Lon+180)/360)*360) - 180

	latRad := deg2rad(lonLat.Lat)
	longRad := deg2rad(longTemp)

	zoneNumber := int((longTemp+180)/6) + 1

	if lonLat.Lat >= 56.0 && lonLat.Lat < 64.0 && longTemp >= 3.0 && longTemp < 12.0 {
		zoneNumber = 32
	}

	// Special zones for Svalbard
	if lonLat.Lat >= 72.0 && lonLat.Lat < 84.0 {
		if longTemp >= 0.0 && longTemp < 9.0 {
			zoneNumber = 31
		} else if longTemp >= 9.0 && longTemp < 21.0 {
			zoneNumber = 33
		} else if longTemp >= 21.0 && longTemp < 33.0 {
			zoneNumber = 35
		} else if longTemp >= 33.0 && longTemp < 42.0 {
			zoneNumber = 37
		}
	}

	// +3 puts origin in middle of zone
	longOrigin := (zoneNumber-1)*6 - 180 + 3
	longOriginRad := deg2rad(cast.ToFloat64(longOrigin))

	// compute the UTM Zone from the latitude and longitude
	// sprintf(UTMZone, "%d%c", zoneNumber, UTMLetterDesignator(Lat));

	eccPrimeSquared := (eccSquared) / (1 - eccSquared)

	n := a / math.Sqrt(1-eccSquared*math.Sin(latRad)*math.Sin(latRad))
	t := math.Tan(latRad) * math.Tan(latRad)
	c := eccPrimeSquared * math.Cos(latRad) * math.Cos(latRad)
	a1 := math.Cos(latRad) * (longRad - longOriginRad)

	M := a * ((1-eccSquared/4-3*eccSquared*eccSquared/64-5*eccSquared*eccSquared*eccSquared/256)*latRad -
		(3*eccSquared/8+3*eccSquared*eccSquared/32+45*eccSquared*eccSquared*eccSquared/1024)*math.Sin(2*latRad) +
		(15*eccSquared*eccSquared/256+45*eccSquared*eccSquared*eccSquared/1024)*math.Sin(4*latRad) -
		(35*eccSquared*eccSquared*eccSquared/3072)*math.Sin(6*latRad))

	utmEasting := cast.ToFloat64(k0*n*(a1+(1-t+c)*a1*a1*a1/6+(5-18*t+t*t+72*c-58*eccPrimeSquared)*a1*a1*a1*a1*a1/120) + 500000.0)

	utmNorth := cast.ToFloat64(k0 * (M + n*math.Tan(latRad)*(a1*a1/2+(5-t+9*c+4*c*c)*a1*a1*a1*a1/24+(61-58*t+t*t+600*c-330*eccPrimeSquared)*a1*a1*a1*a1*a1*a1/720)))

	if lonLat.Lat < 0 {
		//10000000 meter offset for southern hemisphere
		utmNorth += 10000000.0
	}

	// 获取zone
	_, _, zone := utm.ToUTM(lonLat.Lat, lonLat.Lon)

	return UTM{
		East:  utmEasting,
		North: utmNorth,
		Zone: Zone{
			Number: zoneNumber,
			Letter: zone.Letter,
			North:  zone.North,
		},
	}

}

func UTMConvertLonLatByZoneString(east, north float64, zoneStr string) (LonLat, error) {
	zone, b := utm.ParseZone(zoneStr)
	if !b {
		return LonLat{}, errors.New("invalid utm coord")
	}
	lat, lon := zone.ToLatLon(east, north)
	return LonLat{
		Lon: lon,
		Lat: lat,
	}, nil
}

func UTMConvertLonLat(u UTM) (LonLat, error) {
	zone, b := utm.ParseZone(fmt.Sprintf("%d%c", u.Zone.Number, u.Zone.Letter))
	if !b {
		return LonLat{}, errors.New("invalid utm coord")
	}
	lat, lon := zone.ToLatLon(u.East, u.North)
	return LonLat{
		Lon: lon,
		Lat: lat,
	}, nil
}

// UTMConvertLonLat2 公司内部utm转经纬度 【算法直译】
func UTMConvertLonLat2(u UTM) (LonLat, error) {
	k0 := utmK0
	a := wgs84A
	eccSquared := utmE2
	e1 := (1 - math.Sqrt(1-eccSquared)) / (1 + math.Sqrt(1-eccSquared))

	y := u.North
	x := u.East - 500000.0

	if u.Zone.North != true {
		y = y - 10000000.0
	}

	LongOrigin := (u.Zone.Number-1)*6 - 180 + 3
	eccPrimeSquared := (eccSquared) / (1 - eccSquared)

	M := y / k0
	mu := M / (a * (1 - eccSquared/4 - 3*eccSquared*eccSquared/64 - 5*eccSquared*eccSquared*eccSquared/256))

	phi1Rad := mu + ((3*e1/2-27*e1*e1*e1/32)*math.Sin(2*mu) + (21*e1*e1/16-55*e1*e1*e1*e1/32)*math.Sin(4*mu) + (151*e1*e1*e1/96)*math.Sin(6*mu))

	N1 := a / math.Sqrt(1-eccSquared*math.Sin(phi1Rad)*math.Sin(phi1Rad))
	T1 := math.Tan(phi1Rad) * math.Tan(phi1Rad)
	C1 := eccPrimeSquared * math.Cos(phi1Rad) * math.Cos(phi1Rad)
	R1 := a * (1 - eccSquared) / math.Pow(1-eccSquared*math.Sin(phi1Rad)*math.Sin(phi1Rad), 1.5)
	D := x / (N1 * k0)
	Lat := phi1Rad - ((N1 * math.Tan(phi1Rad) / R1) * (D*D/2 - (5+3*T1+10*C1-4*C1*C1-9*eccPrimeSquared)*D*D*D*D/24 + (61+90*T1+298*C1+45*T1*T1-252*eccPrimeSquared-3*C1*C1)*D*D*D*D*D*D/720))

	Lat = Lat * raoTODeg

	Long := (D - (1+2*T1+C1)*D*D*D/6 + (5-2*C1+28*T1-3*C1*C1+8*eccPrimeSquared+24*T1*T1)*D*D*D*D*D/120) / math.Cos(phi1Rad)
	Long = cast.ToFloat64(LongOrigin) + Long*raoTODeg

	return LonLat{
		Lon: Long,
		Lat: Lat,
	}, nil
}

// GetDistanceByLonLat 返回米
func GetDistanceByLonLat(lonLat1, lonLat2 LonLat) float64 {
	radius := 6371000 // 地球半径
	rad := math.Pi / 180.0

	lat1 := lonLat1.Lat * rad
	lng1 := lonLat1.Lon * rad
	lat2 := lonLat2.Lat * rad
	lng2 := lonLat2.Lon * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * float64(radius)
}

func GetDistanceByUtm(utm1, utm2 UTM) float64 {
	return math.Sqrt(math.Pow(utm1.East-utm2.East, 2) + math.Pow(utm1.North-utm2.North, 2))
}

func LongLatToMercator(lonLat LonLat) Mercator {
	mercator := Mercator{
		Lon: 0,
		Lat: 0,
	}
	earthRed := 6378137.0 //地球半径
	mercator.Lat = lonLat.Lon * math.Pi / 180 * earthRed
	param := lonLat.Lat * math.Pi / 180
	mercator.Lon = earthRed / 2 * math.Log((1.0+math.Sin(param))/(1.0-math.Sin(param)))
	return mercator
}
