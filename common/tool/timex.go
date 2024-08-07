package tool

import (
	"errors"
	"time"
)

var (
	CstZone = time.FixedZone("CST", int((8 * time.Hour).Seconds()))
)

// ToFormatTime 获取上海时区特定格式的数据
func ToFormatTime(tm, tmFormat string) (time.Time, error) {
	if tm == "" {
		return time.Time{}, errors.New("value is null")
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation(tmFormat, tm, loc)
}

// GetDateCst 获取 cst 时区今天时间为 0 的 Time
func GetDateCst() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, CstZone)
}

// GetMondayOfWeek 获取给定时间当周周一
func GetMondayOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return t.AddDate(0, 0, offset)
}

// GetMondayOfNextWeek 获取给定时间下周周一
func GetMondayOfNextWeek(t time.Time) time.Time {
	return GetMondayOfWeek(t).AddDate(0, 0, 7)
}
