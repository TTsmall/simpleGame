package mmap

import (
	"fmt"

	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	hmcom "github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/bags"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/building"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/configs"
	"github.com/TTsmall/wawaji_pub_hmhelper/log"
)

const (
	X_LEN = 98 //posY= slice(index)/98
	//posX= slice(index)%98
)

type IMapUser interface {
	building.IUser
	bags.IResBagUser
	bags.IPeopleUser
	GetMapBuildingInfos() (result map[int]*building.Building) //(result map[int]map[int]building.IBuilding)
	GetUnLockedArea() (result common.MapIntSlice)
	GetClearBox() (result map[int]map[int]int)
	GetRunis() (result common.RuinsItemsMap)
	GetCityIds() []int
	GetExploreItems() (result common.ExploreItemsMap)
	CheckConditon(common.Conditions) bool
}

type IMap interface {
	Build(user IMapUser, id, heroId int, peopleIds []int, pos common.PosInfo, isWatchVideos bool, makeUpGold int) (building *building.BuildingInfo, heromd *bags.HeroMD, err error)
	ClearBuilding(uid int) error
	MoveBuilding(user IMapUser, uid int, newPos common.PosInfo) (*building.BuildingInfo, error)
	PackageBuilding(uid int) error
	UnPackageBuilding(id int, newPos common.PosInfo) (building *building.BuildingInfo, err error)
	ChangeDirect(uid, direct int) error
	UnlockedArea(areaId int) error
	ClearRuins(item common.RuinsItem) error
}

type MapManager struct {
	infos map[int]*MapObj //key:cityId,value:map
}

func NewMapManager() *MapManager {
	return &MapManager{
		infos: make(map[int]*MapObj),
	}
}

func (m *MapManager) Init(user IMapUser) {
	m.infos = make(map[int]*MapObj)
	openCityIds := user.GetCityIds()
	for _, cityId := range openCityIds {
		m.infos[cityId] = NewMapObj(cityId)
	}
	m.InitBuildingsFromDb(user)
	m.InitUnlockedAreaFromDb(user)
	m.InitClearBox(user)
	m.InitRuinsFromDb(user)
	m.setMapEmpityPos()
	m.InitExploreItemFromDb(user)

}

