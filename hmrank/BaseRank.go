package hmrank

import (
	"encoding/json"
	"fmt"

	"wawaji/hmhelper/redislib"
)

//一般排行榜命令

//写入排行榜
//获取排行前N名
//获取排行榜自己的排名

type RankDayEnum int

const (
	RankDayEnum_CurrSeason RankDayEnum = iota //当前周期
	RankDayEnum_PreSeason                     //上一周期

)

//从小到大
func getRankScoreToXByRange(rdmd *redislib.RedisHandleModel, rkcfmd *RankConfig, key string, start, end int) (result []*RankModel, err error) {
	result = make([]*RankModel, 0, end-start)
	rdkey := fmt.Sprint(RankEx.gamename, "_", key)
	rdinfokey := fmt.Sprint(RankEx.gamename, "_info_", key)
	if vals, rderr := rdmd.Zrange(rdkey, start, end); rderr != nil {
		err = rderr
		return
	} else {
		scoreli, _ := rdmd.Hmget(rdinfokey, vals...)
		for index, ukey := range vals {
			rkmd := new(RankModel)
			result = append(result, rkmd)
			rkmd.RankNo = start + index + 1
			rkmd.UserKey = ukey
			rkmd.Score, _ = scoreli[ukey]

		}
	}
	return
}

//从大到小
func getRankScoreToXByRevRange(rdmd *redislib.RedisHandleModel, rkcfmd *RankConfig, key string, start, end int) (result []*RankModel, err error) {
	result = make([]*RankModel, 0, end-start)
	rdkey := fmt.Sprint(RankEx.gamename, "_", key)
	rdinfokey := fmt.Sprint(RankEx.gamename, "_info_", key)
	if vals, rderr := rdmd.ZrevRange(rdkey, start, end); rderr != nil {
		err = rderr
		return
	} else {
		scoreli, _ := rdmd.Hmget(rdinfokey, vals...)
		for index, ukey := range vals {
			rkmd := new(RankModel)
			result = append(result, rkmd)
			rkmd.RankNo = start + index + 1
			rkmd.UserKey = ukey
			rkmd.Score, _ = scoreli[ukey]

		}
	}
	return
}

//从小到大
func getRankScoreToXMyByRange(rdmd *redislib.RedisHandleModel, rkcfmd *RankConfig, key, userkey string, start, end int) (result []*RankModel, my *RankModel, err error) {
	result = make([]*RankModel, 0, end-start)
	rdkey := fmt.Sprint(RankEx.gamename, "_", key)
	rdinfokey := fmt.Sprint(RankEx.gamename, "_info_", key)
	if vals, rderr := rdmd.Zrange(rdkey, start, end); rderr != nil {
		err = rderr
		return
	} else {
		scoreli, _ := rdmd.Hmget(rdinfokey, vals...)
		for index, ukey := range vals {
			rkmd := new(RankModel)
			result = append(result, rkmd)
			rkmd.RankNo = start + index + 1
			rkmd.UserKey = ukey
			rkmd.Score, _ = scoreli[ukey]
			if my == nil && ukey == userkey {
				my = rkmd
			}
		}
	}
	if my == nil {
		my = new(RankModel)
		my.UserKey = userkey
		if rankno, rderr := rdmd.Zrank(rdkey, userkey); rderr == nil {
			my.RankNo = rankno + 1
		} else {
			my.RankNo = -1
		}
		my.Score, _ = rdmd.Hget(rdinfokey, userkey)
		if my.Score == "" {
			my.Score = "0"
		}
	}
	return
}

//从大到小
func getRankScoreToXMyByRevRange(rdmd *redislib.RedisHandleModel, rkcfmd *RankConfig, key, userkey string, start, end int) (result []*RankModel, my *RankModel, err error) {
	result = make([]*RankModel, 0, end-start)
	rdkey := fmt.Sprint(RankEx.gamename, "_", key)
	rdinfokey := fmt.Sprint(RankEx.gamename, "_info_", key)
	if vals, rderr := rdmd.ZrevRange(rdkey, start, end); rderr != nil {
		err = rderr
		return
	} else {
		scoreli, _ := rdmd.Hmget(rdinfokey, vals...)
		b, _ := json.Marshal(scoreli)
		fmt.Println(string(b))
		for index, ukey := range vals {
			rkmd := new(RankModel)
			result = append(result, rkmd)
			rkmd.RankNo = start + index + 1
			rkmd.UserKey = ukey
			rkmd.Score, _ = scoreli[ukey]
			if my == nil && ukey == userkey {
				my = rkmd
			}
		}
	}
	if my == nil {
		my = new(RankModel)
		my.UserKey = userkey
		if rankno, rderr := rdmd.ZrevRank(rdkey, userkey); rderr == nil {
			my.RankNo = rankno + 1
		} else {
			my.RankNo = -1
		}
		my.Score, _ = rdmd.Hget(rdinfokey, userkey)
		if my.Score == "" {
			my.Score = "0"
		}
	}
	return
}

//从小到大 只拿自己
func getRankScoreMyByRange(rdmd *redislib.RedisHandleModel, rkcfmd *RankConfig, key, userkey string) (my *RankModel, err error) {
	rdkey := fmt.Sprint(RankEx.gamename, "_", key)
	rdinfokey := fmt.Sprint(RankEx.gamename, "_info_", key)
	my = new(RankModel)
	my.UserKey = userkey
	if rankno, rderr := rdmd.Zrank(rdkey, userkey); rderr == nil {
		my.RankNo = rankno + 1
	} else {
		my.RankNo = -1
	}
	my.Score, _ = rdmd.Hget(rdinfokey, userkey)
	if my.Score == "" {
		my.Score = "0"
	}
	return
}

//从大到小 只拿自己
func getRankScoreMyByRevRange(rdmd *redislib.RedisHandleModel, rkcfmd *RankConfig, key, userkey string) (my *RankModel, err error) {
	rdkey := fmt.Sprint(RankEx.gamename, "_", key)
	rdinfokey := fmt.Sprint(RankEx.gamename, "_info_", key)
	my = new(RankModel)
	my.UserKey = userkey
	if rankno, rderr := rdmd.ZrevRank(rdkey, userkey); rderr == nil {
		my.RankNo = rankno + 1
	} else {
		my.RankNo = -1
	}
	my.Score, _ = rdmd.Hget(rdinfokey, userkey)
	if my.Score == "" {
		my.Score = "0"
	}
	return
}
