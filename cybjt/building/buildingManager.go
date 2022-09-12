package building

import (
	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
)

type BuildingManager struct {
	// operateBuldingInfos []int //正在操作的建筑物
}

func NewBuildingManager() *BuildingManager {
	//这里要加他的依赖管理器是不是实例化的检查
	if bag.BagEx == nil {
		panic("BagEx is nil.")
	}
	return nil
}

func IsBuildingType(tp int) bool {
	switch tp {
	case DECORATE_BUILDING, AMUSEMENT_BUILDING, PRODUCTION_BUILDING, OCCUP_BUILDING:
		return true
	}
	return false
}

func NewIBuilding(tp int) IBuilding {
	switch tp {
	case DECORATE_BUILDING, AMUSEMENT_BUILDING:
		return NewBaseBuilding()
	case PRODUCTION_BUILDING:
		return NewProductionBuilding()
	case OCCUP_BUILDING:
		return NewOccupBuilding()
	}
	return nil
}
