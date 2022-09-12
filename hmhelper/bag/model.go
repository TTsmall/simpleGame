package bag

import (
	"encoding/json"
	"fmt"
)

//背包中用到的类型

type ResultItem struct {
	ItemID int //物品ID
	Num    int //变化后的量
	Delta  int //变化的量
}

//拿到物品ID
func (md *ResultItem) GetItemID() int {
	return md.ItemID
}

//BagExItemData 基础仓库数据
type BagExItemData struct {
	Data map[int]IBagExItem
}

//NewBagExItemDataString 用字符串初始化一个数据
func NewBagExItemData() *BagExItemData {
	result := new(BagExItemData)
	result.Data = make(map[int]IBagExItem)

	return result
}

//GetNumByKey指定数据的值
func (this *BagExItemData) GetNumByItemID(itemid int) (result int) {
	for _, v := range this.Data {
		if v.GetItemID() == itemid {
			result++
		}
	}
	return result
}

//ToString 字符串化
func (this *BagExItemData) ToString() string {
	if result, err := json.Marshal(this.Data); err != nil {
		fmt.Print("BagExItemData.ToString ", err)
	} else {
		return string(result)
	}
	return ""
}

//Count 总数量
func (this *BagExItemData) Count() (result int) {
	return len(this.Data)
}

//Clear清数据
func (this *BagExItemData) Clear() {
	this.Data = make(map[int]IBagExItem)
}

//移除一个实例
func (this *BagExItemData) Remove(uid int) IBagExItem {
	result, ok := this.Data[uid]
	if ok {
		return result
	} else {
		return nil
	}
}

//更新实例
func (this *BagExItemData) Update(item IBagExItem) {
	this.Data[item.GetUid()] = item
}
