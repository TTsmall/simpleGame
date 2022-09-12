package bags

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/configs"
)

//资源背包
var (
	BagKey_Resources = "资源背包"
)

type ResBagMD struct {
	Bags map[int]*bag.BagMD
	user IResBagUser
}

// 实现driver.Valuer接口
func (md ResBagMD) Value() (driver.Value, error) {
	if md.Bags == nil {
		md.Bags = make(map[int]*bag.BagMD)
	}
	return json.Marshal(md.Bags)
}

func (md *ResBagMD) String() string {
	return fmt.Sprintf("%+v", *md)
}

// 实现sql.Scanner接口
func (md *ResBagMD) Scan(val interface{}) (err error) {
	if buf, ok := val.([]byte); ok {
		if len(buf) == 0 {
			buf = []byte("")
		}
		md.Bags = make(map[int]*bag.BagMD)
		json.Unmarshal(buf, &md.Bags)
	}

	return
}

//拿到对应城市的背包
func (md *ResBagMD) GetBagByCity(cid int) *bag.BagMD {
	var bagmd *bag.BagMD
	var ok bool
	if bagmd, ok = md.Bags[cid]; !ok {
		bagmd = bag.NewBagMD()().(*bag.BagMD)
		md.Bags[cid] = bagmd
	}
	fmt.Println("GetBagByCity:", cid)
	return bagmd
}

func (md *ResBagMD) AddItem(cf *bag.ItemCfg, itemid, num, logtype int) (result bag.IBagItem, nitems common.ItemInfos) {
	if cf == nil {
		cf = configs.GetDb().GetItemCfg(itemid)
	}
	cid := md.user.GetCurrCityID()
	bagmd := md.GetBagByCity(cid)
	fmt.Println("GetBagByCity.AddItem:", cid, "itemid:", itemid, "num:", num)

	return bagmd.AddItem(cf, itemid, num, logtype)

}

//删除物品,以ItemID
func (md *ResBagMD) DelItemByID(cf *bag.ItemCfg, itemid, num, logtype int) (result bag.IBagItem) {
	if cf == nil {
		cf = configs.GetDb().GetItemCfg(itemid)
	}
	cid := md.user.GetCurrCityID()
	bagmd := md.GetBagByCity(cid)

	return bagmd.DelItemByID(cf, itemid, num, logtype)
}

//检查道具数量，返回true表示满足，返回false表示不满足
func (md *ResBagMD) CheckItem(itemid, num int) bool {
	cid := md.user.GetCurrCityID()
	return md.CheckItemByCity(cid, itemid, num)
}

func (md *ResBagMD) CheckItemByCity(cid, itemid, num int) bool {
	bagmd := md.GetBagByCity(cid)
	return bagmd.CheckItem(itemid, num)
}

//道具ID对应的数量
func (md *ResBagMD) GetItemNum(itemid int) int {
	cid := md.user.GetCurrCityID()
	return md.GetItemNumByCity(cid, itemid)
}
func (md *ResBagMD) GetItemNumByCity(cid, itemid int) int {
	bagmd := md.GetBagByCity(cid)
	return bagmd.GetItemNum(itemid)
}

//添加资源
func (md *ResBagMD) AddRes(itemid, num int) (result int, item *bag.ResultItem) {
	cid := md.user.GetCurrCityID()
	return md.AddResByCity(cid, itemid, num)
}

func (md *ResBagMD) AddResByCity(cid, itemid, num int) (result int, item *bag.ResultItem) {
	bagmd := md.GetBagByCity(cid)

	max := md.user.GetResMax(cid)
	item = new(bag.ResultItem)
	item.ItemID = itemid
	item.Num = bagmd.GetItemNum(itemid)
	if item.Num >= max {
		return num, item
	} else if item.Num+num >= max {
		result = item.Num + num - max
		item.Delta = max - item.Num
		//因为这里把返回的数据都算完了，所以不使用方法的返回
		bagmd.AddItem(configs.GetDb().GetItemCfg(itemid), itemid, item.Delta, 0)
		item.Num = bagmd.GetItemNum(itemid)
	}
	return
}

func (md *ResBagMD) DelRes(itemid, num int) (result bag.IBagItem) {
	cid := md.user.GetCurrCityID()
	return md.DelResByCity(cid, itemid, num)
}
func (md *ResBagMD) DelResByCity(cid, itemid, num int) (result bag.IBagItem) {
	bagmd := md.GetBagByCity(cid)
	return bagmd.DelItemByID(nil, itemid, num, 0)
}

//------------------初始化

//新资源背包
func NewResBagMD(user IResBagUser) *ResBagMD {
	result := new(ResBagMD)
	result.Bags = make(map[int]*bag.BagMD)
	// result.BagMD = *bag.NewBagMD()().(*bag.BagMD)
	result.user = user
	return result
}

func (md *ResBagMD) SetBagMD(user IResBagUser) {
	md.user = user
}

//------------------外部接口

type IResBagUser interface {
	//获取当前城市ID
	GetCurrCityID() int
	//资源背包的上限值，是一个动态获取的
	GetResMax(cid int) int
}
