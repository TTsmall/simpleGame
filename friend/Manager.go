package friend

import (
	"context"
	"fmt"

	"github.com/TTsmall/wawaji_pub_hmhelper/log"
	"github.com/TTsmall/wawaji_pub_hmhelper/redislib"
	"github.com/buguang01/util"
	"github.com/buguang01/util/threads"
)

//好友系统

/*
通过redis进行用户数据交互
所有用户都通过管理器进行异步数据的写入
在写入的时候，会有一个数据长度的判断，过长了就会进行数据整理
从redis读出后，就会从redis清空，所以要把数据写入DB
*/
var (
	FriendEx *FriendManager
)

type NewFriendFunc (func(key string) IFriend)

type ScanFriendFunc func(buf []byte) (result map[string]IFriend)

type FriendOption func(md *FriendManager)

func SetFriendMgByNewf(f NewFriendFunc) FriendOption {
	return func(md *FriendManager) {
		md.nFriendf = f
	}
}

func SetFriendMgByScanf(f ScanFriendFunc) FriendOption {
	return func(md *FriendManager) {
		md.scanFriendf = f
	}
}

func SetFriendMgByMsgID(msgid FriendMsgIDEnum, f MsgHandler, g GetRedisData) FriendOption {
	return func(md *FriendManager) {
		md.msgHandlerList[msgid] = f
		md.getmsgList[msgid] = g
	}
}

func NewFriendManager(opts ...FriendOption) *FriendManager {
	md := new(FriendManager)
	md.nFriendf = NewFriendMD
	md.scanFriendf = ScanFriendMD
	md.msgChan = make(chan *FriendMsg, 1024)
	md.msgHandlerList = make(map[FriendMsgIDEnum]MsgHandler)
	md.getmsgList = make(map[FriendMsgIDEnum]GetRedisData)
	md.msgHandlerList[FriendMsgIDEnum_AddFriend] = addFriendMsgHandler
	md.getmsgList[FriendMsgIDEnum_AddFriend] = addFriendGet
	md.msgHandlerList[FriendMsgIDEnum_AddFriendReply] = addFriendReplyMsgHandler
	md.getmsgList[FriendMsgIDEnum_AddFriendReply] = addFriendReplyGet
	md.msgHandlerList[FriendMsgIDEnum_DelFriend] = delFriendMsgHandler
	md.getmsgList[FriendMsgIDEnum_DelFriend] = delFriendGet
	for i := range opts {
		opts[i](md)
	}
	return md
}

type FriendManager struct {
	nFriendf       NewFriendFunc                    //新建好友实例
	scanFriendf    ScanFriendFunc                   //DB新建好友实例
	msgChan        chan *FriendMsg                  //消息队列
	redismd        *redislib.RedisAccess            //redis
	msgHandlerList map[FriendMsgIDEnum]MsgHandler   //消息列表
	getmsgList     map[FriendMsgIDEnum]GetRedisData //获取数据
	thg            *threads.ThreadGo
}

func (this *FriendManager) Init() error {
	this.thg = threads.NewThreadGo()
	return nil
}

func (this *FriendManager) Run() {
	this.thg.Go(this.handler)
}

func (this *FriendManager) Stop() {
	this.thg.CloseWait()
}

func (this *FriendManager) handler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(this.msgChan)
		case msg, ok := <-this.msgChan:
			if !ok {
				return
			}
			if f, ok := this.msgHandlerList[msg.MsgID]; !ok {
				log.Info("friend msg not exist id:%d", msg.MsgID)
			} else {
				rdmd := this.redismd.GetConn()
				util.Using(rdmd, func() {
					f(rdmd, msg)
				})
			}
		}
	}
}

//写入redis的数据
func (this *FriendManager) SendMsg(openid string, msgid FriendMsgIDEnum, data IMsgData) {
	msg := NewFriendMsg(msgid, openid, data)
	this.msgChan <- msg
}

//从redis中读出自定义的数据
func (this *FriendManager) GetData(user IUser, msgid FriendMsgIDEnum) (result []IMsgData) {
	g := this.getmsgList[msgid]
	rdmd := this.redismd.GetConn()
	util.Using(rdmd, func() {
		result = g(rdmd, user)
	})
	return
}

//加好友
func (this *FriendManager) FriendAdd(user IUser, otheropenid string) {
	fdata := new(FriendMsgModel)
	fdata.UserKey = user.GetUserRedisKey()
	this.SendMsg(otheropenid, FriendMsgIDEnum_AddFriend, fdata)
}

//加好友确认
func (this *FriendManager) FriendAddReply(user IUser, otheropenid string) {
	fdata := new(FriendMsgModel)
	fdata.UserKey = user.GetUserRedisKey()
	this.SendMsg(otheropenid, FriendMsgIDEnum_AddFriendReply, fdata)
}

//删好友
func (this *FriendManager) FriendDel(user IUser, otheropenid string) {
	fdata := new(FriendMsgModel)
	fdata.UserKey = user.GetUserRedisKey()
	this.SendMsg(otheropenid, FriendMsgIDEnum_DelFriend, fdata)
}

//获取好友访问数据
func (this *FriendManager) GetFriendInfo(openid string) (result string) {
	key := fmt.Sprintf(RedisKey_FriendInfo, openid)
	rdmd := this.redismd.GetConn()
	util.Using(rdmd, func() {
		if data, err := rdmd.Get(key); err != nil {
			log.Error("GetFriendInfo err.", err)
		} else {
			result = data
		}
	})
	return
}

//设置好友访问数据
func (this *FriendManager) SetFriendInfo(openid, data string) {
	key := fmt.Sprintf(RedisKey_FriendInfo, openid)
	rdmd := this.redismd.GetConn()
	util.Using(rdmd, func() {
		if _, err := rdmd.Set(key, data, -1, redislib.Set_No_NX_XX); err != nil {
			log.Error("SetFriendInfo err.", err)
		}
	})
	return
}

//设置好友添加信息
func (this *FriendManager) SetUserFriend(user IUser, fli *UserFriend) {
	addlist := this.GetData(user, FriendMsgIDEnum_AddFriend)
	addfriend := this.GetData(user, FriendMsgIDEnum_AddFriendReply)
	delfriend := this.GetData(user, FriendMsgIDEnum_DelFriend)
	//申请列表
	for _, iv := range addlist {
		md := iv.(*FriendMsgModel)
		fli.AddList[md.UserKey] = "1"
	}
	//好友列表
	for _, iv := range addfriend {
		md := iv.(*FriendMsgModel)
		if ifmd, ok := fli.Items[md.UserKey]; !ok {
			ifmd = this.nFriendf(md.UserKey)
			fli.Items[ifmd.GetUserKey()] = ifmd
		}
		delete(fli.AddList, md.UserKey)
	}
	//删除好友列表
	for _, iv := range delfriend {
		md := iv.(*FriendMsgModel)
		delete(fli.Items, md.UserKey)
	}
}