func (m *MapManager) InitBuildingsFromDb(user IMapUser) {
	mapBuildings := user.GetMapBuildingInfos()
	hmconfDb := configs.GetDb()
	decorateBuildings := make(map[int]*building.BaseBuilding)
	amusementBuildings := make(map[int]*building.BaseBuilding)
	productionBuildings := make(map[int]*building.ProductionBuilding)
	occupBuildings := make(map[int]*building.OccupBuilding)
	log.Info("=====InitBuildingsFromDb  mapBuildings=%v=======", mapBuildings)
	for tp, v := range mapBuildings {
		switch tp {
		case building.DECORATE_BUILDING:
			for _, buildingInfo := range v.Buildings {
				log.Info("tp=%d,buildingInfo=%v", tp, buildingInfo)
				buildCfg := hmconfDb.GetBuild(buildingInfo.Id)
				if buildCfg == nil {
					continue
				}
				if decorateBuildings[buildCfg.City] == nil {
					decorateBuildings[buildCfg.City] = building.NewBaseBuilding()
				}
				decorateBuildings[buildCfg.City].Buildings[buildingInfo.Uid] = buildingInfo
				m.infos[buildCfg.City].SetHappiness(buildCfg.BuildValue)
			}
		case building.AMUSEMENT_BUILDING:
			for _, buildingInfo := range v.Buildings {
				log.Info("tp=%d,buildingInfo=%v", tp, buildingInfo)
				buildCfg := hmconfDb.GetBuild(buildingInfo.Id)
				if buildCfg == nil {
					continue
				}
				if amusementBuildings[buildCfg.City] == nil {
					amusementBuildings[buildCfg.City] = building.NewBaseBuilding()
				}
				amusementBuildings[buildCfg.City].Buildings[buildingInfo.Uid] = buildingInfo
				m.infos[buildCfg.City].SetHappiness(buildCfg.BuildValue)
			}
		case building.PRODUCTION_BUILDING:
			for _, buildingInfo := range v.Buildings {
				log.Info("tp=%d,buildingInfo=%v", tp, buildingInfo)
				buildCfg := hmconfDb.GetBuild(buildingInfo.Id)
				if buildCfg == nil {
					continue
				}
				if productionBuildings[buildCfg.City] == nil {
					productionBuildings[buildCfg.City] = building.NewProductionBuilding()
				}
				productionBuildings[buildCfg.City].Buildings[buildingInfo.Uid] = buildingInfo
				m.infos[buildCfg.City].SetHappiness(buildCfg.BuildValue)
			}
		case building.OCCUP_BUILDING:
			for _, buildingInfo := range v.Buildings {
				log.Info("tp=%d,buildingInfo=%v", tp, buildingInfo)
				buildCfg := hmconfDb.GetBuild(buildingInfo.Id)
				if buildCfg == nil {
					continue
				}
				if occupBuildings[buildCfg.City] == nil {
					occupBuildings[buildCfg.City] = building.NewOccupBuilding()
				}
				occupBuildings[buildCfg.City].Buildings[buildingInfo.Uid] = buildingInfo
				m.infos[buildCfg.City].SetHappiness(buildCfg.BuildValue)
				if buildCfg.Itemtype == building.BUILD_ITEMTYPE_CAP {
					caps := configs.GetDb().GetStorageByLvl(buildCfg.Itemtype, buildCfg.Level)
					m.infos[buildCfg.City].AddResourceRepCap(caps)
				}
			}
		}

	}
	log.Info("+++++++++++------------ len decorateBuildings=%d", len(decorateBuildings))
	for cityId, v := range decorateBuildings {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		log.Info("+++++++++++------------ decorateBuildings cityId=%d, v=%v", cityId, v)
		m.infos[cityId].AddBuldingInfos(building.DECORATE_BUILDING, v, v.Buildings) //.buldingInfos[building.DECORATE_BUILDING] = v
	}
	log.Info("+++++++++++------------len amusementBuildings=%d", len(amusementBuildings))
	for cityId, v := range amusementBuildings {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		log.Info("+++++++++++------------ amusementBuildings cityId=%d, v=%v", cityId, v)
		m.infos[cityId].AddBuldingInfos(building.AMUSEMENT_BUILDING, v, v.Buildings)
	}
	log.Info("+++++++++++------------len productionBuildings=%d", len(productionBuildings))
	for cityId, v := range productionBuildings {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		log.Info("+++++++++++------------productionBuildings cityId=%d, v=%v", cityId, v)
		m.infos[cityId].AddBuldingInfos(building.PRODUCTION_BUILDING, v, v.Buildings)
		a := m.infos[1].buldingInfos[building.PRODUCTION_BUILDING].(*building.ProductionBuilding)
		for _, v := range a.Buildings {
			log.Info("+++++++++++------------ProductionBuilding building=%v", v)
		}
	}

	log.Info("--1=============================================")
	fmt.Println("-1=============================================")
	if m.infos[1] != nil && m.infos[1].buldingInfos[building.PRODUCTION_BUILDING] != nil {
		for _, v := range m.infos[1].buldingInfos[building.PRODUCTION_BUILDING].(*building.ProductionBuilding).Buildings {
			log.Info("--1building=%v", v)
			fmt.Println("-1building=", v)

		}
	}
	log.Info("--2=============================================")
	fmt.Println("-2=============================================")

	log.Info("+++++++++++------------len occupBuildings=%d", len(occupBuildings))
	for cityId, v := range occupBuildings {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		log.Info("+++++++++++------------occupBuildings cityId=%d, v=%v", cityId, v)
		m.infos[cityId].AddBuldingInfos(building.OCCUP_BUILDING, v, v.Buildings)
		// a := m.infos[1].buldingInfos[building.OCCUP_BUILDING].(*building.OccupBuilding)
		// for _, v := range a.Buildings {
		// 	log.Info("+++++++++++------------occupBuildings building=%v", v)
		// }
	}

	log.Info("--11=============================================")
	fmt.Println("-11=============================================")
	if m.infos[1] != nil && m.infos[1].buldingInfos[building.PRODUCTION_BUILDING] != nil {
		for _, v := range m.infos[1].buldingInfos[building.PRODUCTION_BUILDING].(*building.ProductionBuilding).Buildings {
			log.Info("--11building=%v", v)
			fmt.Println("-11building=", v)

		}
	}
	log.Info("--22=============================================")
	fmt.Println("-22=============================================")
}

func (m *MapManager) InitUnlockedAreaFromDb(user IMapUser) {
	for cityId, v := range user.GetUnLockedArea() {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		m.infos[cityId].SetUnlockedArea(v)
	}
}

func (m *MapManager) InitClearBox(user IMapUser) {
	for cityId, v := range user.GetClearBox() {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		m.infos[cityId].SetClearBox(v)
	}
}

func (m *MapManager) InitRuinsFromDb(user IMapUser) {
	for cityId, v := range user.GetRunis() {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		m.infos[cityId].SetRuins(v)
	}
}

