package bags

import (
	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/configs"
)

//建筑背包
var (
	BagKey_Building = "建筑背包"
)

//------------------初始化

func NewBuilding() *bag.BagMD {
	result := bag.NewBagMD(
		bag.BagSetTotalMax(configs.GetDb().GameEx.BuildingBagNumMax),
	)().(*bag.BagMD)

	return result
}

func SetBuilding(md *bag.BagMD) {
	md.SetBagMD(
		bag.BagSetTotalMax(configs.GetDb().GameEx.BuildingBagNumMax),
	)
}
