package bags

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"wawaji_pub/hmhelper/common"
)

//平民
// var (
// 	BagKey_People string = "平民背包"
// )

type PeopleMD struct {
	PeopleID int              `json:"uid"`             //平民主键
	CityID   int              `json:"ctid"`            //城市ID
	Status   HeroStatusEnum   `json:"stid,omitempty"`  //状态
	EventLi  common.ItemInfos `json:"evtli,omitempty"` //事件ID
}

func NewPeople(uid int) *PeopleMD {
	result := new(PeopleMD)
	result.PeopleID = uid
	result.Status = HeroStatus_Default
	result.EventLi = make(common.ItemInfos, 0, 10)
	return result
}

func (md *PeopleMD) SetEventID(eid, status int) {
	var item *common.ItemInfo
	if item = md.EventLi.Get(eid); item == nil {
		item = new(common.ItemInfo)
		item.ItemId = eid
		md.EventLi = append(md.EventLi, item)
	}
	item.Count = status
}

//平民背包，不走背包接口
type PeopleBagMD struct {
	Items map[int]*PeopleMD
}

// 实现driver.Valuer接口
func (md PeopleBagMD) Value() (driver.Value, error) {
	if md.Items == nil {
		md.Items = make(map[int]*PeopleMD)
	}
	return json.Marshal(md.Items)
}

func (md *PeopleBagMD) String() string {
	return fmt.Sprintf("%+v", *md)
}

// 实现sql.Scanner接口
func (md *PeopleBagMD) Scan(val interface{}) (err error) {
	if buf, ok := val.([]byte); ok {
		if len(buf) == 0 {
			buf = []byte("{}")
		}
		err = json.Unmarshal(buf, &md.Items)
	}
	return
}

//地图上插入新的NPC
func (md *PeopleBagMD) Insert(uid int) *PeopleMD {
	if result, ok := md.Items[uid]; ok {
		return result
	} else {
		result := NewPeople(uid)
		md.Items[uid] = result
		return result
	}

}

//拿主键获取信息
func (md *PeopleBagMD) GetPeopleByUID(uid int) (result *PeopleMD, ok bool) {
	result, ok = md.Items[uid]
	return result, ok
}

func (md *PeopleBagMD) GetPeoplesByUID(uids []int) []*PeopleMD {
	results := make([]*PeopleMD, 0)
	for _, uid := range uids {
		if result, ok := md.Items[uid]; ok {
			results = append(results, result)
		}
	}

	return results
}

//拿状态的NPC,指定返回的数量,返回时，只会返回少于等于指定数量的数据
func (md *PeopleBagMD) GetPeopleByStatusID(sid HeroStatusEnum, cityid, num int) ([]*PeopleMD, bool) {
	result := make([]*PeopleMD, 0, num)
	for _, v := range md.Items {
		if v.Status == sid { //&& v.CityID == cityid {
			result = append(result, v)
			if len(result) >= num {
				break
			}
		}
	}
	return result, len(result) >= num
}

//------------------初始化
//创建英雄背包
func NewNpcHeroBag() *PeopleBagMD {
	result := new(PeopleBagMD)
	result.Items = make(map[int]*PeopleMD)

	return result
}

func SetNpcHeroBag(md *PeopleBagMD) {
	//还没有需要设置的东西
}

//---------外部接口

//平民背包，用户接口
type IPeopleUser interface {
	//获取平民背包
	GetPeopleBag() *PeopleBagMD
}
