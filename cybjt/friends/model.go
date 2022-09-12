package friends

import (
	"encoding/json"
	"time"

	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/friend"
)

type FriendMD struct {
	friend.FriendMD //好友系统基类

	SendTime time.Time `json:"sdt"`    //送礼物的时间
	IsGift   bool      `json:"isgift"` //是否收到的礼物
	// Gift            common.ItemInfos `json:"gift"` //收到的礼物

}

func NewFriendMD(key string) friend.IFriend {
	result := new(FriendMD)
	result.UserKey = key
	result.SendTime = common.Time1970
	result.IsGift = false
	return result
}
func ScanFriendMD(buf []byte) (result map[string]friend.IFriend) {
	result = make(map[string]friend.IFriend)
	tmpli := make(map[string]*FriendMD)
	json.Unmarshal(buf, &tmpli)
	for k, v := range tmpli {
		result[k] = v
	}
	return result
}

//礼物结构
type GiftModel struct {
	UserKey string `json:"key"` //好友ID
}

func (md *GiftModel) String() string {
	buf, _ := json.Marshal(md)
	return string(buf)
}
