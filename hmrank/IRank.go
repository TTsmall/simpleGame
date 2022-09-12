package hmrank

import (
	"fmt"
	"math"
	"time"

	"github.com/buguang01/util"
)

//要写入排行榜的接口
type IRankWriteModel interface {
	GetRedisKey() []string //要写入的rediskey会有多个，因为有效
	GetUserKey() string    //用户的key
	GetRedisScore() string //拿写入的数据，包时间戳
	GetScore() string      //拿真数据
}

type RankWriteModel struct {
	UserKey    string       //对应用户的KEY
	Score      *util.String //写入的数据（值部分应该是整数部分，小数部分为时间数据）
	UpTime     time.Time    //写入时间
	RankConfig RankConfig   //需要的配置
}

func (md *RankWriteModel) GetRedisKey() []string {
	//因为我们是8时区的
	day := md.RankConfig.GetDayNum(md.UpTime) //int((md.UpTime.Unix() + 28800) / 86400)
	cf := md.RankConfig
	if cf.RkCycle == -1 && cf.RkValidity == -1 {
		//永久
		return []string{md.RankConfig.GetRedisKey()}
	}
	if cf.RkCycle != -1 {
		//写入的数据是有周期的
		day = cf.GetDayCycleNum(md.UpTime)
		// day = day / cf.RkCycle
		day = day % cf.RkValidity
		return []string{md.RankConfig.GetRedisKeyByDay(day)}
	}
	result := make([]string, cf.RkValidity)
	for i := 0; i < cf.RkValidity; i++ {
		tmp := (day + i) % cf.RkValidity
		result[i] = md.RankConfig.GetRedisKeyByDay(tmp)
	}
	return result
}

func (md *RankWriteModel) GetUserKey() string {
	return md.UserKey
}

func (md *RankWriteModel) GetRedisScore() string {
	score := util.NewString(md.Score.ToString())

	if index := score.Index("."); index != -1 {
		score = score.SubstringEnd(index)

	}
	switch md.RankConfig.RkType {
	case RkTypeEnum_MaxMin:
		{
			t := math.MaxInt32 - md.UpTime.Unix()
			return fmt.Sprintf("%s.%010d", score.String(), t)
		}
	case RkTypeEnum_MaxMin_TimeRev:
		{
			t := md.UpTime.Unix()
			return fmt.Sprintf("%s.%010d", score.String(), t)
		}
	case RkTypeEnum_MinMax:
		{
			t := md.UpTime.Unix()
			return fmt.Sprintf("%s.%010d", score.String(), t)
		}
	case RkTypeEnum_MinMax_TimeRev:
		{
			t := math.MaxInt32 - md.UpTime.Unix()
			return fmt.Sprintf("%s.%010d", score.String(), t)
		}
	}
	return "0"
}
func (md *RankWriteModel) GetScore() string {
	return md.Score.String()
}

//score

type RankModel struct {
	UserKey string //用户主键
	RankNo  int    //排名
	Score   string //用户数据

}
