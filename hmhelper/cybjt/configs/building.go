package configs

import (
	"wawaji_pub/hmhelper/common"
)

//所有building

type BuildNumConditionCfg struct {
	Id                 int             `col:"id" client:"id"`                   //日志ID
	BuildType          int             `col:"buildType" client:"buildType"`     //
	Name               string          `col:"name" client:"name"`               //
	MaxBulidNum        int             `col:"maxBulidNum" client:"maxBulidNum"` //
	BuildNum           int             `col:"buildNum" client:"buildNum"`
	NumUnlockCondition string          `col:"numUnlockCondition" client:"numUnlockCondition"`
	LockInfo           string          `col:"lockInfo" client:"lockInfo"`
	Npc                common.IntSlice `col:"npc" client:"npc"`
}

type BuildCfg struct {
	Id            int    `col:"id" client:"id"`     //ID
	City          int    `col:"city" client:"city"` //名字                     //头像
	Name          string `col:"name" client:"name"` //角色立绘
	Level         int    `col:"level" client:"level"`
	MaxLevel      int    `col:"maxLevel" client:"maxLevel"`
	Icon          string `col:"icon" client:"icon"` //地基模式icon
	NormalModel   int    `col:"normalModel" client:"normalModel"`
	BaceModel     string `col:"baceModel" client:"baceModel"` //地基模式icon
	TitleType     int    `col:"titleType" client:"titleType"`
	WorkType      int    `col:"workType" client:"workType"`
	Itemtype      int    `col:"buildType" client:"buildType"`
	BuildingModel string `col:"buildingModel" client:"buildingModel"` //未建成地表模型
	WorkModel     string `col:"workModel" client:"workModel"`         //生产中的模型
	WorkOutModel  string `col:"workOutModel" client:"workOutModel"`   //完成生产的模型
	//UpLevelModel     string            `col:"upLevelModel" client:"upLevelModel"`   //升级中的模型
	IfTouch          int               `col:"ifTouch" client:"ifTouch"`
	IfMove           int               `col:"ifMove" client:"ifMove"`
	Info             string            `col:"info" client:"info"`
	SpecialPicture   string            `col:"specialPicture" client:"specialPicture"`
	LockInfo         string            `col:"lockInfo" client:"lockInfo"`
	UnbulidInfo      string            `col:"unbulidInfo" client:"unbulidInfo"`
	Workicon         int               `col:"workIcon" client:"workIcon"`
	BulidInfo        string            `col:"bulidInfo" client:"bulidInfo"`
	UnLockCondition  common.Conditions `col:"unLockCondition" client:"unLockCondition"`
	BuildCondition   common.Conditions `col:"buildCondition" client:"buildCondition"`
	LevelupCondition common.Conditions `col:"levelupCondition" client:"levelupCondition"`
	LevelupInfo      string            `col:"levelupInfo" client:"levelupInfo"`
	NeedHero         int               `col:"needHero" client:"needHero"`
	PeopleNum        int               `col:"peopleNum" client:"peopleNum"`
	BuildTime        int               `col:"buildTime" client:"buildTime"`
	BuildCost        common.ItemInfos  `col:"buildCost" client:"buildCost"`
	BuildValue       common.ItemInfos  `col:"buildValue" client:"buildValue"`
	BulidExp         int               `col:"bulidExp" client:"bulidExp"`
	WorkHero         int               `col:"workHero" client:"workHero"`
	WorkPeopleNum    int               `col:"workPeopleNum" client:"workPeopleNum"`
	WorkTime         int               `col:"workTime" client:"workTime"`
	WorkCost         common.ItemInfos  `col:"workCost" client:"workCost"`
	WorkOut          common.ItemInfos  `col:"workOut" client:"workOut"`
	WorkHeroExp      int               `col:"workHeroExp" client:"workHeroExp"`
	// SpecialWorkCd    int               `col:"specialWorkCd" client:"specialWorkCd"`
	// SpecialWorkOut string `col:"specialWorkOut" client:"specialWorkOut"`
	MovingOperate string `col:"movingOperate" client:"movingOperate"`
	// BuffArea         common.IntSlice   `col:"buffArea" client:"buffArea"`
	// BuffBulidSubtype int              `col:"buffBulidSubtype" client:"buffBulidSubtype"`
	//Buffvalue common.ItemInfos `col:"buffvalue" client:"buffvalue"`
	// MReward   common.ItemInfos `col:"mReward" client:"mReward"`
	// MaxReward common.ItemInfos `col:"maxReward" client:"maxReward"`
	HouseWorkSpeed common.ItemInfo `col:"houseWorkSpeed" client:"houseWorkSpeed"`
	HousePeopleNum int             `col:"housePeopleNum" client:"housePeopleNum"`
	GoldCost       common.ItemInfo `col:"goldCost" client:"goldCost"`
}

