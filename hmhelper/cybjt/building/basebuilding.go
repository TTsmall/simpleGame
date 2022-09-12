package building

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"wawaji_pub/hmhelper/bag"
	"wawaji_pub/hmhelper/common"
	"wawaji_pub/hmhelper/cybjt/bags"
	"wawaji_pub/hmhelper/cybjt/configs"
	"wawaji_pub/hmhelper/errorx"
)

func (m *Building) StartBuild(user IUser, id, heroId, useBuildingTime int, peopleIds []int, pos common.PosInfo) (*BuildingInfo, error) {
	buildcfg := configs.GetDb().GetBuild(id)
	if buildcfg == nil {
		return nil, errorx.Build_Item_Not_Error
	}
	upEndTime := int(time.Now().Unix()) + useBuildingTime
	newBuild := NewBuildingInfo(id, user.GetBuildIncrId(), pos, upEndTime, heroId, peopleIds)
	m.Buildings[newBuild.Uid] = newBuild
	return newBuild, nil
}

func (m *Building) EndBuild(user IUser, uid int) (*BuildingInfo, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	if buildingInfo.BuildCompletedTime > int(time.Now().Unix()) || buildingInfo.BuildCompletedTime == 0 {
		return nil, errorx.Build_Item_Time_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingInfo.Id)
	if buildCfg == nil {
		return nil, errorx.Build_Item_Not_Error
	}
	buildingInfo.SetBuilder(0)
	// buildingInfo.Builder = 0
	buildingInfo.BuilderPeople = make(common.IntSlice, 0)
	buildingInfo.BuildCompletedTime = 0
	if buildCfg.Itemtype == BUILD_ITEMTYPE_MINJU {
		buildingInfo.ProduceCompletedTime = int(time.Now().Unix()) + buildCfg.WorkTime
	}
	return buildingInfo, nil
}

func (m *Building) InitBuild(user IUser, id int, pos common.PosInfo) (*BuildingInfo, error) {
	buildcfg := configs.GetDb().GetBuild(id)
	if buildcfg == nil {
		return nil, errorx.Build_Item_Not_Error
	}
	newBuild := NewBuildingInfo(id, user.GetBuildIncrId(), pos, 0, 0, make(common.IntSlice, 0))
	if buildcfg.Itemtype == BUILD_ITEMTYPE_MINJU {
		newBuild.ProduceCompletedTime = int(time.Now().Unix()) + buildcfg.WorkTime
	}
	m.Buildings[newBuild.Uid] = newBuild
	return newBuild, nil
}

func (m *Building) DelteBuild(user IUser, uid int) (int, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return 0, errorx.Build_Item_User_Not_Error
	}
	delete(m.Buildings, uid)
	return buildingInfo.Uid, nil
}

func (m *Building) GetCount(id int) int {
	count := 0
	for _, v := range m.Buildings {
		if v.Id == id {
			count++
		}
	}
	return count
}

func (m *Building) GetCountByItemType(itemType int) int {
	fmt.Println("Building GetCountByItemType itemType=,m.Buildings=", itemType, m.Buildings)
	count := 0
	conf := configs.GetDb()
	for _, v := range m.Buildings {
		buildCfg := conf.GetBuild(v.Id)
		if buildCfg == nil {
			continue
		}
		if buildCfg.Itemtype == itemType {
			count++
		}
	}
	return count
}

func (m *Building) CheckLimit(user IUser, id int) bool {
	buildCfg := configs.GetDb().GetBuild(id)
	if buildCfg == nil {
		return true
	}
	limits := configs.GetDb().BuildNumLimit(buildCfg.Itemtype)
	if limits <= 0 {
		//no limit
		return false
	}
	pk := user.GetBag(bags.BagKey_Building).(*bag.BagMD)
	currentPKNum := pk.GetItemNum(id)
	currentNum := currentPKNum + m.GetCountByItemType(buildCfg.Itemtype)
	if currentNum >= limits {
		return true
	}
	return false
}

func (m *Building) GetPos(uid int) common.PosInfo {
	return m.Buildings[uid].Pos
}

func (m *Building) GetAllPos() []*BuildingPos {
	allPos := make([]*BuildingPos, 0)
	for _, v := range m.Buildings {
		buildCfg := configs.GetDb().GetBuild(v.Id)
		if buildCfg == nil {
			continue
		}
		allPos = append(allPos, &BuildingPos{Id: v.Id, Pos: v.Pos})
	}
	return allPos
}

