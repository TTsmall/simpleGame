package bag

import (
	"github.com/TTsmall/wawaji_pub_hmhelper/common"
)

//仓库接口文件

type IUser interface {
	//用户对像需要实现的接口，使用KEY，拿到对应的仓库
	GetBag(key string) (result IBag)
}

//可以进入仓库管理器进行操作的接口
type IBag interface {
	//添加物品，背包内需要实现logtype的逻辑
	AddItem(cf *ItemCfg, itemid, num, logtype int) (result IBagItem, nitems common.ItemInfos)
}

//无实例的背包
type IBagByItem interface {
	IBag
	//删除物品,以ItemID
	DelItemByID(cf *ItemCfg, itemid, num, logtype int) (result IBagItem)
	//检查数量是否足够,返回true表示满足，返回false表示不满足
	CheckItem(itemid, num int) bool
	//拿物品数量
	GetItemNum(itemid int) int
}

//有实例的背包
type IBagByEx interface {
	IBag
	//删除物品,以UID
	DelItemByUID(uid, logtype int) (result IBagExItem)
}

//在仓库中的Item接口
type IBagItem interface {
	//拿到物品ID
	GetItemID() int
}

//仓库中的带实例的接口
type IBagExItem interface {
	IBagItem
	//拿实例UID
	GetUid() int
}
