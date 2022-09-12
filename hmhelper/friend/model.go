package friend

import "encoding/json"

type FriendMsgModel struct {
	UserKey string `json:"key"` //好友ID
}

func (md *FriendMsgModel) String() string {
	buf, _ := json.Marshal(md)
	return string(buf)
}
