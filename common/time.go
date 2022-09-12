package common

import "time"

const (
	DAILY_RESET_HOUR = 5
)

var Time1970 = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
var Time2028 = time.Date(2028, 1, 1, 0, 0, 0, 0, time.Local)

//是否是同一天
func SameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func IsToday(ts int64) bool {
	return SameDay(time.Unix(ts, 0), time.Now())
}

func CheckTimeFormat(src, layout string) bool {
	_, err := time.Parse(layout, src)
	return err == nil
}

func GetDateWithOffset(t time.Time) string {
	return GetTimeWithOffset(t).Format("060102")
}

func GetTimeWithOffset(t time.Time) time.Time {
	return t.Add(-time.Hour * DAILY_RESET_HOUR)
}

func GetDateNoOffset(t time.Time) string {
	return t.Format("060102")
}

//GetNowDateWithOffset
//得到带重置时间的当天的YYYYMMHH
func GetNowDateWithOffset() string {
	return GetDateWithOffset(time.Now())
}

//得到当前时间到2028年的时间
//一般用在排行榜上,比如同等级,判断时间先后
func Time2TenYears() time.Duration {
	return Time2028.Sub(time.Now())
}

//等到tm1-tm2的天数，
func GetPassedDays(tm1, tm2 time.Time) int64 {
	return tm1.Unix()/86400 - tm2.Unix()/86400
}

func FormatSqlDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func IsSameDay(time1 time.Time, time2 time.Time) bool {
	return IsTmSameDay(time1.Unix(), time2.Unix(), DAILY_RESET_HOUR)
}

func IsTmSameDay(tm1 int64, tm2 int64, startHour int) bool {
	return GetPassedDaysVisTs(tm1, tm2, startHour) == 0
}

func GetPassedDaysVisTs(tm1 int64, tm2 int64, startHour int) int64 {
	return (tm1-Time1970.Unix()-int64(startHour)*3600)/86400 - (tm2-Time1970.Unix()-int64(startHour)*3600)/86400
}

// 获取今天的零点时间
func GetTodayBegin() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	return t.Unix()
}

func GetNowTimeUnixNano() int64 {
	return time.Now().UnixNano() / 1e6
}
