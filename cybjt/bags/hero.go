package bags

import (
	"encoding/json"

	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/configs"
)

var (
	BagKey_Hero string = "英雄背包"
)

//--------------初始化英雄背包

func NewHeroItem(uid int, conf *bag.ItemCfg) bag.IBagExItem {
	cfg := configs.GetDb().HeroCfgs[conf.Id]

	md := new(HeroMD)
	md.HeroID = cfg.Id
	md.CityID = 0
	md.IsHome = false
	md.EquipUID = 0
	md.HeroLv = 1
	md.HeroExp = 0
	md.StarID = cfg.MinStar
	md.Status = HeroStatus_Default
	return md
}

func ExistHeroChange(cf *bag.ItemCfg, item bag.IBagExItem) common.ItemInfos {
	if herocf, ok := configs.GetDb().HeroCfgs[cf.Id]; ok {
		return append(common.ItemInfos{}, herocf.RecoveryReward...)
	}
	return nil
}

//创建英雄背包
func NewHeroBag() *bag.BagExMD {
	result := bag.NewBagExMD(
		bag.BagSetNewItemF(NewHeroItem),
		bag.BagSetExistChange(ExistHeroChange),
	)
	return result().(*bag.BagExMD)
}
func SetHeroBag(md *bag.BagExMD) {
	md.SetBag(
		bag.BagSetNewItemF(NewHeroItem),
		bag.BagSetExistChange(ExistHeroChange),
		bag.BagSetUnmarshalF(HeroUnmarshal),
	)

}

func HeroUnmarshal(buf []byte) (result map[int]bag.IBagExItem) {
	result = make(map[int]bag.IBagExItem)
	tmpli := make(map[int]*HeroMD)
	json.Unmarshal(buf, &tmpli)
	for k, v := range tmpli {
		result[k] = v
	}
	return
}
