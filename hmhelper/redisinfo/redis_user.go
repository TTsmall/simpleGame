package redisinfo

import (
	"encoding/json"
	"sync"
	"time"

	"wawaji_pub/hmhelper/helper"
	"wawaji_pub/hmhelper/log"
	"wawaji_pub/hmhelper/redislib"
)

//写入redis的用户信息

const (
	//用户信息
	RedisKey_UserInfo = "UserInfos"
)

var (
	//用户信息模型的池子
	ModelPoolEx *sync.Pool
)

type ModelPoolOption func(md *sync.Pool)

func SetModelPoolFnew(fnew func() IModel) ModelPoolOption {
	return func(md *sync.Pool) {
		md.New = func() interface{} {
			return fnew()
		}
	}
}
func NewModelPool(opts ...ModelPoolOption) {
	ModelPoolEx = new(sync.Pool)
	ModelPoolEx.New = func() interface{} {
		return NewUserRedisModel()
	}
	for i := range opts {
		opts[i](ModelPoolEx)
	}
}

type IModel interface {
	//用户的key
	GetRedisKey() string
	//用户redis的数据
	GetUserRedisValue() string
	// //清数据
	// Clear()
}

//用户信息写入redis
func RedisUserInfoSet(user IModel) {
	helper.UsingRedis(func(rdmd *redislib.RedisHandleModel) {
		if _, err := rdmd.Hmset(RedisKey_UserInfo, user.GetRedisKey(), user.GetUserRedisValue()); err != nil {
			log.Error("RedisUserInfoSet err:", err)
		}
	})
	// ModelPoolEx.Put(user)
	//
}

//用户信息读出redis
func RedisUserInfoGet(key string) (result string, err error) {
	// restr := ""
	helper.UsingRedis(func(rdmd *redislib.RedisHandleModel) {
		if result, err = rdmd.Hget(RedisKey_UserInfo, key); err != nil {
			log.Error("RedisUserInfoGet err:", err)
		}
	})
	// result = ModelPoolEx.Get().(IModel)
	// result.Clear()
	// if err = json.Unmarshal([]byte(restr), result); err != nil {
	// 	log.Error("RedisUserInfoGet err:", err)
	// }
	return
}

//用户信息列表读出redis
func RedisUserInfoGets(keys []string) (result map[string]string, err error) {
	helper.UsingRedis(func(rdmd *redislib.RedisHandleModel) {
		if result, err = rdmd.Hmget(RedisKey_UserInfo, keys...); err != nil {
			log.Error("RedisUserInfoGets err:", err)
		}
	})
	return
}

func NewUserRedisModel() IModel {
	result := new(UserRedisModel)
	return result
}

type UserRedisModel struct {
	UserKey string    `json:"UserKey,omitempty"` //主键
	Name    string    `json:"name,omitempty"`    //名字
	Stage   int       `json:"stage,omitempty"`   //关卡
	Level   int       `json:"level,omitempty"`   //等级
	Avatar  string    `json:"avatar,omitempty"`  //头像
	Online  bool      `json:"online,omitempty"`  //是否在线
	UpTime  time.Time `json:"uptime,moitempty"`  //更新时间
}

//用户的key
func (md *UserRedisModel) GetRedisKey() string {
	return md.UserKey
}

//用户redis的数据
func (md *UserRedisModel) GetUserRedisValue() string {
	buf, _ := json.Marshal(md)
	return string(buf)
}

func (md *UserRedisModel) Clear() {
	md.UserKey = ""
	md.Name = ""
	md.Stage = 0
	md.Level = 0
	md.Avatar = ""

}