func (m *Building) GetBuildings() map[int]*BuildingInfo {
	return m.Buildings
}

func (m *Building) SetBuildingInfo(building *BuildingInfo) {
	m.Buildings[building.Uid] = building
}

func (m *Building) GetBuildingBuildEnd(building *BuildingInfo) {
	m.Buildings[building.Uid] = building
}

func (m *Building) BuildingSpeedUp(uid, speedTime int) (*BuildingInfo, error) {
	building := m.Buildings[uid]
	if building == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	building.BuildCompletedTime -= speedTime
	return building, nil
}

func (m *Building) LevelUpSpeedUp(uid, speedTime int) (*BuildingInfo, error) {
	building := m.Buildings[uid]
	if building == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	building.UpLevelCompletedTime -= speedTime
	return building, nil
}

func (m *Building) ProduceSpeedUp(uid, speedTime int) (*BuildingInfo, error) {
	building := m.Buildings[uid]
	if building == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	building.ProduceCompletedTime -= speedTime
	return building, nil
}

func (m *Building) MoveBuilding(uid int, newPos common.PosInfo) (*BuildingInfo, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	buildingInfo.Pos = newPos
	return buildingInfo, nil
}

func (m *Building) ReleasePeople(uid, releaseType int) (*BuildingInfo, int, []int, error) {
	buildingInfo := m.Buildings[uid]
	if buildingInfo == nil {
		return nil, 0, nil, errorx.Build_Item_User_Not_Error
	}
	buildHeroId := buildingInfo.Builder
	buildPeople := buildingInfo.BuilderPeople
	fmt.Println("======ReleasePeople======buildingInfo=", buildingInfo)
	fmt.Println("buildingInfo.ProduceCompletedTime=", buildingInfo.ProduceCompletedTime)
	fmt.Println("time.Now().Unix()=", int(time.Now().Unix()))
	switch releaseType {
	case RELEASETYPE_BUILDING:
		if buildingInfo.BuildCompletedTime > int(time.Now().Unix()) || buildingInfo.BuildCompletedTime == 0 {
			return nil, 0, nil, errorx.Produce_Time_Not_Error
		}
	case RELEASETYPE_LEVELUP:
		if buildingInfo.UpLevelCompletedTime > int(time.Now().Unix()) || buildingInfo.UpLevelCompletedTime == 0 {
			return nil, 0, nil, errorx.Produce_Time_Not_Error
		}
	case RELEASETYPE_PRODUCE:
		if buildingInfo.ProduceCompletedTime > int(time.Now().Unix()) || buildingInfo.ProduceCompletedTime == 0 {
			return nil, 0, nil, errorx.Produce_Time_Not_Error
		}
	case RELEASETYPE_EXPLORE:
	default:
		return nil, 0, nil, errorx.Release_Building_No_Err
	}

	buildCfg := configs.GetDb().GetBuild(buildingInfo.Id)
	if buildCfg == nil {
		return nil, 0, nil, errorx.Build_Item_Not_Error
	}

	buildingInfo.SetBuilder(0)
	buildingInfo.BuilderPeople = make(common.IntSlice, 0)
	return buildingInfo, buildHeroId, buildPeople, nil
}

func (m *HeroSites) Scan(value interface{}) error {
	fmt.Println("=============HeroSites Scan data=", value)
	err := json.Unmarshal(value.([]byte), m)
	if err != nil {
		fmt.Println("HeroSites Scan err=", err)
		return err
	}
	return nil
	//return this.UnmarshalJSON(value.([]byte))
}

func (m HeroSites) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Building) Scan(value interface{}) error {
	fmt.Println("=============Building Scan data=", value)
	err := json.Unmarshal(value.([]byte), m)
	if err != nil {
		fmt.Println("Building Scan err=", err)
		return err
	}
	return nil
	//return this.UnmarshalJSON(value.([]byte))
}

func (m Building) Value() (driver.Value, error) {
	return json.Marshal(m)
}

//基础建筑(装饰类)
type BaseBuilding struct {
	Building
}

func (m *BaseBuilding) Init(building map[int]*BuildingInfo) {
	m.Buildings = building
}

func NewBaseBuilding() *BaseBuilding {
	a := &BaseBuilding{}
	a.Building.Buildings = make(map[int]*BuildingInfo)
	return a
}
