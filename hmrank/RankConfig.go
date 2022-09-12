package hmrank

import (
	"fmt"
	"time"
)

type RkTypeEnum int

const (
	RkTypeEnum_Default        RkTypeEnum = iota //默认，未设置
	RkTypeEnum_MinMax                           //值 从小到大，时间 从小到大
	RkTypeEnum_MinMax_TimeRev                   //值 从小到大，时间 从大到小
	RkTypeEnum_MaxMin                           //值 从大到小，时间 从小到大
	RkTypeEnum_MaxMin_TimeRev                   //值 从大到小，时间 从大到小
)

type CycleEnum int

const (
	CycleEnum_Default CycleEnum = iota //默认，按天设置，RkCycle表示几天一个周期
	CycleEnum_Week                     //按周设置，RkCycle表示几周一个周期
)

//排行榜配置
type RankConfig struct {
	RkName      string        //排行榜名字
	RkRedisKey  string        //排行榜放Redis的Key
	RkType      RkTypeEnum    //排行榜类型
	CycleType   CycleEnum     //周期类型
	CycleOffset time.Duration //周期偏移量
	RkCycle     int           //排行榜清排周期（单位看CycleType）-1表示没有周期
	RkValidity  int           //排行榜有效时间（天）-1表示无限，如果周期!=-1，有效时间单位=周期
}

//拿rediskey
func (cf *RankConfig) GetRedisKey() string {
	return cf.RkRedisKey
}

func (cf *RankConfig) GetRedisKeyByDay(day int) string {
	return fmt.Sprint(cf.RkRedisKey, "_", day)
}

//北京 时区的无周期时间天数
func (cf *RankConfig) GetDayNum(t time.Time) int {
	return int((t.Unix() + 28800) / 86400)
}

//按时间算出周期
func (cf *RankConfig) GetDayCycleNum(t time.Time) int {
	switch cf.CycleType {
	case CycleEnum_Default:
		day := cf.GetDayNum(t) //int((t.Unix() + 28800) / 86400)
		day = day / cf.RkCycle
		return day
	case CycleEnum_Week:
		day := cf.GetDayNum(t.Add(-cf.CycleOffset)) //int((t.Add(-cf.CycleOffset).Unix() + 28800) / 86400)
		day = day / 7 / cf.RkCycle
		return day
	}
	return 0

}

//算周期时间长度
func (cf *RankConfig) GetCycleDuration() time.Duration {
	switch cf.CycleType {
	case CycleEnum_Default:
		return time.Duration(cf.RkCycle) * 24 * time.Hour
	case CycleEnum_Week:
		return time.Duration(cf.RkCycle) * 24 * time.Hour * 7
	}
	return 0
}

func (cf *RankConfig) GetCycleDateByCyid(cyid int) time.Time {
	switch cf.CycleType {
	case CycleEnum_Default:
		return time.Unix(int64((cyid*cf.RkCycle*86400 - 28800)), 0).Add(cf.CycleOffset)
	case CycleEnum_Week:
		return time.Unix(int64((cyid*7*cf.RkCycle*86400 - 28800)), 0).Add(cf.CycleOffset)
	}
	return time.Now()
}
