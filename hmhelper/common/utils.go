package common

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//TenThousand 1万
const TenThousand = 10000

const MaxInt = int(^uint(0) >> 1)

func RandNum(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func Sample(arr []int) (error, int) {
	l := len(arr)
	if l == 0 {
		return errors.New("common/util.go Sample can't take empty array"), 0
	}
	return nil, arr[rand.Intn(l)]
}
func SampleMustInt(arr []int) int {
	l := len(arr)
	if l == 0 {
		return 0
	}
	return arr[rand.Intn(l)]
}

//拿到时间0点
func ZeroTimeOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

//返回当前时间是t的第几天，t自己为第一天
func TimeSubOfDay(t time.Time) int {
	dt := ZeroTimeOfDay(t)
	nt := ZeroTimeOfDay(time.Now())
	return int(nt.Sub(dt).Hours())/24 + 1
}

func NormalizeTimeOfDay(t time.Time, startHour int) time.Time {
	if t.Hour() < startHour {
		year, month, day := t.AddDate(0, 0, -1).Date()
		return time.Date(year, month, day, startHour, 0, 0, 0, time.Local)
	} else {
		year, month, day := t.Date()
		return time.Date(year, month, day, startHour, 0, 0, 0, time.Local)
	}
}

func DiffDays(endTime time.Time, startTime time.Time) int {
	year2, month2, day2 := endTime.Date()
	year1, month1, day1 := startTime.Date()
	d2 := time.Date(
		year2, month2, day2,
		0, 0, 0, 0, time.Local,
	)
	d1 := time.Date(
		year1, month1, day1,
		0, 0, 0, 0, time.Local,
	)
	return int(d2.Sub(d1) / (24 * time.Hour))
}

func GetIpAddress(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(forwardedFor, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if ip != "" {
				return ip
			}
		}
	}
	ip := r.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip
	}
	index := strings.LastIndex(r.RemoteAddr, ":")
	if index < 0 {
		return r.RemoteAddr
	}
	return r.RemoteAddr[:index]
}

//HitRateTenThousand
//是否命中万分比
func HitRateTenThousand(rate int) bool {
	return rand.Intn(TenThousand) < rate
}

func GetOutboundIp() string {
	netAddr, err := net.ResolveTCPAddr("tcp", "www.baidu.com:80")
	if err != nil {
		panic("can not get outbound ip " + err.Error())
	}
	conn, err := net.Dial("udp", netAddr.IP.String()+":80")
	if err != nil {
		panic("can not get outbound ip " + err.Error())
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]

}

