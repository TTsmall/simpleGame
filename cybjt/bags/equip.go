package bags

//装备背包

var (
	BagKey_Equip string = "装备背包"
)

type EquipMD struct {
	UID    int `json:"uid"`    //主键
	ItemID int `json:"itemid"` //物品ID
}
