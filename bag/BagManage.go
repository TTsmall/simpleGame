package bag

import (
	"fmt"
	"math"

	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/errorx"
)

/*
要可以从用户身上拿到需要的仓库
也就是每个仓库应该有一个自己的名字，可以通过物品表的配置拿到这个信息
拿到这个仓库后，就可以把东西放到这个仓库中
仓库对象本身应该有仓库接口可以使用
不同类型的仓库在创建的时候，需要使用不同初始化
这个仓库还要有创建存储实例的方法，方法中有初始化信息，可以传入对应的物品信息
*/

var (
	BagKey_Gold = "金币背包"
	BagKey_Base = "普通背包"
)
var (
	BagEx *BagManage
)

//初始化
func NewBagManage(helpcf *HelperConf, goldcf map[int]*GoldLogCfg, itemcf map[int]*ItemCfg) (result *BagManage) {
	result = new(BagManage)
	// result.bagconfs = make(map[string]InitBag)
	result.GameEx = helpcf
	result.GoldConfEx = goldcf
	result.ItemConfEx = itemcf
	return
}

//需要别的背包实现的新建包的方法
type InitBag func() (result IBag)

//仓库管理器
type BagManage struct {
	// bagconfs   map[string]InitBag  //背包生成器
	GameEx     *HelperConf         //基础的全局配置表
	GoldConfEx map[int]*GoldLogCfg //金币日志配置
	ItemConfEx map[int]*ItemCfg    //物品配置
}

// func (mg *BagManage) SetInit(key string, ibfunc InitBag) {
// 	mg.bagconfs[key] = ibfunc
// 	return
// }

// func (mg *BagManage) CreateBag(key string) (result IBag) {
// 	result = mg.bagconfs[key]()
// 	return
// }

//添加物品
func (mg *BagManage) AddItems(user IUser, items common.ItemInfos, logtype int) (result []IBagItem) {
	for _, item := range items {
		conf := mg.ItemConfEx[item.ItemId]
		if bag := user.GetBag(conf.BagKey); bag != nil {
			additem, resultitem := bag.AddItem(conf, item.ItemId, item.Count, logtype)
			if additem != nil {
				result = append(result, additem)
			}
			if resultitem != nil {
				fmt.Println("resultitem not nil.", resultitem)
				result = append(result, mg.AddItems(user, resultitem, logtype)...)
			}

		} else {
			//礼包检查
		}

	}
	return
}

//单物品
func (mg *BagManage) AddItem(user IUser, itemid, num int, logtype int) (result []IBagItem) {
	conf := mg.ItemConfEx[itemid]
	if bag := user.GetBag(conf.BagKey); bag != nil {
		additem, resultitem := bag.AddItem(conf, itemid, num, logtype)
		if additem != nil {
			result = append(result, additem)
		}

		if resultitem != nil {
			result = append(result, mg.AddItems(user, resultitem, logtype)...)
		}

	} else {
		//礼包检查
	}

	return
}

//通过物品ID扣物品
func (mg *BagManage) DelItems(user IUser, items common.ItemInfos, logtype int) (result []IBagItem, err error) {
	for _, item := range items {
		conf := mg.ItemConfEx[item.ItemId]
		bag := user.GetBag(conf.BagKey).(IBagByItem)
		result = append(result, bag.DelItemByID(conf, item.ItemId, item.Count, logtype))
	}
	return
}

func (mg *BagManage) DelItem(user IUser, itemid, num int, logtype int) (result []IBagItem, err error) {
	conf := mg.ItemConfEx[itemid]
	bag := user.GetBag(conf.BagKey).(IBagByItem)

	result = append(result, bag.DelItemByID(conf, itemid, num, logtype))
	return
}

//通过物品UID扣物品
func (mg *BagManage) DelItemByUID(bag IBagByEx, uids []int, logtype int) (result []IBagExItem) {
	for _, item := range uids {
		result = append(result, bag.DelItemByUID(item, logtype))
	}
	return
}

//检查物品数量
func (mg *BagManage) CheckItem(user IUser, items common.ItemInfos) (err error) {
	for _, item := range items {
		if item.Count < 0 {
			//非法的
			return errorx.Bag_Item_Illegal
		}
		conf := mg.ItemConfEx[item.ItemId]
		bag := user.GetBag(conf.BagKey).(IBagByItem)
		if !bag.CheckItem(item.ItemId, item.Count) {
			fmt.Println("CheckItem ", item.ItemId, " bagitemnum: ", bag.GetItemNum(item.ItemId), " count: ", item.Count)
			return errorx.Bag_Item_Not_Enough.CloneWithArgs(conf.Name)
		}
	}
	return
}

//拆分物品
func (mg *BagManage) SplitItem(user IUser, items common.ItemInfos) (owned, unowned *common.BaseData) {
	owned = common.NewBaseDataString("")
	unowned = common.NewBaseDataString("")
	for i := range items {
		conf := mg.ItemConfEx[items[i].ItemId]
		bag := user.GetBag(conf.BagKey).(IBagByItem)
		oldnum := bag.GetItemNum(conf.Id)
		num := items[i].Count
		if oldnum > num {
			owned.UpData(conf.Id, num)
		} else {
			owned.UpData(conf.Id, oldnum)
			unowned.UpData(conf.Id, num-oldnum)
		}
	}
	return
}

//物品转成金币
func (mg *BagManage) GoldReplace(user IUser, unowned *common.BaseData) (result *common.BaseData) {
	result = common.NewBaseDataString("")
	for k, n := range unowned.Data {
		conf := mg.ItemConfEx[k]
		num := int(math.Ceil(conf.Price * float64(n)))
		result.UpData(ItemType_Gold, num)
	}
	return
}

//-----------------拿背包配置
