package friend

import (
	"encoding/json"
	"fmt"

	"github.com/buguang01/util"
	"wawaji_pub/hmhelper/log"
	"wawaji_pub/hmhelper/redislib"
)

func anyFriendMsgHandler(rdmd *redislib.RedisHandleModel, msg *FriendMsg, redis_key string) {
	key := fmt.Sprintf(redis_key, msg.OpenID)
	num, err := rdmd.Append(key, fmt.Sprint(msg.Data.String(), ","))
	if err != nil {
		log.Error("GiftMsgHandler append err.", err)
		return
	}
	max := 4000
	if num > max {
		v, err := rdmd.GetSet(key, "")
		if err != nil {
			log.Error("GiftMsgHandler GetSet err.", err)
			return
		}
		if len(v) == 0 {
			return
		}
		v = fmt.Sprintf("[%s]", v[:len(v)-1])
		li := make([]*FriendMsgModel, 0, 50)
		json.Unmarshal([]byte(v), &li)
		tmpmap := make(map[string]*FriendMsgModel)
		for _, item := range li {
			tmpmap[item.UserKey] = item
		}
		sli := util.NewStringBuilder()
		for _, item := range tmpmap {
			sli.Append(item.String())
			sli.Append(",")
		}
		rdmd.Append(key, sli.ToString())
	}
}

func anyFriendGet(rdmd *redislib.RedisHandleModel, user IUser, redis_key string) (result []IMsgData) {
	result = make([]IMsgData, 0, 50)
	key := fmt.Sprintf(redis_key, user.GetUserRedisKey())
	v, err := rdmd.GetSet(key, "")
	if err != nil {
		log.Error("giftGet GetSet err.", err)
		return
	}
	if len(v) == 0 {
		return
	}
	v = fmt.Sprintf("[%s]", v[:len(v)-1])
	li := make([]*FriendMsgModel, 0, 50)
	json.Unmarshal([]byte(v), &li)
	for i := range li {
		result = append(result, li[i])
	}
	return result
}

func addFriendMsgHandler(rdmd *redislib.RedisHandleModel, msg *FriendMsg) {
	anyFriendMsgHandler(rdmd, msg, RedisKey_AddFriend)
}

func addFriendGet(rdmd *redislib.RedisHandleModel, user IUser) (result []IMsgData) {
	return anyFriendGet(rdmd, user, RedisKey_AddFriend)
}
func addFriendReplyMsgHandler(rdmd *redislib.RedisHandleModel, msg *FriendMsg) {
	anyFriendMsgHandler(rdmd, msg, RedisKey_AddFriendReply)
}

func addFriendReplyGet(rdmd *redislib.RedisHandleModel, user IUser) (result []IMsgData) {
	return anyFriendGet(rdmd, user, RedisKey_AddFriendReply)
}

func delFriendMsgHandler(rdmd *redislib.RedisHandleModel, msg *FriendMsg) {
	anyFriendMsgHandler(rdmd, msg, RedisKey_DelFriend)
}

func delFriendGet(rdmd *redislib.RedisHandleModel, user IUser) (result []IMsgData) {
	return anyFriendGet(rdmd, user, RedisKey_DelFriend)
}
