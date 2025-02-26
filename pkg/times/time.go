package times

import (
	"errors"
	"time"
)

const (
	LayoutDate     = "2006-01-02"
	LayoutTime     = "15:04:05"
	layoutDateTime = "2006-01-02 15:04:05"
)

var cstZone *time.Location

func init() {
	cstZone = time.FixedZone("UTC", 8*3600)
}

func Location() *time.Location {
	return cstZone
}

// ParseTimeToStr 格式化时间
func ParseTimeToStr(tm time.Time) string {
	return parseTimeToStr(tm.In(cstZone), LayoutTime) //.In(cstZone)表示设置成中国标准时间
}

// ParseDateToStr 格式化日期
func ParseDateToStr(tm time.Time) string {
	return parseTimeToStr(tm.In(cstZone), LayoutDate)
}

// ParseDateTimeToStr 格式化日期时间
func ParseDateTimeToStr(tm time.Time) string {
	return parseTimeToStr(tm.In(cstZone), layoutDateTime)
}

// GetNowDateTimeStr 获取当前日期时间的字符串
func GetNowDateTimeStr() string {
	return parseTimeToStr(time.Now().In(cstZone), layoutDateTime)
}

// GetNowDateStr 获取当前日期的字符串
func GetNowDateStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutDate)
}

// GetNowTimeStr 获取当前时间的字符串
func GetNowTimeStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutTime)
}

// ParseDateTime 解析日期时间
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return parseStrToTime(dateTimeStr, layoutDateTime)
}

// ParseDate 解析日期
func ParseDate(dateStr string) (time.Time, error) {
	return parseStrToTime(dateStr, LayoutDate)
}

// ParseTime 解析时间
func ParseTime(timeStr string) (time.Time, error) {
	return parseStrToTime(timeStr, LayoutTime)
}

func parseTimeToStr(t time.Time, layout string) string {
	return t.Format(layout)
}

func parseStrToTime(str, layout string) (time.Time, error) {
	if str == "" {
		return time.Now(), errors.New("args It can't be empty")
	}
	return time.ParseInLocation(layout, str, cstZone)
}

func GetNowTime() time.Time {
	return time.Now().In(cstZone)
}

// FuncTiming 程序运行时间
func FuncTiming(fn func()) time.Duration {
	startT := GetNowTime()
	fn()
	return GetNowTime().Sub(startT) //返回时间差
}

// 1970-01-01 08:00:00 +0800 CST
var zeroTime = time.Unix(0, 0)

// IsZero reports whether t represents（代表） the zero time instant（时刻）
// 报告 t 是否表示零时刻
func IsZero(t time.Time) bool {
	return t.IsZero() || zeroTime.Equal(t)
}
