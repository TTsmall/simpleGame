package hmrank

import (
	"time"

	"github.com/buguang01/util"
	"wawaji_pub/hmrank/errors"
)

//排行榜因为是其他系统的基础，所以他们单例是本地放在本地，方便引用
var (
	RankEx *RankManage
)

//初始化
func Init(rankmg *RankManage) {
	RankEx = rankmg
}

//写入排行榜
func SetRankScore(rkname, userkey string, score int, upt time.Time) (err error) {
	if rkcfmd, ok := RankEx.rankconflist[rkname]; ok {
		//找到配置
		rkmd := new(RankWriteModel)
		rkmd.Score = util.NewStringAny(score)
		rkmd.UserKey = userkey
		rkmd.UpTime = upt
		rkmd.RankConfig = *rkcfmd
		RankEx.SendRankModel(rkmd)
		return
	} else {
		return errors.Err_Rank_Config
	}
}

//拿排行榜,day 拿指定日期的排行榜
func GetRankScoreMyToX(rkname, userkey string, day time.Time, x int) (result []*RankModel, my *RankModel, err error) {
	start, end := 0, x
	if end == -1 {
		end = 99
	} else if end > 99 {
		start = x - 99
	}
	rdmd := RankEx.redismd.GetConn()
	util.Using(rdmd, func() {
		if rkcfmd, ok := RankEx.rankconflist[rkname]; ok {
			//找到配置
			key := ""
			if rkcfmd.RkCycle == -1 && rkcfmd.RkValidity == -1 {
				//永久的与时间无关
				key = rkcfmd.GetRedisKey()

			} else if rkcfmd.RkCycle == -1 {
				//没有周期
				tmday := rkcfmd.GetDayNum(day) //int((day.Unix() + 28800) / 86400)
				tmday = tmday % rkcfmd.RkValidity
				key = rkcfmd.GetRedisKeyByDay(tmday)
			} else {
				//有周期
				tmday := rkcfmd.GetDayCycleNum(day)
				tmday = tmday % rkcfmd.RkValidity
				key = rkcfmd.GetRedisKeyByDay(tmday)
			}
			switch rkcfmd.RkType {
			case RkTypeEnum_MinMax, RkTypeEnum_MinMax_TimeRev:
				result, my, err = getRankScoreToXMyByRange(rdmd, rkcfmd, key, userkey, start, end)
			case RkTypeEnum_MaxMin, RkTypeEnum_MaxMin_TimeRev:
				result, my, err = getRankScoreToXMyByRevRange(rdmd, rkcfmd, key, userkey, start, end)
			}
			return
		} else {
			err = errors.Err_Rank_Config
			return
		}
	})
	return
}

//获取排行榜数据，没有自己
func GetRankScoreToX(rkname string, day time.Time, x int) (result []*RankModel, err error) {
	start, end := 0, x
	if end == -1 {
		end = 99
	} else if end > 99 {
		start = x - 99
	}
	rdmd := RankEx.redismd.GetConn()
	util.Using(rdmd, func() {
		if rkcfmd, ok := RankEx.rankconflist[rkname]; ok {
			//找到配置
			key := ""
			if rkcfmd.RkCycle == -1 && rkcfmd.RkValidity == -1 {
				//永久的与时间无关
				key = rkcfmd.GetRedisKey()

			} else if rkcfmd.RkCycle == -1 {
				//没有周期
				tmday := rkcfmd.GetDayNum(day)
				tmday = tmday % rkcfmd.RkValidity
				key = rkcfmd.GetRedisKeyByDay(tmday)
			} else {
				//有周期
				tmday := rkcfmd.GetDayCycleNum(day)
				// tmday := int(day.Unix() / 86400)
				// tmday = tmday / rkcfmd.RkCycle
				tmday = tmday % rkcfmd.RkValidity
				key = rkcfmd.GetRedisKeyByDay(tmday)
			}
			switch rkcfmd.RkType {
			case RkTypeEnum_MinMax, RkTypeEnum_MinMax_TimeRev:
				result, err = getRankScoreToXByRange(rdmd, rkcfmd, key, start, end)
			case RkTypeEnum_MaxMin, RkTypeEnum_MaxMin_TimeRev:
				result, err = getRankScoreToXByRevRange(rdmd, rkcfmd, key, start, end)
			}
			return
		} else {
			err = errors.Err_Rank_Config
			return
		}
	})
	return
}

