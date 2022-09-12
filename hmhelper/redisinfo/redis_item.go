package redisinfo

import (
	"encoding/json"
	"fmt"

	"wawaji_pub/hmhelper/common"
	"wawaji_pub/hmhelper/helper"
	"wawaji_pub/hmhelper/log"
	"wawaji_pub/hmhelper/redislib"
)

const (
	//用户信息
	RedisKey_UserReward = "UserRewards:"
)

type RedisItem struct {
	RedisKey string `json:"-"`
	Reward   common.ItemInfos
	LogType  int
}

//用户奖励写入redis
func RedisUserRewardSet(md *RedisItem) {
	key := fmt.Sprint(RedisKey_UserReward, md.RedisKey)
	buf, _ := json.Marshal(md)
	helper.UsingRedis(func(rdmd *redislib.RedisHandleModel) {
		if _, err := rdmd.Append(key, fmt.Sprint(string(buf), ",")); err != nil {
			log.Error("RedisUserRewardSet err:", err)
		}
	})
	// ModelPoolEx.Put(user)
	//
}

//用户奖励读出redis
func RedisUserRewardGet(rediskey string) (result []*RedisItem) {
	result = make([]*RedisItem, 0, 50)
	key := fmt.Sprint(RedisKey_UserReward, rediskey)
	data := "[]"
	helper.UsingRedis(func(rdmd *redislib.RedisHandleModel) {
		v, err := rdmd.GetSet(key, "")
		if err != nil {
			log.Error("giftGet GetSet err.", err)
			return
		}
		if len(v) == 0 {
			return
		}
		data = fmt.Sprintf("[%s]", v[:len(v)-1])
	})
	json.Unmarshal([]byte(data), &result)
	return result
}
