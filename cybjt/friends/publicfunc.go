package friends

import (
	"time"

	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/friend"
)

//从redis中读出礼物数据填充
func SetFriendGift(user friend.IUser, fli *friend.UserFriend) {
	giftli := friend.FriendEx.GetData(user, FriendMsgIDEnum_Gift)
	for _, imd := range giftli {
		giftmd := imd.(*GiftModel)
		if ifdmd, ok := fli.Items[giftmd.UserKey]; !ok {
			continue
		} else {
			fdmd := ifdmd.(*FriendMD)
			fdmd.IsGift = true
		}
	}
}

//把好友的礼物都收下来，会修改标记
func GetFriendGiftToItemInfo(user friend.IUser, fli *friend.UserFriend, reward common.ItemInfos) (result *common.BaseData) {
	addbc := common.NewBaseDataString("")
	for _, item := range reward {
		addbc.UpData(item.ItemId, item.Count)
	}
	SetFriendGift(user, fli)
	result = common.NewBaseDataString("")
	for _, ifdmd := range fli.Items {
		fdmd := ifdmd.(*FriendMD)
		if fdmd.IsGift {
			fdmd.IsGift = false
			result.UpDataBc(addbc, nil)
		}
	}
	return result
}

//送礼物给全好友
func SendFriendsGift(user friend.IUser, fli *friend.UserFriend) (result []string) {
	result = make([]string, 0)
	for key, ifdmd := range fli.Items {
		fdmd := ifdmd.(*FriendMD)
		if common.SameDay(time.Now(), fdmd.SendTime) {
			giftmd := new(GiftModel)
			giftmd.UserKey = user.GetUserRedisKey()
			fdmd.SendTime = time.Now()
			friend.FriendEx.SendMsg(fdmd.UserKey, FriendMsgIDEnum_Gift, giftmd)
			result = append(result, key)
		}
	}
	return
}
func SendFriendGift(user friend.IUser, fdmd *FriendMD) (result []string) {
	result = make([]string, 0)

	if common.SameDay(time.Now(), fdmd.SendTime) {
		giftmd := new(GiftModel)
		giftmd.UserKey = user.GetUserRedisKey()
		fdmd.SendTime = time.Now()
		friend.FriendEx.SendMsg(fdmd.UserKey, FriendMsgIDEnum_Gift, giftmd)
		result = append(result, fdmd.UserKey)
	}
	return
}
