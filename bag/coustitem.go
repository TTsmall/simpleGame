package bag

type ItemTypeEnum int

const (
	//物品ID，金币
	ItemTypeEnum_Gold ItemTypeEnum = 1
)

var (
	//物品ID，金币
	ItemType_Gold int = int(ItemTypeEnum_Gold)
)

//SetCountItem 设置道具的ID，如果上层道具ID冲突了，可以通过这个方法设置
func SetCountItem(itemid ItemTypeEnum, newitem int) {
	switch ItemTypeEnum(itemid) {
	case ItemTypeEnum_Gold:
		ItemType_Gold = newitem
	}
}