func GetTomorrowStamp() time.Time {
	tomorrow := time.Now().Add(24 * time.Hour)
	year, month, day := tomorrow.Date()
	//tomorrow_str := fmt.Sprintf("%d-%d-%d 00:00:00", year, month, day)
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func ConvertMap(m map[int]int) map[int32]int32 {
	newMap := make(map[int32]int32, len(m))
	for k, v := range m {
		newMap[int32(k)] = int32(v)
	}
	return newMap
}

func ConvertMapTo64(m map[int]int) map[int32]int64 {
	newMap := make(map[int32]int64, len(m))
	for k, v := range m {
		newMap[int32(k)] = int64(v)
	}
	return newMap
}

func SubString(source string, start, end int) string {
	var r = []rune(source)
	length := len(r)
	if start == 0 && end >= length {
		return source
	}

	if start < 0 || start > end {
		return ""
	}
	return string(r[start:end])
}

func WaitWriteChanTimeout(ch chan interface{}, data interface{}, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		ch <- data
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

type PosInfo struct {
	PosX   int
	PoxY   int
	Direct int
}

func (m *PosInfo) Scan(value interface{}) error {
	fmt.Println("=============PosInfo Scan data=", value)
	err := json.Unmarshal(value.([]byte), m)
	if err != nil {
		fmt.Println("PosInfo Scan err=", err)
		return err
	}
	return nil
	//return this.UnmarshalJSON(value.([]byte))
}

func (m PosInfo) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// func (m *PosInfo) Scan(value interface{}) error {
// 	if len(value.([]byte)) == 0 {
// 		return nil
// 	}
// 	var posInfo PosInfo
// 	err := json.Unmarshal(value.([]byte), &posInfo)
// 	if err != nil {
// 		return err
// 	}
// 	*m = posInfo
// 	return nil
// }

// func (m PosInfo) Value() (driver.Value, error) {
// 	return json.Marshal(m)
// }

// func (m PosInfo) MarshalJSON() ([]byte, error) {
// 	fmt.Println("MarshalJSON PosInfo")
// 	var buf bytes.Buffer
// 	buf.WriteByte('{')
// 	buf.WriteString(fmt.Sprintf(`%d,%d,%d`, m.PosX, m.PoxY, m.Direct))
// 	buf.WriteByte('}')
// 	return buf.Bytes(), nil
// }

// func (m *PosInfo) UnmarshalJSON(data []byte) error {
// 	*m = PosInfo{}
// 	if len(data) == 0 {
// 		return nil
// 	}
// 	mpHs := make([]int, 0)
// 	err := json.Unmarshal(data, &mpHs)
// 	if err != nil {
// 		return err
// 	}

// 	for _, v := range mpHs {
// 		mpHs = append(mpHs, v)
// 	}
// 	return nil
// }

type IntKv map[int]int

func (m IntKv) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for k, v := range m {
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf(`"%d":%d`, k, v))
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (m *IntKv) UnmarshalJSON(data []byte) error {
	*m = make(map[int]int)
	if len(data) <= 2 {
		return nil
	}
	var mp map[string]int
	err := json.Unmarshal(data, &mp)
	if err != nil {
		return err
	}
	for k, v := range mp {
		key, _ := strconv.Atoi(k)
		(*m)[key] = v
	}
	return nil
}

type RuinsItemsMap map[int][]RuinsItem
type ExploreItemsMap map[int][]ExploreItem
type RuinsItems []RuinsItem
type ExploreItems []ExploreItem

type ExploreItem struct {
	Id        int
	EndTime   int
	Hero      int
	WorkerIds IntSlice
	TileX     int
	TileY     int
	BuildId   int
}

type RuinsItem struct {
	Id    int
	Areas IntKv
	Uid   int
}

func (m *ExploreItemsMap) Scan(value interface{}) error {
	//fmt.Println("=============ExploreItems Scan data=", value)
	err := json.Unmarshal(value.([]byte), m)
	if err != nil {
		fmt.Println("UnmarshalJSON ExploreItems nil")
		return err
	}
	return nil
}

func (m ExploreItemsMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *RuinsItemsMap) Scan(value interface{}) error {
	//fmt.Println("=============RuinsItemsMap Scan data=", value)
	err := json.Unmarshal(value.([]byte), m)
	if err != nil {
		fmt.Println("UnmarshalJSON RuinsItemsMap nil")
		return err
	}
	return nil
}

func (m RuinsItemsMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m IntSlice) MarshalJSON() ([]byte, error) {
	var result = make([]string, len(m))
	for i, v := range m {
		result[i] = strconv.Itoa(v)
	}
	return []byte("[" + strings.Join(result, ",") + "]"), nil
}

func (m *IntSlice) UnmarshalJSON(data []byte) error {
	if len(data) <= 2 {
		*m = make([]int, 0)
		return nil
	}
	strs := strings.Split(string(data[1:len(data)-1]), ",")
	*m = make([]int, len(strs))
	for i, str := range strs {
		(*m)[i], _ = strconv.Atoi(strings.TrimSpace(str))
	}
	return nil
}

func (m *IntSlice) Scan(value interface{}) error {
	return m.UnmarshalJSON(value.([]byte))
}

func (m IntSlice) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MapIntSlice) Scan(value interface{}) error {
	//fmt.Println("=============ExploreItems Scan data=", value)
	err := json.Unmarshal(value.([]byte), m)
	if err != nil {
		fmt.Println("UnmarshalJSON ExploreItems nil")
		return err
	}
	return nil
}

func (m MapIntSlice) Value() (driver.Value, error) {
	return json.Marshal(m)
}
