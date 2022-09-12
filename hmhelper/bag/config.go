package bag

import "wawaji_pub/hmhelper/common"

//物品配置表
type ItemCfg struct {
	Id          int    `col:"id" client:"id"`                   //物品ID
	BagKey      string `col:"bagkey" client:"bagkey"`           //背包key
	ItemType    int    `col:"type" client:"type"`               //物品类型
	SubType     int    `col:"subtype" client:"subtype"`         //子类型
	MaxNum      int    `col:"maxNum" client:"maxNum"`           //物品上限，0为没有上限
	Info        string `col:"info" client:"info"`               //物品说明
	ChangeParam string `col:"changeparam" client:"changeparam"` //物品转换为其他时的参数

	Name      string           `col:"name" client:"name"`           //物品名字
	TitleBage int              `col:"titleBage" client:"titleBage"` //物品背包分页-前端
	Icon      string           `col:"icon" client:"icon"`           //ICON
	Color     int              `col:"color" client:"color"`         //品质
	CanSell   bool             `col:"canSell" client:"canSell"`     //是否能出售给系统
	SellGet   common.ItemInfos `col:"sellGet" client:"sellGet"`     //出售获得
	DropID    string           `col:"dropId" client:"dropId"`       //掉落途径
	Price     float64          `col:"price" client:"price"`         //代币快捷购买,货币id,价格

}

//金币日志配置
type GoldLogCfg struct {
	Id        int    `col:"id" client:"id"`               //日志ID
	CenterKey string `col:"centerkey" client:"centerkey"` //中台Key
	ClientKey string `col:"clientkey" client:"clientkey"` //客户端显示Key
	IsDaily   bool   `col:"isdaily" client:"isdaily"`     //是否记入每日
}
