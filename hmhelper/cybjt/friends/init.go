package friends

import (
	"encoding/json"
	"fmt"

	"github.com/buguang01/util"
	"wawaji_pub/hmhelper/friend"
	"wawaji_pub/hmhelper/log"
	"wawaji_pub/hmhelper/redislib"
)

//
func NewFriendManager() *friend.FriendManager {
	result := friend.NewFriendManager(
		friend.SetFriendMgByNewf(NewFriendMD),
		friend.SetFriendMgByScanf(ScanFriendMD),
		friend.SetFriendMgByMsgID(FriendMsgIDEnum_Gift, giftMsgHandler, giftGet),
	)
	friend.FriendEx = result
	return result
}

const (
	//发礼物
	FriendMsgIDEnum_Gift friend.FriendMsgIDEnum = 100 + iota
	//好友商店
	FriendMsgIDEnum_Shop
)

const (
	//发礼物的redis key
	RedisKey_Gift string = "bjtGift_%s"
)

//发礼物
func giftMsgHandler(rdmd *redislib.RedisHandleModel, msg *friend.FriendMsg) {
	key := fmt.Sprintf(RedisKey_Gift, msg.OpenID)
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
		li := make([]*GiftModel, 0, 50)
		json.Unmarshal([]byte(v), &li)
		tmpmap := make(map[string]*GiftModel)
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

//礼物的获取
func giftGet(rdmd *redislib.RedisHandleModel, user friend.IUser) (result []friend.IMsgData) {
	result = make([]friend.IMsgData, 0, 50)
	key := fmt.Sprintf(RedisKey_Gift, user.GetUserRedisKey())
	v, err := rdmd.GetSet(key, "")
	if err != nil {
		log.Error("giftGet GetSet err.", err)
		return
	}
	if len(v) == 0 {
		return
	}
	v = fmt.Sprintf("[%s]", v[:len(v)-1])
	li := make([]*GiftModel, 0, 50)
	json.Unmarshal([]byte(v), &li)
	for i := range li {
		result = append(result, li[i])
	}
	return result
}
