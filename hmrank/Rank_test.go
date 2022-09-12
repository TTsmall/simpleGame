package hmrank_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/buguang01/util"
	"wawaji_pub/hmrank"
)

func TestRank(t *testing.T) {
	f := 987452419845698745612846513298645186453.45645126451
	s := util.NewStringAny(f)
	fmt.Println(f)
	fmt.Println(s.ToString())
}

func TestRankTimer(t *testing.T) {
	//周期排行榜，删明天要用的列表
	next := time.Now().Add(time.Second * 10)
	fmt.Println("autoTask", time.Until(next).Seconds())
	tk := time.NewTicker(time.Until(next))
	fmt.Println(time.Now().String())
	for {
		select {
		case t := <-tk.C:
			fmt.Println(t.String())
			next = next.Add(time.Minute)
			// tk.Reset(time.Until(next))
			fmt.Println(time.Now().String())

		}
	}
}

func TestWeek(t *testing.T) {
	td := time.Now()
	td = td.Add(-2 * 24 * time.Hour)
	fmt.Println(td)
	fmt.Println(td.Weekday())
	day := int(td.Unix() / 86400)
	fmt.Println(day / 7)
	t1970 := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	fmt.Println(t1970.Weekday())
}

var (
	ActGuessRankCf *hmrank.RankConfig
)

func init() {
	ActGuessRankCf = new(hmrank.RankConfig)
	ActGuessRankCf.RkCycle = 1
	ActGuessRankCf.RkName = "全民猜成语"
	ActGuessRankCf.RkRedisKey = "GuessingIdiomsScore"
	ActGuessRankCf.RkType = hmrank.RkTypeEnum_MaxMin
	ActGuessRankCf.RkValidity = 2
	ActGuessRankCf.CycleType = hmrank.CycleEnum_Week
	ActGuessRankCf.CycleOffset = time.Hour * 24 * 4 //周一开始
}

func TestDate(t *testing.T) {
	td := time.Now()
	td2 := td.Add(-time.Hour * 4)
	cy := ActGuessRankCf.GetDayCycleNum(td)
	cy2 := ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 5)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 8)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 8)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 8)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 8)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 4)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 4)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 4)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 4)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 4)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 4)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)
	td = td.AddDate(0, 0, 1)
	td2 = td.Add(-time.Hour * 4)
	cy = ActGuessRankCf.GetDayCycleNum(td)
	cy2 = ActGuessRankCf.GetDayCycleNum(td2)
	fmt.Println(cy, td, cy2, td2)

}

func TestDateDay(t *testing.T) {
	cy := ActGuessRankCf.GetDayCycleNum(time.Now())
	day := ActGuessRankCf.GetCycleDateByCyid(cy)
	fmt.Println(day)
	// td := ActGuessRankCf.GetCycleDuration() * time.Duration(cy+1) //- ActGuessRankCf.CycleOffset

	// day := util.GetMinDateTime()
	// fmt.Println(day, day.Unix(), day.UTC())
	// day = day.Add(td)
	// fmt.Println(day)
	// day = day.Add(ActGuessRankCf.GetCycleDuration())
	// fmt.Println(day)
}
