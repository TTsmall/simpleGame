package configs

type ItemTypeEnum int

const (
	//加速卡
	ItemTypeEnum_QuickCard ItemTypeEnum = 27
	ItemTypeEnum_Happines  ItemTypeEnum = 11
)

var (
	//加速卡
	ItemType_QuickCard int = int(ItemTypeEnum_QuickCard)
	ItemType_Happines  int = int(ItemTypeEnum_Happines)
)

//SetCountItem 设置道具的ID，如果上层道具ID冲突了，可以通过这个方法设置
func SetCountItem(itemid ItemTypeEnum, newitem int) {
	switch ItemTypeEnum(itemid) {
	case ItemTypeEnum_QuickCard:
		ItemType_QuickCard = newitem
	case ItemTypeEnum_Happines:
		ItemType_Happines = newitem
	}
}
