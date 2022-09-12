package building

import (
	"time"

	"wawaji_pub/hmhelper/common"
	"wawaji_pub/hmhelper/cybjt/configs"
	"wawaji_pub/hmhelper/errorx"
)

type OccupBuilding struct {
	ProductionBuilding
}

func NewOccupBuilding() *OccupBuilding {
	a := &OccupBuilding{}
	a.Building.Buildings = make(map[int]*BuildingInfo)
	return a
}

func (m *OccupBuilding) AwardMinju(uid int) (*BuildingInfo, *common.ItemInfos, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildcfg := configs.GetDb().GetBuild(buildingInfo.Id)
	if buildcfg == nil {
		return nil, nil, errorx.Build_Item_Not_Error
	}
	nowTime := int(time.Now().Unix())
	if buildingInfo.ProduceCompletedTime > nowTime {
		return nil, nil, errorx.Produce_Time_Not_Error
	}
	buildingInfo.ProduceCompletedTime = nowTime + buildcfg.WorkTime
	return buildingInfo, &buildcfg.WorkOut, nil
}

func (m *OccupBuilding) OccupBuilding(uid, heroId int) (*BuildingInfo, int, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, 0, errorx.Build_Item_User_Not_Error
	}

	//update old building
	for _, v := range m.Buildings {
		if v.Occupiers[0] == heroId {
			v.Occupiers[0] = 0
			break
		}
	}

	oldOccuperId := buildingInfo.Occupiers[0]
	buildcfg := configs.GetDb().GetBuild(buildingInfo.Id)
	if buildcfg == nil {
		return nil, 0, errorx.Build_Item_Not_Error
	}
	if buildcfg.Itemtype != BUILD_ITEMTYPE_MINJU {
		return buildingInfo, 0, errorx.Occup_Minju_Not_Error
	}
	buildingInfo.Occupiers[0] = heroId

	return buildingInfo, oldOccuperId, nil
}
