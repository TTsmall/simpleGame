package configs

import "github.com/TTsmall/wawaji_pub_hmhelper/bag"

var gameDb *GameDb

type GameDb struct {
	GameEx             *HelperConf                //单例配置
	HeroCfgs           map[int]*HeroCfg           //英雄配置
	HeroStarCfgs       map[int]*HeroStarCfg       //星级
	GeniusCfgs         map[int]*GeniusCfg         //天赋
	HeroLevelProps     map[int]*HeroLevelProp     //等级属性
	HeroPropertyCfgs   map[int]*HeroPropertyCfg   //英雄属性
	HeroPropGeniusCfgs map[int]*HeroPropGeniusCfg //英雄属性折算
	// HeroMaxCfgs    map[int]*HeroMaxCfg    //品质

	NpcHeroCfgs map[int]*NpcCfg //NPC平民配置

	BuildNumConditions   map[int]*BuildNumConditionCfg
	BuildMinjuWrokers    map[int][]int
	BuildNumLimits       map[int]int
	BuildMinjuByNumType  map[int]*BuildNumConditionCfg
	Builds               map[int]*BuildCfg
	BuildLevelByItemType map[int]map[int]*BuildCfg
	BuildingModels       map[int]*BuildingModelCfg
	Areas                map[int]*AreaCfg
	AreaOpenCostCfgs     map[int]*AreaOpenCostCfg
	StorageCfgs          map[int]*StorageCfg

	GridMaps map[int]map[string]*GridMap //key:cityId,subkey:posX_posY,value:(type,area)
	AreaMaps map[int]map[int]*AreaMapCfg //key:cityId,subkey:area,value:posX_posY

	MapCfgs map[int]*MapCfg
}

func GetDb() *GameDb {
	if gameDb == nil {
		gameDb = new(GameDb)
	}
	return gameDb
}
func (m *GameDb) SetBuildings(buildings map[int]*BuildCfg) {
	m.Builds = buildings
}

func (m *GameDb) SetBuildingModels(buildingModels map[int]*BuildingModelCfg) {
	m.BuildingModels = buildingModels
}

func (m *GameDb) SetAreas(areas map[int]*AreaCfg) {
	m.Areas = areas
}

//bags有关
func (m *GameDb) GetItemCfg(itemId int) *bag.ItemCfg {
	return bag.BagEx.ItemConfEx[itemId]
}

func (m *GameDb) GetAreaMap(cityId, area int) *AreaMapCfg {
	if _, ok := m.AreaMaps[cityId]; !ok {
		return nil
	}
	return m.AreaMaps[cityId][area]
}
