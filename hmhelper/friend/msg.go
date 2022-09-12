package friend

import "wawaji_pub/hmhelper/redislib"

//好友消息号
type FriendMsgIDEnum int

const (
	//消息默认类型，不使用
	FriendMsgIDEnum_Defualt FriendMsgIDEnum = iota
	//添加好友
	FriendMsgIDEnum_AddFriend
	//添加好友确认
	FriendMsgIDEnum_AddFriendReply
	//删好友
	FriendMsgIDEnum_DelFriend
)
const (
	//好友信息
	RedisKey_FriendInfo = "FriendInfo_%s"
	//添加好友redis key
	RedisKey_AddFriend = "AddFriend_%s"
	//添加好友确认redis key
	RedisKey_AddFriendReply = "AddFriendReply_%s"
	//删好友redis key
	RedisKey_DelFriend = "DelFriend_%s"
)

//好友系统消息
type FriendMsg struct {
	MsgID  FriendMsgIDEnum //消息号
	OpenID string          //用户ID
	Data   IMsgData        //数据
}

func NewFriendMsg(msgid FriendMsgIDEnum, openid string, data IMsgData) *FriendMsg {
	result := new(FriendMsg)
	result.MsgID = msgid
	result.OpenID = openid
	result.Data = data
	return result
}

type MsgHandler func(rdmd *redislib.RedisHandleModel, msg *FriendMsg)

type GetRedisData func(rdmd *redislib.RedisHandleModel, user IUser) (result []IMsgData)

type IMsgData interface {
	//得到字符串
	String() string
}