type BuildingModelCfg struct {
	ModelId int             `col:"modelId" client:"modelId"` //ID
	Name    string          `col:"name" client:"name"`       //名字
	BuildID int             `col:"buildId" client:"buildId"`
	Size    common.IntSlice `col:"size" client:"size"`
	Sacle   string          `col:"sacle" client:"sacle"`
	Res     string          `col:"res" client:"res"`
	Rect    common.IntSlice `col:"rect" client:"rect"`
}

//仓库配置
type StorageCfg struct {
	Id           int             `col:"id" client:"id"`
	BuildType    int             `col:"buildType" client:"buildType"`
	Desc         string          `col:"desc" client:"desc"`
	Info         string          `col:"info" client:"info"`
	StorageLimit common.IntSlice `col:"storageLimit" client:"storageLimit"`
}

type AreaCfg struct {
	Id         int              `col:"id" client:"id"`     //
	City       int              `col:"city" client:"city"` //ID
	Area       string           `col:"area" client:"area"` //名字
	Info       string           `col:"info" client:"info"`
	AreaBuild  string           `col:"areaBuild" client:"areaBuild"`
	AreaStage  common.PropInfos `col:"areaStage" client:"areaStage"`
	NearbyArea string           `col:"nearbyArea" client:"nearbyArea"`
}

type AreaOpenCostCfg struct {
	Id              int               `col:"id" client:"id"`     //
	City            int               `col:"city" client:"city"` //ID
	OpenNum         int               `col:"areaOpenNum" client:"areaOpenNum"`
	UnLockCondition common.Conditions `col:"unLockCondition" client:"unLockCondition"`
	NeedHero        int               `col:"needHero" client:"needHero"`
	PeopleNum       int               `col:"peopleNum" client:"peopleNum"`
	UseTime         int               `col:"time" client:"time"`
	Cost            common.ItemInfos  `col:"cost" client:"cost"`
	Gold            common.ItemInfos  `col:"gold" client:"gold"`
	HeroExp         int               `col:"reward" client:"reward"` //奖励英雄经验
}

type MapCfg struct {
	Id        int    `col:"id" client:"id"`
	Name      string `col:"name" client:"name"`
	Res       string `col:"res" client:"res"`
	Condition string `col:"condition" client:"condition"`
}

//----------------有关方法

func (m *GameDb) BuildNumLimit(itemType int) int {
	return m.BuildNumLimits[itemType]
}

func (m *GameDb) GetBuild(id int) *BuildCfg {
	//fmt.Println("====MapObj StartUpLevel m.Builds=", m.Builds)
	return m.Builds[id]
}

func (m *GameDb) GetArea(id int) *AreaCfg {
	return m.Areas[id]
}

func (m *GameDb) GetBuildingModel(id int) *BuildingModelCfg {
	return m.BuildingModels[id]
}

func (m *GameDb) GetBuildByLevel(itemType, level int) *BuildCfg {
	if len(m.BuildLevelByItemType[itemType]) > 0 {
		return m.BuildLevelByItemType[itemType][level]
	}
	return nil
}

func (m *GameDb) GetBuildingMinju(num int) *BuildNumConditionCfg {
	return m.BuildMinjuByNumType[num]
}

func (m *GameDb) GetBuildAddHappinessV(id int) int {
	if m.Builds[id] != nil {
		for _, v := range m.Builds[id].BuildValue {
			if v.ItemId == ItemType_Happines {
				return v.Count
			}
		}
	}
	return 0
}

func (m *GameDb) GetHappinesV(items common.ItemInfos) int {
	for _, v := range items {
		if v.ItemId == ItemType_Happines {
			return v.Count
		}
	}
	return 0
}

type HappinessCfg struct {
	Id             int `col:"id" client:"id"`                         //
	HappinessLevel int `col:"happinessLevel" client:"happinessLevel"` //ID
	Happiness      int `col:"happiness" client:"happiness"`           //
}

func (m *GameDb) GetAreaCfg(id int) *AreaCfg {
	return m.Areas[id]
}

func (m *GameDb) GetAreaOpenCostCfg(id int) *AreaOpenCostCfg {
	return m.AreaOpenCostCfgs[id]
}

func (m *GameDb) GetNpcByOne(zeroNpcId int) common.IntSlice {
	return m.BuildMinjuWrokers[zeroNpcId]
}

func (m *GameDb) GetStorageCfg(buildType int) *StorageCfg {
	return m.StorageCfgs[buildType]
}

func (m *GameDb) GetStorageByLvl(buildType, lvl int) int {
	storageCfg := m.StorageCfgs[buildType]
	if storageCfg == nil {
		return 0
	}
	maxLvl := len(storageCfg.StorageLimit)
	if lvl >= maxLvl {
		return storageCfg.StorageLimit[maxLvl-1]
	}
	return storageCfg.StorageLimit[lvl-1]
}
