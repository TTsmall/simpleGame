package building

import (
	"time"

	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/configs"
	"github.com/TTsmall/wawaji_pub_hmhelper/errorx"
)

//生产类
type ProductionBuilding struct {
	BaseBuilding
}

//NewProductionBuilding new
func NewProductionBuilding() *ProductionBuilding {
	a := &ProductionBuilding{}
	a.Building.Buildings = make(map[int]*BuildingInfo)
	return a
}

//LvlUpStartBuild lvlup
func (m *ProductionBuilding) LvlUpStartBuild(uid, id, heroId, useLevelUpTime int, peopleId []int) (*BuildingInfo, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingInfo.Id)
	if buildCfg == nil {
		return nil, errorx.Build_Item_Not_Error
	}
	if (buildingInfo.BuildCompletedTime != 0 || buildingInfo.ProduceCompletedTime != 0) && buildCfg.Itemtype != BUILD_ITEMTYPE_MINJU {
		return nil, errorx.Building_Operate_Err
	}
	buildingInfo.UpLevelCompletedTime = int(time.Now().Unix()) + useLevelUpTime
	buildingInfo.SetBuilder(heroId)
	// buildingInfo.Builder = heroId
	buildingInfo.BuilderPeople = peopleId
	buildingInfo.Id = id
	return buildingInfo, nil
}

func (m *ProductionBuilding) LvlUpEndBuild(uid int) (*BuildingInfo, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	if buildingInfo.UpLevelCompletedTime > int(time.Now().Unix()) || buildingInfo.UpLevelCompletedTime == 0 {
		return nil, errorx.Produce_Time_Not_Error
	}
	buildingInfo.UpLevelCompletedTime = 0
	buildingInfo.SetBuilder(0)
	// buildingInfo.Builder = 0
	buildingInfo.BuilderPeople = make(common.IntSlice, 0)
	return buildingInfo, nil
}

func (m *ProductionBuilding) AwardStart(uid, heroId, useProduceTime int, peopleId []int) (*BuildingInfo, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingInfo.Id)
	if buildCfg == nil {
		return nil, errorx.Build_Item_Not_Error
	}
	if buildCfg.Itemtype == BUILD_ITEMTYPE_MINJU {
		return nil, errorx.Building_Minju_Error
	}
	if buildingInfo.BuildCompletedTime != 0 || buildingInfo.UpLevelCompletedTime != 0 {
		return nil, errorx.Building_Operate_Err
	}
	buildingInfo.ProduceCompletedTime = int(time.Now().Unix()) + useProduceTime
	buildingInfo.SetBuilder(heroId)
	// buildingInfo.Builder = heroId
	buildingInfo.BuilderPeople = peopleId
	return buildingInfo, nil
}

func (m *ProductionBuilding) AwardEnd(uid int) (*BuildingInfo, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	if buildingInfo.ProduceCompletedTime > int(time.Now().Unix()) || buildingInfo.ProduceCompletedTime == 0 {
		return nil, errorx.Produce_Time_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingInfo.Id)
	if buildCfg == nil {
		return nil, errorx.Build_Item_Not_Error
	}
	if buildCfg.Itemtype == BUILD_ITEMTYPE_MINJU {
		return nil, errorx.Building_Minju_Error
	}
	buildingInfo.ProduceCompletedTime = 0
	return buildingInfo, nil
}