func (m *MapManager) InitExploreItemFromDb(user IMapUser) {
	for cityId, v := range user.GetExploreItems() {
		if m.infos[cityId] == nil {
			m.infos[cityId] = NewMapObj(cityId)
		}
		m.infos[cityId].SetExploreItem(v)
	}
}

func (m *MapManager) GetBuildingByCityId(cityId int) map[int]building.IBuilding {
	if m.infos[cityId] == nil {
		fmt.Println("MapManager GetBuildingByCityId cityid nil", cityId)
	}
	return m.infos[cityId].GetBuldingInfos()
}

func (m *MapManager) GetUnlockedAreaByCity(cityId int) common.IntSlice {
	return m.infos[cityId].GetUnlockedArea()
}
func (m *MapManager) GetClearBoxByCity(cityId int) map[int]int {
	return m.infos[cityId].GetClearBox()
}

func (m *MapManager) GetExploreItemByCity(cityId int) common.ExploreItems {
	return m.infos[cityId].GetExploreItems()
}

func (m *MapManager) GetResourceRepCapByCity(cityId int) int {
	return m.infos[cityId].GetResourceRepCap()
}

func (m *MapManager) ReturnAll() (map[int]building.IBuilding, common.MapIntSlice, map[int]map[int]int, common.RuinsItemsMap, common.ExploreItemsMap) {
	relustBuilding := m.GetBuildings()
	unlockedArea := m.GetUnlockedArea()
	clearBox := m.GetClearBox()
	ruins := m.GetRuins()
	return relustBuilding, unlockedArea, clearBox, ruins, m.GetExploreItem()
}

func (m *MapManager) GetBuildings() map[int]building.IBuilding {
	relustBuilding := make(map[int]building.IBuilding)
	for _, v := range m.infos {
		if len(relustBuilding) == 0 {
			relustBuilding = v.GetBuldingInfos()
			continue
		}
		for tp, IBuildingInfo := range v.GetBuldingInfos() {
			switch tp {
			case building.DECORATE_BUILDING:
				buildingInfos := IBuildingInfo.(*building.BaseBuilding)
				for _, buildingInfo := range buildingInfos.Buildings {
					relustBuilding[tp].SetBuildingInfo(buildingInfo)
				}
			case building.AMUSEMENT_BUILDING:
				buildingInfos := IBuildingInfo.(*building.BaseBuilding)
				for _, buildingInfo := range buildingInfos.Buildings {
					relustBuilding[tp].SetBuildingInfo(buildingInfo)
				}
			case building.PRODUCTION_BUILDING:
				buildingInfos := IBuildingInfo.(*building.ProductionBuilding)
				for _, buildingInfo := range buildingInfos.Buildings {
					relustBuilding[tp].SetBuildingInfo(buildingInfo)
				}
			case building.OCCUP_BUILDING:
				buildingInfos := IBuildingInfo.(*building.OccupBuilding)
				for _, buildingInfo := range buildingInfos.Buildings {
					relustBuilding[tp].SetBuildingInfo(buildingInfo)
				}
			}
		}
	}
	return relustBuilding
}

func (m *MapManager) GetUnlockedArea() common.MapIntSlice {
	result := make(common.MapIntSlice)
	for cityId, v := range m.infos {
		result[cityId] = v.GetUnlockedArea()
	}
	return result
}

func (m *MapManager) GetClearBox() map[int]map[int]int {
	result := make(map[int]map[int]int)
	for cityId, v := range m.infos {
		result[cityId] = v.GetClearBox()
	}
	return result
}

func (m *MapManager) GetRuins() common.RuinsItemsMap {
	result := make(common.RuinsItemsMap)
	for cityId, v := range m.infos {
		result[cityId] = v.GetRuins()
	}
	return result
}

func (m *MapManager) GetExploreItem() common.ExploreItemsMap {
	result := make(common.ExploreItemsMap)
	for cityId, v := range m.infos {
		result[cityId] = v.GetExploreItems()
	}
	return result
}

func (m *MapManager) setMapEmpityPos() {
	for _, v := range m.infos {
		v.setMapEmpityPos()
	}
}

func (m *MapManager) GetHappiness(cityId int) int {
	return m.infos[cityId].GetHappiness()
}

func (m *MapManager) CheckBuild(user IMapUser) {

}

func (m *MapManager) InitBuild(user IMapUser, id int, pos common.PosInfo) (err error) {
	return m.infos[user.GetCurrCityID()].BuildByInit(user, id, pos)
}

func (m *MapManager) Build(user IMapUser, id, heroId int, peopleIds []int, pos common.PosInfo, isWatchVideos bool, makeUpGold int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	return m.infos[user.GetCurrCityID()].Build(user, id, heroId, peopleIds, pos, isWatchVideos, makeUpGold)
}