//获取排行榜数据，只有自己
func GetRankScoreMy(rkname, userkey string, day time.Time) (my *RankModel, err error) {
	rdmd := RankEx.redismd.GetConn()
	util.Using(rdmd, func() {
		if rkcfmd, ok := RankEx.rankconflist[rkname]; ok {
			//找到配置
			key := ""
			if rkcfmd.RkCycle == -1 && rkcfmd.RkValidity == -1 {
				//永久的与时间无关
				key = rkcfmd.GetRedisKey()

			} else if rkcfmd.RkCycle == -1 {
				//没有周期
				tmday := rkcfmd.GetDayNum(day)
				tmday = tmday % rkcfmd.RkValidity
				key = rkcfmd.GetRedisKeyByDay(tmday)
			} else {
				//有周期
				tmday := rkcfmd.GetDayCycleNum(day)
				// tmday := int(day.Unix() / 86400)
				// tmday = tmday / rkcfmd.RkCycle
				tmday = tmday % rkcfmd.RkValidity
				key = rkcfmd.GetRedisKeyByDay(tmday)
			}
			switch rkcfmd.RkType {
			case RkTypeEnum_MinMax, RkTypeEnum_MinMax_TimeRev:
				my, err = getRankScoreMyByRange(rdmd, rkcfmd, key, userkey)
			case RkTypeEnum_MaxMin, RkTypeEnum_MaxMin_TimeRev:
				my, err = getRankScoreMyByRevRange(rdmd, rkcfmd, key, userkey)
			}
			return
		} else {
			err = errors.Err_Rank_Config
			return
		}
	})
	return
}

//获取排行榜数据，包括自己和列表
func GetRankScoreMyToXByCycle(rkname, userkey string, season RankDayEnum, x int) (result []*RankModel, my *RankModel, err error) {
	if rkcfmd, ok := RankEx.rankconflist[rkname]; ok {
		day := time.Now()
		if season == RankDayEnum_PreSeason {
			if rkcfmd.RkCycle != -1 {
				//改成支持周
				day = day.Add(-rkcfmd.GetCycleDuration())
				// day = day.Add(-time.Duration(rkcfmd.RkCycle) * time.Hour * 24)
			}
		}
		result, my, err = GetRankScoreMyToX(rkname, userkey, day, x)
		return
	} else {
		err = errors.Err_Rank_Config
		return
	}
}

//获取周期排行榜列表数据
func GetRankScoreToXMyByCycle(rkname string, season RankDayEnum, x int) (result []*RankModel, err error) {
	if rkcfmd, ok := RankEx.rankconflist[rkname]; ok {
		day := time.Now()
		if season == RankDayEnum_PreSeason {
			if rkcfmd.RkCycle != -1 {
				day = day.Add(-rkcfmd.GetCycleDuration())
				// day = day.Add(-time.Duration(rkcfmd.RkCycle) * time.Hour * 24)
			}
		}
		result, err = GetRankScoreToX(rkname, day, x)
		return
	} else {
		err = errors.Err_Rank_Config
		return
	}
}

//获取周期排行榜自己的数据
func GetRankScoreMyByCycle(rkname, userkey string, season RankDayEnum) (my *RankModel, err error) {
	if rkcfmd, ok := RankEx.rankconflist[rkname]; ok {
		day := time.Now()
		if season == RankDayEnum_PreSeason {
			if rkcfmd.RkCycle != -1 {
				day = day.Add(-rkcfmd.GetCycleDuration())
				// day = day.Add(-time.Duration(rkcfmd.RkCycle) * time.Hour * 24)
			}
		}
		my, err = GetRankScoreMy(rkname, userkey, day)
		return
	} else {
		err = errors.Err_Rank_Config
		return
	}
}
