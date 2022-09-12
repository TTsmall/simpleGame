package friend

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/buguang01/util"
)

func NewUserFriend() *UserFriend {
	result := new(UserFriend)
	result.Items = make(map[string]IFriend)
	result.AddList = make(map[string]string)
	return result
}

//用户好友系统
type UserFriend struct {
	Items   map[string]IFriend
	AddList map[string]string
}

// 实现driver.Valuer接口
func (md UserFriend) Value() (driver.Value, error) {
	if md.Items == nil {
		md.Items = make(map[string]IFriend)
	}
	m := make(map[string]util.StringJson, 0)
	if buf, err := json.Marshal(md.Items); err != nil {
		return nil, err
	} else {
		m["Items"] = util.StringJson(buf)
	}
	if buf, err := json.Marshal(md.AddList); err != nil {
		return nil, err
	} else {
		m["AddList"] = util.StringJson(buf)
	}
	return json.Marshal(&m)
}

func (md *UserFriend) String() string {
	return fmt.Sprintf("%+v", *md)
}

// 实现sql.Scanner接口
func (md *UserFriend) Scan(val interface{}) (err error) {
	if buf, ok := val.([]byte); ok {
		if len(buf) == 0 {
			buf = []byte("{}")
		}
		m := make(map[string]util.StringJson, 0)
		json.Unmarshal(buf, &m)
		//加载
		if v, ok := m["Items"]; ok {
			md.Items = FriendEx.scanFriendf([]byte(v))
		}
		md.AddList = make(map[string]string)
		if v, ok := m["AddList"]; ok {
			return json.Unmarshal([]byte(v), &md.AddList)
		}
	}
	return
}

func (md *UserFriend) GetAddFriendList() (result []string) {
	result = make([]string, len(md.AddList))
	i := 0
	for k := range md.AddList {
		result[i] = k
		i++
	}
	return
}

func (md *UserFriend) GetFriendKeys() (result []string) {
	result = make([]string, len(md.Items))
	i := 0
	for k := range md.Items {
		result[i] = k
		i++
	}
	return
}

type IFriend interface {
	//OpenID
	GetUserKey() string
}

func NewFriendMD(key string) IFriend {
	result := new(FriendMD)
	result.UserKey = key
	return result
}

//基础方法，从数据读入
func ScanFriendMD(buf []byte) (result map[string]IFriend) {
	result = make(map[string]IFriend)
	tmpli := make(map[string]*FriendMD)
	json.Unmarshal(buf, &tmpli)
	for k, v := range tmpli {
		result[k] = v
	}
	return result
}

type FriendMD struct {
	UserKey string `json:"key"` //好友的ID

}

//OpenID
func (fmd *FriendMD) GetUserKey() string {
	return fmd.UserKey
}

//用户接口
type IUser interface {
	//用户的key
	GetUserRedisKey() string
}