func (m *MapManager) FinishBuild(user IMapUser, uid int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, peoples []*bags.PeopleMD, err error) {
	return m.infos[user.GetCurrCityID()].FinishBuild(user, uid)
}

func (m *MapManager) StartProduce(user IMapUser, uid, heroId int, peopleIds []int, isWatchVideos bool, makeUpGold int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	return m.infos[user.GetCurrCityID()].StartProduce(user, uid, heroId, peopleIds, isWatchVideos, makeUpGold)
}

func (m *MapManager) FinishProduce(user IMapUser, uid int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	return m.infos[user.GetCurrCityID()].FinishProduce(user, uid)
}

func (m *MapManager) StartUpLevel(user IMapUser, uid, heroId int, peopleIds []int, isWatchVideos bool, makeUpGold int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	return m.infos[user.GetCurrCityID()].StartUpLevel(user, uid, heroId, peopleIds, isWatchVideos, makeUpGold)
}

func (m *MapManager) FinishUpLevel(user IMapUser, uid int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, peoples []*bags.PeopleMD, err error) {
	return m.infos[user.GetCurrCityID()].FinishLevelUp(user, uid)
}

func (m *MapManager) BuildingSpeedUp(user IMapUser, uid, costType int) (*building.BuildingInfo, []bag.IBagItem, error) {
	return m.infos[user.GetCurrCityID()].BuildingSpeedUp(user, uid, costType)
}

func (m *MapManager) LevelUpSpeedUp(user IMapUser, uid, costType int) (*building.BuildingInfo, []bag.IBagItem, error) {
	return m.infos[user.GetCurrCityID()].LevelUpSpeedUp(user, uid, costType)
}

func (m *MapManager) ProduceSpeedUp(user IMapUser, uid, costType int) (*building.BuildingInfo, []bag.IBagItem, error) {
	return m.infos[user.GetCurrCityID()].ProduceSpeedUp(user, uid, costType)
}

func (m *MapManager) MoveBuilding(user IMapUser, uid int, newPos common.PosInfo) (*building.BuildingInfo, error) {
	return m.infos[user.GetCurrCityID()].MoveBuilding(user, uid, newPos)
}

func (m *MapManager) AwardMinju(user IMapUser, uid int) (*building.BuildingInfo, []bag.IBagItem, error) {
	return m.infos[user.GetCurrCityID()].AwardMinju(user, uid)
}

func (m *MapManager) OccupMinju(user IMapUser, uid, heorId int) (*building.BuildingInfo, []*bags.HeroMD, error) {
	return m.infos[user.GetCurrCityID()].OccupMinju(user, uid, heorId)
}

func (m *MapManager) ReleasePeople(user IMapUser, uid, releaseType int) (*building.BuildingInfo, *bags.HeroMD, []*bags.PeopleMD, error) {
	return m.infos[user.GetCurrCityID()].ReleasePeople(user, uid, releaseType)
}

func (m *MapManager) DelBuilding(user IMapUser, uid int) (int, error) {
	return m.infos[user.GetCurrCityID()].DelBuilding(user, uid)
}

func (m *MapManager) ExploreArea(user IMapUser, cityId, id, heroId, titleX, titleY, buildId int, peopleIds common.IntSlice) ([]bag.IBagItem, common.ExploreItem, *bags.HeroMD, []*bags.PeopleMD, error) {
	return m.infos[cityId].ExploreArea(user, id, heroId, titleX, titleY, buildId, peopleIds)
}

func (m *MapManager) UnlockArea(user IMapUser, cityId, id int, isUseGold bool) ([]bag.IBagItem, *bags.HeroMD, []*bags.PeopleMD, error) {
	return m.infos[cityId].UnlockArea(user, cityId, id, isUseGold)
}

func (m *MapManager) ClearRuins(user IMapUser, cityId, id int, ruins *hmcom.RuinsItem) (buildingInfo *building.BuildingInfo, err error) {
	// buildcfg := configs.GetDb().GetBuild(id)
	// if buildcfg == nil {
	// 	return nil, errorx.Build_Item_Not_Error
	// }
	// return m.infos[buildcfg.City].Build(user, id, heroId, peopleIds, pos, isWatchVideos)

	// buildingInfo, err = m.buldingInfos[buildcfg.TitleType].StartBuild(user, id, heroId, peopleIds, pos)

	return nil, nil
}

func (m *MapManager) InitUnlockArea(user IMapUser, cityId, id int) error {
	return m.infos[cityId].InitUnlockArea(user, cityId, id)
}
