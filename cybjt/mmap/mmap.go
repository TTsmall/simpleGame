package mmap

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/TTsmall/wawaji_pub_hmhelper/log"

	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
	"github.com/TTsmall/wawaji_pub_hmhelper/common"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/bags"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/building"
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/configs"
	"github.com/TTsmall/wawaji_pub_hmhelper/errorx"
)

type MapObj struct {
	cityId       int
	buldingInfos map[int]building.IBuilding //地图上所有建筑物
	mapEmpityPos map[string]bool            //解锁区域-已放置区域
	unlockedArea common.IntSlice            //解锁区域
	clearBox     map[int]int                //已清除格子信息
	exploreItems common.ExploreItems        //探索区域
	ruins        common.RuinsItems          //废墟信息
	happyV       int                        //幸福值

	buldingInfosByUid map[int]*building.BuildingInfo //temp
	resourceRepCap    int                            //资源仓库容量
}

func NewMapObj(cityId int) *MapObj {
	return &MapObj{
		buldingInfos:      make(map[int]building.IBuilding),
		mapEmpityPos:      make(map[string]bool),
		unlockedArea:      make(common.IntSlice, 0),
		clearBox:          make(map[int]int),
		ruins:             make(common.RuinsItems, 0),
		buldingInfosByUid: make(map[int]*building.BuildingInfo),
		exploreItems:      make(common.ExploreItems, 0),
		cityId:            cityId,
	}
}

func (m *MapObj) GetBuldingInfos() map[int]building.IBuilding {
	return m.buldingInfos
}

func (m *MapObj) GetUnlockedArea() common.IntSlice {
	return m.unlockedArea
}

func (m *MapObj) GetClearBox() map[int]int {
	return m.clearBox
}

func (m *MapObj) GetExploreItems() common.ExploreItems {
	return m.exploreItems
}

func (m *MapObj) GetRuins() common.RuinsItems {
	return m.ruins
}

func (m *MapObj) GetResourceRepCap() int {
	return m.resourceRepCap
}

func (m *MapObj) AddBuldingInfos(tp int, buildingInfos building.IBuilding, buildingInfosB map[int]*building.BuildingInfo) {
	m.buldingInfos[tp] = buildingInfos
	// switch tp {
	// case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
	// 	a := m.buldingInfos[tp].(*building.BaseBuilding)
	// 	for _, v := range a.Buildings {
	// 		log.Info("++++++++ before  AddBuldingInfos  BaseBuilding building=%v", v)
	// 	}
	// case building.PRODUCTION_BUILDING:
	// 	a := m.buldingInfos[tp].(*building.ProductionBuilding)
	// 	for _, v := range a.Buildings {
	// 		log.Info("++++++++ before  AddBuldingInfos ProductionBuilding building=%v", v)
	// 	}
	// case building.OCCUP_BUILDING:
	// 	a := m.buldingInfos[tp].(*building.OccupBuilding)
	// 	for _, v := range a.Buildings {
	// 		log.Info("++++++++ before  AddBuldingInfos OccupBuilding building=%v", v)
	// 	}
	// }

	for uid, v := range buildingInfosB {
		m.buldingInfosByUid[uid] = v
	}

	// switch tp {
	// case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
	// 	a := m.buldingInfos[tp].(*building.BaseBuilding)
	// 	for _, v := range a.Buildings {
	// 		log.Info("++++++++ after  AddBuldingInfos  BaseBuilding building=%v", v)
	// 	}
	// case building.PRODUCTION_BUILDING:
	// 	a := m.buldingInfos[tp].(*building.ProductionBuilding)
	// 	for _, v := range a.Buildings {
	// 		log.Info("++++++++ after  AddBuldingInfos ProductionBuilding building=%v", v)
	// 	}
	// case building.OCCUP_BUILDING:
	// 	a := m.buldingInfos[tp].(*building.OccupBuilding)
	// 	for _, v := range a.Buildings {
	// 		log.Info("++++++++ after  AddBuldingInfos OccupBuilding building=%v", v)
	// 	}
	// }
}

func (m *MapObj) SetUnlockedArea(unlockedArea common.IntSlice) {
	m.unlockedArea = unlockedArea
}

func (m *MapObj) SetClearBox(clearbox map[int]int) {
	m.clearBox = clearbox
}

func (m *MapObj) SetExploreItem(exploreItems common.ExploreItems) {
	m.exploreItems = exploreItems
}

func (m *MapObj) SetRuins(ruins common.RuinsItems) {
	m.ruins = ruins
}

func (m *MapObj) GetPosKey(posX, posY int) string {
	return strconv.Itoa(posX) + "_" + strconv.Itoa(posY)
}

func (m *MapObj) CheckUsedMapPos(areas map[int]int) bool {
	for k, v := range areas {
		posKey := m.GetPosKey(k, v)
		if _, ok := m.mapEmpityPos[posKey]; !ok {
			return false
		}
	}
	return true
}

func (m *MapObj) AddResourceRepCap(delta int) {
	m.resourceRepCap += delta
}

func (m *MapObj) SetHappiness(values common.ItemInfos) {
	m.happyV += configs.GetDb().GetHappinesV(values)
}

func (m *MapObj) DelHappiness(values common.ItemInfos) {
	m.happyV -= configs.GetDb().GetHappinesV(values)
}

func (m *MapObj) GetHappiness() int {
	return m.happyV
}

func (m *MapObj) setEmpityByUnlockedArea(areaId int) {
	areaCfg := configs.GetDb().GetAreaMap(m.cityId, areaId)
	if areaCfg == nil {
		return
	}
	for _, posStr := range areaCfg.AreaMap {
		m.mapEmpityPos[posStr] = true
	}
	log.Info("=============setMapEmpityPos m.mapEmpityPos[%d]=%v", m.cityId, m.mapEmpityPos)
	log.Info("len =%d", len(m.mapEmpityPos))
}

func (m *MapObj) delMapEmpity(buildingPos *building.BuildingPos) {
	areas := m.GetPosAreaByIdPos(buildingPos.Id, buildingPos.Pos)
	for k, v := range areas {
		posKey := m.GetPosKey(k, v)
		delete(m.mapEmpityPos, posKey)
	}
}

func (m *MapObj) addMapEmpity(buildingPos *building.BuildingPos) {
	areas := m.GetPosAreaByIdPos(buildingPos.Id, buildingPos.Pos)
	for k, v := range areas {
		posKey := m.GetPosKey(k, v)
		m.mapEmpityPos[posKey] = true
	}
}

func (m *MapObj) setMapEmpityPos() {
	//set
	for _, areaId := range m.unlockedArea {
		m.setEmpityByUnlockedArea(areaId)
	}

	//del building pos
	for tp, v := range m.buldingInfos {
		// switch tp {
		// case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
		// 	a := m.buldingInfos[tp].(*building.BaseBuilding)
		// 	for _, v := range a.Buildings {
		// 		log.Info("************ before  setMapEmpityPos  BaseBuilding building=%v", v)
		// 	}
		// case building.PRODUCTION_BUILDING:
		// 	a := m.buldingInfos[tp].(*building.ProductionBuilding)
		// 	for _, v := range a.Buildings {
		// 		log.Info("*********** before  setMapEmpityPos ProductionBuilding building=%v", v)
		// 	}
		// case building.OCCUP_BUILDING:
		// 	a := m.buldingInfos[tp].(*building.OccupBuilding)
		// 	for _, v := range a.Buildings {
		// 		log.Info("*********** before  setMapEmpityPos OccupBuilding building=%v", v)
		// 	}
		// }
		buildingsPos := v.GetAllPos()
		for _, buildingPos := range buildingsPos {
			log.Info("---------------setMapEmpityPos buildingPos=%v", tp, buildingPos)
			m.delMapEmpity(buildingPos)
			log.Info("=========del cityId=%d,buildingPos=%v,mapEmpityPos=%v", m.cityId, buildingPos, m.mapEmpityPos)
		}
		// switch tp {
		// case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
		// 	a := m.buldingInfos[tp].(*building.BaseBuilding)
		// 	for _, v := range a.Buildings {
		// 		log.Info("************ after  setMapEmpityPos  BaseBuilding building=%v", v)
		// 	}
		// case building.PRODUCTION_BUILDING:
		// 	a := m.buldingInfos[tp].(*building.ProductionBuilding)
		// 	for _, v := range a.Buildings {
		// 		log.Info("*********** after  setMapEmpityPos ProductionBuilding building=%v", v)
		// 	}
		// case building.OCCUP_BUILDING:
		// 	a := m.buldingInfos[tp].(*building.OccupBuilding)
		// 	for _, v := range a.Buildings {
		// 		log.Info("*********** after  setMapEmpityPos OccupBuilding building=%v", v)
		// 	}
		// }
	}

}

func (m *MapObj) GetPosAreaByIdPos(id int, pos common.PosInfo) (areas map[int]int) {
	areas = make(map[int]int)
	log.Info("========MapObj GetPosAreaByIdPos id=%d,pos=%d", id, pos)
	buildCfg := configs.GetDb().GetBuild(id)
	if buildCfg == nil {
		log.Info("========MapObj GetPosAreaByIdPos 11")
		return
	}
	moduleCfg := configs.GetDb().GetBuildingModel(buildCfg.NormalModel)
	if moduleCfg == nil {
		log.Info("========MapObj GetPosAreaByIdPos 222222 buildCfg.NormalModel=%v", buildCfg.NormalModel)
		return
	}
	if len(moduleCfg.Size) < 2 {
		log.Info("========MapObj GetPosAreaByIdPos 222222 moduleCfg.Size=%v", moduleCfg.Size)
		return
	}
	if pos.Direct == 0 {
		for i := pos.PosX; i < pos.PosX+moduleCfg.Size[0]; i++ {
			for j := pos.PoxY; j > pos.PoxY-moduleCfg.Size[1]; j-- {
				areas[i] = j
			}
		}
	} else {
		for i := pos.PosX; i < pos.PosX-moduleCfg.Size[1]; i++ {
			for j := pos.PoxY; j > pos.PoxY-moduleCfg.Size[0]; j-- {
				areas[i] = j
			}
		}
	}
	log.Info("=====================area=%v", areas)
	return
}

func (m *MapObj) SetMinjuOccup(user IMapUser, cityId, buildingLvl int, buildingInfo *building.BuildingInfo, isInit, isBuilding bool) {
	currentMaxNpc := 2
	npcList := make(common.IntSlice, 0)
	if isBuilding {
		minjuCount := m.buldingInfos[building.OCCUP_BUILDING].(*building.OccupBuilding).GetCountByItemType(building.BUILD_ITEMTYPE_MINJU)
		log.Info("-----------FinishBuild-------minjuCount=%v", minjuCount)
		minjuCfg := configs.GetDb().GetBuildingMinju(minjuCount)
		log.Info("-----------FinishBuild-------minjuCfg=%v", minjuCfg)
		if minjuCfg == nil {
			return
		}
		npcList = minjuCfg.Npc
	} else {
		if buildingLvl == 4 {
			currentMaxNpc = 3
		} else if buildingLvl == 5 {
			currentMaxNpc = 4
		}
		if buildingInfo.Occupiers[1] == 0 {
			return
		}
		npcList = configs.GetDb().GetNpcByOne(buildingInfo.Occupiers[1])
	}
	var index int
	for k, v := range buildingInfo.Occupiers {
		if v == 0 && k != 0 {
			index = k
			break
		}
	}
	if index == 0 || index > currentMaxNpc {
		return
	}
	for k, v := range npcList {
		if k < index-1 {
			continue
		}
		if k == currentMaxNpc {
			break
		}
		peopleR := user.GetPeopleBag().Insert(v)
		peopleR.CityID = cityId
		buildingInfo.Occupiers[k+1] = v
		log.Info("-----------FinishBuild-------v=%v ,", v)
		log.Info("buildingInfo.Occupiers=%v", buildingInfo.Occupiers)
	}
}

func (m *MapObj) FinishBuild(user IMapUser, uid int) (*building.BuildingInfo, *bags.HeroMD, []*bags.PeopleMD, error) {
	var buildingInfo *building.BuildingInfo
	var heromd *bags.HeroMD
	var err error
	currentBuilding := m.buldingInfosByUid[uid]
	if currentBuilding == nil {
		return nil, nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(currentBuilding.Id)
	if buildCfg == nil {
		return nil, nil, nil, errorx.Build_Item_Not_Error
	}
	fmt.Println("====FinishBuild currentBuilding=", currentBuilding)
	if currentBuilding.Builder > 0 {
		heromd, err = bags.GetHeroMDBycityID(user, m.cityId, currentBuilding.Builder)
		if err != nil {
			return nil, nil, nil, err
		}
		heromd.AddExp(buildCfg.BulidExp)
		heromd.Status = bags.HeroStatus_Default
	}
	buildingInfo, err = m.buldingInfos[buildCfg.TitleType].EndBuild(user, uid)
	if err != nil {
		return nil, nil, nil, err
	}
	//update hero

	//update peoples
	bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Default, currentBuilding.BuilderPeople...)
	if buildCfg.Itemtype == building.BUILD_ITEMTYPE_MINJU {
		m.SetMinjuOccup(user, buildCfg.City, buildCfg.Level, buildingInfo, false, true)
	}
	m.SetHappiness(buildCfg.BuildValue)
	currentBuilding = buildingInfo
	peoples := user.GetPeopleBag().GetPeoplesByUID(currentBuilding.BuilderPeople)
	return buildingInfo, heromd, peoples, err
}

func (m *MapObj) BuildByInit(user IMapUser, id int, pos common.PosInfo) (err error) {
	log.Info("MapObj Build id=%d", id)
	buildcfg := configs.GetDb().GetBuild(id)
	if buildcfg == nil {
		return errorx.Build_Item_Not_Error
	}

	if !building.IsBuildingType(buildcfg.TitleType) {
		return errorx.Build_Item_Not_Error
	}
	var buildingInfo *building.BuildingInfo

	//check condition unLockCondition,buildCondition

	//check map
	// areas := m.GetPosAreaByIdPos(id, pos)
	// if !m.CheckUsedMapPos(areas) {
	// 	return nil, nil, errorx.Build_Item_No_Space
	// }
	//check num
	if m.buldingInfos[buildcfg.TitleType] == nil {
		m.buldingInfos[buildcfg.TitleType] = building.NewIBuilding(buildcfg.TitleType)
	}

	switch buildcfg.TitleType {
	case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildcfg.TitleType]).(*building.BaseBuilding).InitBuild(user, id, pos)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildcfg.TitleType]).(*building.ProductionBuilding).InitBuild(user, id, pos)
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildcfg.TitleType]).(*building.OccupBuilding).InitBuild(user, id, pos)
	default:
		err = errorx.Release_Building_No_Err
		return
	}
	if err != nil {
		return
	}
	if buildcfg.Itemtype == building.BUILD_ITEMTYPE_MINJU {
		m.SetMinjuOccup(user, buildcfg.City, buildcfg.Level, buildingInfo, true, true)
	}
	m.SetHappiness(buildcfg.BuildValue)
	m.buldingInfosByUid[buildingInfo.Uid] = buildingInfo
	m.delMapEmpity(&building.BuildingPos{Id: buildingInfo.Id, Pos: buildingInfo.Pos})
	return
}

func (m *MapObj) Build(user IMapUser, id, heroId int, peopleIds []int, pos common.PosInfo, isWatchVideos bool, makeUpGold int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	log.Info("MapObj Build id=%d,isWatchVideos=%v,heroId=%d,peopleIds=%v", id, isWatchVideos, heroId, peopleIds)
	buildcfg := configs.GetDb().GetBuild(id)
	if buildcfg == nil {
		return nil, nil, nil, nil, errorx.Build_Item_Not_Error
	}
	if !building.IsBuildingType(buildcfg.TitleType) {
		return nil, nil, nil, nil, errorx.Build_Item_Not_Error
	}
	//check condition unLockCondition,buildCondition

	//check hero
	//var needPeople, useBuildingTime int
	var heroBuildingResult *bags.MakeResult
	//var heromd *bags.HeroMD
	log.Info("MapObj Build buildcfg.NeedHero=%d", buildcfg.NeedHero)
	if buildcfg.NeedHero > 0 {
		if heroId <= 0 {
			return nil, nil, nil, nil, errorx.Build_People_Not_Error
		}
		//check hero bag
		heromd, err = bags.GetHeroMDBycityID(user, buildcfg.City, heroId)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		heroBuildingResult = heromd.CreateBuilding(buildcfg)
	} else {
		heroBuildingResult = bags.NewMakeResult(nil, buildcfg, configs.WorkTypeEnum_CreateBuilding,
			buildcfg.PeopleNum, buildcfg.BuildTime, buildcfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	}
	useBuildingTime := heroBuildingResult.GetUseTime()
	needPeople := heroBuildingResult.GetNpcNum()
	buildCost := heroBuildingResult.GetItems()
	//check people{}
	peoples = make([]*bags.PeopleMD, 0)
	if isWatchVideos == true {
		needPeople--
	}
	if needPeople > 0 {
		var isExist bool
		peoples, isExist = user.GetPeopleBag().GetPeopleByStatusID(bags.HeroStatus_Default, buildcfg.City, needPeople)
		if !isExist {
			return nil, nil, nil, nil, errorx.Build_People_Not_Error
		}
	}

	//check map
	// areas := m.GetPosAreaByIdPos(id, pos)
	// if !m.CheckUsedMapPos(areas) {
	// 	return nil, nil, errorx.Build_Item_No_Space
	// }
	//check num
	if m.buldingInfos[buildcfg.TitleType] == nil {
		m.buldingInfos[buildcfg.TitleType] = building.NewIBuilding(buildcfg.TitleType)
	}
	if m.buldingInfos[buildcfg.TitleType].CheckLimit(user, id) {
		return nil, nil, nil, nil, errorx.Build_Item_Num_Error.CloneWithArgs(buildcfg.Name)
	}
	//check cost item
	err = bag.BagEx.CheckItem(user, buildCost)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//del item
	resultItems, err = bag.BagEx.DelItems(user, buildCost, 0)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	//m.DelUsedPosFromLocked(areas)

	//additem
	extraAward := heroBuildingResult.GetReward()
	addResult := bag.BagEx.AddItems(user, extraAward, 0)
	resultItems = append(resultItems, addResult...)
	//update hero status
	if heromd != nil {
		heromd.Status = bags.HeroStatus_Working
	}

	//update people status
	buildPeopleIds := make([]int, len(peoples))
	for k, v := range peoples {
		buildPeopleIds[k] = v.PeopleID
		v.Status = bags.HeroStatus_Working
	}
	bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Working, buildPeopleIds...)

	//update building
	buildingInfo, err = m.buldingInfos[buildcfg.TitleType].StartBuild(user, id, heroId, useBuildingTime, buildPeopleIds, pos)
	m.buldingInfosByUid[buildingInfo.Uid] = buildingInfo
	m.delMapEmpity(&building.BuildingPos{Id: buildingInfo.Id, Pos: buildingInfo.Pos})
	return
}

func (m *MapObj) StartProduce(user IMapUser, uid, heroId int, peopleIds []int, isWatchVideos bool, makeUpGold int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	buildingTemp := m.buldingInfosByUid[uid]
	if buildingTemp == nil {
		return nil, nil, nil, nil, errorx.Build_Item_User_Not_Error
	}
	log.Info("=====MapObj StartProduce buildingTemp=%v", buildingTemp)
	buildCfg := configs.GetDb().GetBuild(buildingTemp.Id)
	if buildCfg == nil {
		return nil, nil, nil, nil, errorx.Build_Item_Not_Error
	}

	//condition produce

	//check hero
	log.Info("MapObj StartProduce buildcfg.NeedHero=%d", buildCfg.NeedHero)
	var heroProduceResult *bags.MakeResult
	if buildCfg.NeedHero > 0 {
		if heroId <= 0 {
			return nil, nil, nil, nil, errorx.Build_People_Not_Error
		}
		//check hero bag
		heromd, err = bags.GetHeroMDBycityID(user, buildCfg.City, heroId)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		heroProduceResult = heromd.Working(buildCfg)
	} else {
		heroProduceResult = bags.NewMakeResult(nil, buildCfg, configs.WorkTypeEnum_Working,
			buildCfg.WorkHero, buildCfg.WorkTime, buildCfg.WorkCost, common.ItemInfos{}, buildCfg.WorkOut)
	}
	useProduceTime := heroProduceResult.GetUseTime()
	needPeople := heroProduceResult.GetNpcNum()
	buildCost := heroProduceResult.GetItems()

	//check people
	peoples = make([]*bags.PeopleMD, 0)
	if isWatchVideos == true {
		needPeople--
	}
	if needPeople > 0 {
		// if needPeople > len(peopleIds) {
		// 	return nil, nil, nil, nil, errorx.Build_People_Not_Error
		// }
		var isExist bool
		peoples, isExist = user.GetPeopleBag().GetPeopleByStatusID(bags.HeroStatus_Default, buildCfg.City, needPeople)
		if !isExist {
			return nil, nil, nil, nil, errorx.Build_People_Not_Error
		}
	}

	//check item
	err = bag.BagEx.CheckItem(user, buildCost)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//del item
	resultItems, err = bag.BagEx.DelItems(user, buildCost, 0)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//additem
	extraAward := heroProduceResult.GetReward()
	fmt.Printf("StartProduce:%+v.\n", extraAward)
	addResultItems := bag.BagEx.AddItems(user, extraAward, 0)
	resultItems = append(resultItems, addResultItems...)
	//update people
	buildPeopleIds := make([]int, len(peoples))
	for k, v := range peoples {
		buildPeopleIds[k] = v.PeopleID
		v.Status = bags.HeroStatus_Working
	}
	bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Working, buildPeopleIds...)

	if heromd != nil {
		heromd.Status = bags.HeroStatus_Working
	}
	switch buildCfg.TitleType {
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).AwardStart(uid, heroId, useProduceTime, buildPeopleIds)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).AwardStart(uid, heroId, useProduceTime, buildPeopleIds)
	default:
		err = errorx.Release_Building_No_Err
		return
	}
	if err != nil {
		return
	}
	buildingTemp = buildingInfo
	return
}

func (m *MapObj) FinishProduce(user IMapUser, uid int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, nil, nil, errorx.Build_Item_Not_Error
	}

	//buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).AwardEnd(uid) //EndBuild(user, uid)
	//update hero
	var heroProduceResult *bags.MakeResult
	if buildingObj.Builder > 0 {
		heromd, err = bags.GetHeroMDBycityID(user, buildCfg.City, buildingObj.Builder)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		heromd.AddExp(buildCfg.BulidExp)
		heromd.Status = bags.HeroStatus_Default
		heroProduceResult = heromd.Working(buildCfg)
	} else {
		heroProduceResult = bags.NewMakeResult(nil, buildCfg, configs.WorkTypeEnum_Working,
			buildCfg.WorkHero, buildCfg.WorkTime, buildCfg.WorkCost, common.ItemInfos{}, buildCfg.WorkOut)
	}
	switch buildCfg.TitleType {
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).AwardEnd(uid)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).AwardEnd(uid)
	default:
		err = errorx.Release_Building_No_Err
		return
	}
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//update item
	awardItms := heroProduceResult.GetResLi()
	fmt.Printf("FinishProduce:%+v.\n", awardItms.ToInt32Map())
	resultItems = bag.BagEx.AddItems(user, awardItms, 0)

	//update hero
	if heromd != nil {
		heromd.Status = bags.HeroStatus_Default
	}

	//update peoples
	bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Default, buildingObj.BuilderPeople...)
	buildingObj = buildingInfo
	peoples = user.GetPeopleBag().GetPeoplesByUID(buildingObj.BuilderPeople)
	return
}

func (m *MapObj) StartUpLevel(user IMapUser, uid, heroId int, peopleIds []int, isWatchVideos bool, makeUpGold int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, resultItems []bag.IBagItem, peoples []*bags.PeopleMD, err error) {
	fmt.Println("====MapObj StartUpLevel uid=", uid)
	buildingTemp := m.buldingInfosByUid[uid]
	if buildingTemp == nil {
		return nil, nil, nil, nil, errorx.Build_Item_User_Not_Error
	}
	fmt.Println("====MapObj StartUpLevel buildingTemp=", buildingTemp)
	fmt.Println("====MapObj StartUpLevel buildingTemp.Id=", buildingTemp.Id)
	buildCfg := configs.GetDb().GetBuild(buildingTemp.Id)
	if buildCfg == nil {
		return nil, nil, nil, nil, errorx.Build_Item_Not_Error
	}
	var newBuildingId int
	newBuildCfg := configs.GetDb().GetBuildByLevel(buildCfg.Itemtype, buildCfg.Level+1)
	if newBuildCfg == nil {
		return nil, nil, nil, nil, errorx.Build_Item_Not_Error
	}
	newBuildingId = newBuildCfg.Id
	//condition produce

	//check hero
	log.Info("MapObj StartProduce newBuildCfg.NeedHero=%d", newBuildCfg.NeedHero)
	var heroUpLevelResult *bags.MakeResult
	if newBuildCfg.NeedHero > 0 {
		if heroId <= 0 {
			return nil, nil, nil, nil, errorx.Build_People_Not_Error
		}
		//check hero bag
		heromd, err = bags.GetHeroMDBycityID(user, newBuildCfg.City, heroId)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		heroUpLevelResult = heromd.UpBuilding(newBuildCfg)
	} else {
		heroUpLevelResult = bags.NewMakeResult(nil, newBuildCfg, configs.WorkTypeEnum_UpBuilding,
			newBuildCfg.PeopleNum, newBuildCfg.BuildTime, newBuildCfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	}
	useUpLevelTime := heroUpLevelResult.GetUseTime()
	needPeople := heroUpLevelResult.GetNpcNum()
	upLevelCost := heroUpLevelResult.GetItems()

	//check people
	peoples = make([]*bags.PeopleMD, 0)
	if isWatchVideos == true {
		needPeople--
	}
	if needPeople > 0 {
		// if needPeople > len(peopleIds) {
		// 	return nil, nil, nil, nil, errorx.Build_People_Not_Error
		// }
		var isExist bool
		peoples, isExist = user.GetPeopleBag().GetPeopleByStatusID(bags.HeroStatus_Default, newBuildCfg.City, needPeople)
		if !isExist {
			return nil, nil, nil, nil, errorx.Build_People_Not_Error
		}
	}

	//check item
	err = bag.BagEx.CheckItem(user, upLevelCost)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//del item
	resultItems, err = bag.BagEx.DelItems(user, upLevelCost, 0)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//additem
	extraAward := heroUpLevelResult.GetReward()
	addReslut := bag.BagEx.AddItems(user, extraAward, 0)
	resultItems = append(resultItems, addReslut...)
	//update people
	buildPeopleIds := make([]int, len(peoples))
	for k, v := range peoples {
		buildPeopleIds[k] = v.PeopleID
		v.Status = bags.HeroStatus_Working
	}
	bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Working, buildPeopleIds...)

	if heromd != nil {
		heromd.Status = bags.HeroStatus_Working
	}
	switch buildCfg.TitleType {
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[newBuildCfg.TitleType]).(*building.OccupBuilding).LvlUpStartBuild(uid, newBuildingId, heroId, useUpLevelTime, buildPeopleIds)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[newBuildCfg.TitleType]).(*building.ProductionBuilding).LvlUpStartBuild(uid, newBuildingId, heroId, useUpLevelTime, buildPeopleIds)
	default:
		err = errorx.LevelUp_Building_No_Err
		return
	}
	if err != nil {
		return
	}
	buildingTemp = buildingInfo

	return
}

func (m *MapObj) FinishLevelUp(user IMapUser, uid int) (buildingInfo *building.BuildingInfo, heromd *bags.HeroMD, peoples []*bags.PeopleMD, err error) {
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, nil, errorx.Build_Item_Not_Error
	}

	//buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).LvlUpEndBuild(uid) //EndBuild(user, uid)
	//update hero
	if buildingObj.Builder > 0 {
		heromd, err = bags.GetHeroMDBycityID(user, buildCfg.City, buildingObj.Builder)
		if err != nil {
			return nil, nil, nil, err
		}
		heromd.AddExp(buildCfg.BulidExp)
		heromd.Status = bags.HeroStatus_Default
	} else {
		//heroLevelUpResult = bags.NewMakeResult(nil, buildCfg, configs.WorkTypeEnum_UpBuilding,
		//	buildCfg.PeopleNum, buildCfg.BuildTime, buildCfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	}
	switch buildCfg.TitleType {
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).LvlUpEndBuild(uid)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).LvlUpEndBuild(uid)
	default:
		err = errorx.LevelUp_Building_No_Err
		return
	}
	if err != nil {
		return
	}
	//update hero
	if heromd != nil {
		heromd.Status = bags.HeroStatus_Default
	}

	//update peoples
	bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Default, buildingObj.BuilderPeople...)

	m.SetHappiness(buildCfg.BuildValue)
	oldBuildCfg := configs.GetDb().GetBuildByLevel(buildCfg.Itemtype, buildCfg.Level-1)
	if oldBuildCfg != nil {
		m.DelHappiness(buildCfg.BuildValue)
	}
	if buildCfg.Itemtype == building.BUILD_ITEMTYPE_MINJU {
		m.SetMinjuOccup(user, buildCfg.City, buildCfg.Level, buildingInfo, false, false)
	}
	buildingObj = buildingInfo
	peoples = user.GetPeopleBag().GetPeoplesByUID(buildingObj.BuilderPeople)
	return
}

func (m *MapObj) getBuildingEndMakeUpGold(buildingInfo *building.BuildingInfo) (int, int) {
	needTime := buildingInfo.BuildCompletedTime - int(time.Now().Unix())
	needMin := int(math.Ceil(float64(needTime) / 60.0))
	return needMin * configs.GetDb().GameEx.CoinSpeedTimes, needTime * 60
}

func (m *MapObj) BuildingSpeedUp(user IMapUser, uid, costType int) (*building.BuildingInfo, []bag.IBagItem, error) {
	var buildingInfo *building.BuildingInfo
	var err error
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, errorx.Build_Item_Not_Error
	}

	//check cost
	costItems := make(common.ItemInfos, 0)
	costItem := &common.ItemInfo{}
	var speedTime int
	switch costType {
	case building.SPEEDUPCOSTTYPE_WATCH_VIDEO:
		speedTime = configs.GetDb().GameEx.VideoReduceTimes * 60
	case building.SPEEDUPCOSTTYPE_GOLD:
		costItem.ItemId = bag.ItemType_Gold
		costItem.Count, speedTime = m.getBuildingEndMakeUpGold(buildingObj)
	case building.SPEEDUPCOSTTYPE_CARD:
		costItem.ItemId = configs.ItemType_QuickCard
		costItem.Count = 1
		speedTime = configs.GetDb().GameEx.AddSpeedCardTimes * 60
	default:
		return nil, nil, errorx.Building_Speed_up_Type_Error
	}
	if costItem.ItemId != 0 && costItem.Count != 0 {
		costItems = append(costItems, costItem)
	}

	err = bag.BagEx.CheckItem(user, costItems)
	if err != nil {
		return nil, nil, err
	}

	//cost
	resultItems, err := bag.BagEx.DelItems(user, costItems, 0)
	if err != nil {
		return nil, nil, err
	}
	switch buildCfg.TitleType {
	case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.BaseBuilding).BuildingSpeedUp(uid, speedTime)
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).BuildingSpeedUp(uid, speedTime)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).BuildingSpeedUp(uid, speedTime)
	}
	if err != nil {
		return nil, nil, err
	}
	buildingObj = buildingInfo
	return buildingInfo, resultItems, nil
}

func (m *MapObj) getLevelUpEndMakeUpGold(buildingInfo *building.BuildingInfo) (int, int) {
	needTime := buildingInfo.UpLevelCompletedTime - int(time.Now().Unix())
	needMin := int(math.Ceil(float64(needTime) / 60.0))
	return needMin * configs.GetDb().GameEx.CoinSpeedTimes, needTime * 60
}

func (m *MapObj) LevelUpSpeedUp(user IMapUser, uid, costType int) (*building.BuildingInfo, []bag.IBagItem, error) {
	var buildingInfo *building.BuildingInfo
	var err error
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, errorx.Build_Item_Not_Error
	}

	//check cost
	costItems := make(common.ItemInfos, 0)
	costItem := &common.ItemInfo{}
	var speedTime int
	switch costType {
	case building.SPEEDUPCOSTTYPE_WATCH_VIDEO:
		speedTime = configs.GetDb().GameEx.VideoReduceTimes * 60
	case building.SPEEDUPCOSTTYPE_GOLD:
		costItem.ItemId = bag.ItemType_Gold
		costItem.Count, speedTime = m.getLevelUpEndMakeUpGold(buildingObj)
	case building.SPEEDUPCOSTTYPE_CARD:
		costItem.ItemId = configs.ItemType_QuickCard
		costItem.Count = 1
		speedTime = configs.GetDb().GameEx.AddSpeedCardTimes * 60
	default:
		return nil, nil, errorx.Building_Speed_up_Type_Error
	}
	if costItem.ItemId != 0 && costItem.Count != 0 {
		costItems = append(costItems, costItem)
	}

	err = bag.BagEx.CheckItem(user, costItems)
	if err != nil {
		return nil, nil, err
	}

	//cost
	resultItems, err := bag.BagEx.DelItems(user, costItems, 0)
	if err != nil {
		return nil, nil, err
	}
	switch buildCfg.TitleType {
	case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.BaseBuilding).LevelUpSpeedUp(uid, speedTime)
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).LevelUpSpeedUp(uid, speedTime)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).LevelUpSpeedUp(uid, speedTime)
	}
	if err != nil {
		return nil, nil, err
	}
	buildingObj = buildingInfo
	return buildingInfo, resultItems, nil
}

func (m *MapObj) getProduceEndMakeUpGold(buildingInfo *building.BuildingInfo) (int, int) {
	needTime := buildingInfo.ProduceCompletedTime - int(time.Now().Unix())
	needMin := int(math.Ceil(float64(needTime) / 60.0))
	return needMin * configs.GetDb().GameEx.CoinSpeedTimes, needTime * 60
}

func (m *MapObj) ProduceSpeedUp(user IMapUser, uid, costType int) (*building.BuildingInfo, []bag.IBagItem, error) {
	var buildingInfo *building.BuildingInfo
	var err error
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, errorx.Build_Item_Not_Error
	}

	//check cost
	costItems := make(common.ItemInfos, 0)
	costItem := &common.ItemInfo{}
	var speedTime int
	switch costType {
	case building.SPEEDUPCOSTTYPE_WATCH_VIDEO:
		speedTime = configs.GetDb().GameEx.VideoReduceTimes * 60
	case building.SPEEDUPCOSTTYPE_GOLD:
		costItem.ItemId = bag.ItemType_Gold
		costItem.Count, speedTime = m.getProduceEndMakeUpGold(buildingObj)
	case building.SPEEDUPCOSTTYPE_CARD:
		costItem.ItemId = configs.ItemType_QuickCard
		costItem.Count = 1
		speedTime = configs.GetDb().GameEx.AddSpeedCardTimes * 60
	default:
		return nil, nil, errorx.Building_Speed_up_Type_Error
	}
	if costItem.ItemId != 0 && costItem.Count != 0 {
		costItems = append(costItems, costItem)
	}

	err = bag.BagEx.CheckItem(user, costItems)
	if err != nil {
		return nil, nil, err
	}

	//cost
	resultItem, err := bag.BagEx.DelItems(user, costItems, 0)
	if err != nil {
		return nil, nil, err
	}

	switch buildCfg.TitleType {
	case building.DECORATE_BUILDING, building.AMUSEMENT_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.BaseBuilding).ProduceSpeedUp(uid, speedTime)
	case building.OCCUP_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).ProduceSpeedUp(uid, speedTime)
	case building.PRODUCTION_BUILDING:
		buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.ProductionBuilding).ProduceSpeedUp(uid, speedTime)
	}
	if err != nil {
		return nil, nil, err
	}

	buildingObj = buildingInfo
	return buildingInfo, resultItem, nil
}

func (m *MapObj) MoveBuilding(user IMapUser, uid int, newPos common.PosInfo) (*building.BuildingInfo, error) {
	var buildingInfo *building.BuildingInfo
	var err error
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, errorx.Build_Item_Not_Error
	}
	//check map
	// areas := m.GetPosAreaByIdPos(id, pos)
	// if !m.CheckUsedMapPos(areas) {
	// 	return nil, nil, errorx.Build_Item_No_Space
	// }

	buildingInfo, err = (m.buldingInfos[buildCfg.TitleType]).(*building.BaseBuilding).MoveBuilding(uid, newPos)
	if err != nil {
		return nil, err
	}
	buildingObj = buildingInfo
	return buildingInfo, nil
}

func (m *MapObj) AwardMinju(user IMapUser, uid int) (*building.BuildingInfo, []bag.IBagItem, error) {
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, errorx.Build_Item_Not_Error
	}
	buildingInfo, awards, err := (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).AwardMinju(uid) //EndBuild(user, uid)
	if err != nil {
		return nil, nil, err
	}
	//update item
	resultItems := bag.BagEx.AddItems(user, *awards, 0)

	buildingObj = buildingInfo
	return buildingInfo, resultItems, nil
}

func (m *MapObj) OccupMinju(user IMapUser, uid, heroId int) (*building.BuildingInfo, []*bags.HeroMD, error) {
	heromds := make([]*bags.HeroMD, 0)
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, errorx.Build_Item_Not_Error
	}
	if buildCfg.Itemtype != building.BUILD_ITEMTYPE_MINJU {
		return nil, nil, errorx.Build_Type_Occup_Error
	}

	buildingInfo, oldOccuperId, err := (m.buldingInfos[buildCfg.TitleType]).(*building.OccupBuilding).OccupBuilding(uid, heroId) //EndBuild(user, uid)
	if err != nil {
		return nil, nil, err
	}
	heromd, err := bags.GetHeroMD(user, heroId)
	if err != nil {
		return nil, nil, err
	}
	heromd.IsHome = true
	heromds = append(heromds, heromd)
	if oldOccuperId != 0 {
		oldHeromd, err := bags.GetHeroMD(user, oldOccuperId)
		if err != nil {
			return nil, nil, err
		}
		oldHeromd.IsHome = false
		heromds = append(heromds, oldHeromd)
	}

	heromd.IsHome = true

	buildingObj = buildingInfo
	return buildingInfo, heromds, err
}

func (m *MapObj) ReleasePeople(user IMapUser, uid, releaseType int) (*building.BuildingInfo, *bags.HeroMD, []*bags.PeopleMD, error) {
	var heromd *bags.HeroMD
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return nil, nil, nil, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		return nil, nil, nil, errorx.Build_Item_Not_Error
	}
	var buildingInfo *building.BuildingInfo
	var err error
	var heroId int
	buildPeopleIds := make([]int, 0)
	buildingInfo, heroId, buildPeopleIds, err = (m.buldingInfos[buildCfg.TitleType]).ReleasePeople(uid, releaseType)

	if err != nil {
		return nil, nil, nil, err
	}

	if heroId != 0 {
		heromd, err = bags.GetHeroMD(user, heroId)
		if err != nil {
			return nil, nil, nil, err
		}
		heromd.Status = bags.HeroStatus_Default
	}
	if len(buildPeopleIds) > 0 {
		bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Default, buildPeopleIds...)
	}
	buildingObj = buildingInfo
	peoples := user.GetPeopleBag().GetPeoplesByUID(buildPeopleIds)
	return buildingInfo, heromd, peoples, err
}

func (m *MapObj) DelBuilding(user IMapUser, uid int) (int, error) {
	buildingObj := m.buldingInfosByUid[uid]
	if buildingObj == nil {
		return 0, errorx.Build_Item_User_Not_Error
	}
	buildCfg := configs.GetDb().GetBuild(buildingObj.Id)
	if buildCfg == nil {
		log.Info("========MapObj GetPosAreaByIdPos 11")
		return 0, errorx.Build_Item_Not_Error
	}
	delUid, err := (m.buldingInfos[buildCfg.TitleType]).DelteBuild(user, uid)
	if err != nil {
		return 0, err
	}
	m.addMapEmpity(&building.BuildingPos{Id: buildCfg.Id, Pos: buildingObj.Pos})
	delete(m.buldingInfosByUid, uid)
	return delUid, err
}

func (m *MapObj) getExploreTimes() int {
	return len(m.unlockedArea) + len(m.exploreItems) + len(configs.GetDb().GameEx.FirstArea) + 1
}

func (m *MapObj) ExploreArea(user IMapUser, areaId, heroId, titleX, titleY, buildId int, peopleIds common.IntSlice) ([]bag.IBagItem, common.ExploreItem, *bags.HeroMD, []*bags.PeopleMD, error) {
	resultItems := make([]bag.IBagItem, 0)
	var exploreItem common.ExploreItem
	//check explore
	for _, v := range m.exploreItems {
		if v.Id == areaId {
			return nil, exploreItem, nil, nil, errorx.Building_Explore_Err
		}
	}

	//check unlockedArea
	for _, v := range m.unlockedArea {
		if v == areaId {
			return nil, exploreItem, nil, nil, errorx.Building_UnLock_Err
		}
	}

	//check unlockTimes
	unlockTimes := m.getExploreTimes() //len(m.unlockedArea) + len(m.exploreItems) + 1
	areaCostCfg := configs.GetDb().GetAreaOpenCostCfg(unlockTimes)
	if areaCostCfg == nil {
		return nil, exploreItem, nil, nil, errorx.UnLock_Area_No_Err
	}

	var heromd *bags.HeroMD
	var err error
	var heroExploreResult *bags.MakeResult
	if areaCostCfg.NeedHero > 0 {
		if heroId <= 0 {
			return nil, exploreItem, nil, nil, errorx.Build_People_Not_Error
		}
		//check hero bag
		heromd, err = bags.GetHeroMDBycityID(user, areaCostCfg.City, heroId)
		if err != nil {
			return nil, exploreItem, nil, nil, err
		}
		heroExploreResult = heromd.ExploreNew(areaCostCfg)
	} else {
		heroExploreResult = bags.NewMakeResult(areaCostCfg, nil, configs.WorkTypeEnum_ExploreNew,
			areaCostCfg.PeopleNum, areaCostCfg.UseTime, areaCostCfg.Cost, common.ItemInfos{}, common.ItemInfos{})
	}
	exploreTime := heroExploreResult.GetUseTime()
	needPeople := heroExploreResult.GetNpcNum()
	exploreCost := heroExploreResult.GetItems()

	//check people
	peoples := make([]*bags.PeopleMD, 0)
	// if needPeople > len(peopleIds) {
	// 	return nil, exploreItem, nil, nil, errorx.Build_People_Not_Error
	// }
	var isExist bool
	peoples, isExist = user.GetPeopleBag().GetPeopleByStatusID(bags.HeroStatus_Default, areaCostCfg.City, needPeople)
	if !isExist {
		return nil, exploreItem, nil, nil, errorx.Build_People_Not_Error
	}

	//check condition

	//check cost
	err = bag.BagEx.CheckItem(user, exploreCost)
	if err != nil {
		return nil, exploreItem, nil, nil, err
	}
	//del item
	resultItems, err = bag.BagEx.DelItems(user, exploreCost, 0)
	if err != nil {
		return nil, exploreItem, nil, nil, err
	}

	//add extaraward
	extraAward := heroExploreResult.GetReward()
	addResult := bag.BagEx.AddItems(user, extraAward, 0)
	resultItems = append(resultItems, addResult...)

	buildPeopleIds := make([]int, len(peoples))
	for k, v := range peoples {
		buildPeopleIds[k] = v.PeopleID
		v.Status = bags.HeroStatus_Working
	}
	exploreItem = common.ExploreItem{
		Id:        areaId,
		EndTime:   int(time.Now().Unix()) + exploreTime,
		Hero:      heroId,
		WorkerIds: buildPeopleIds,
		TileX:     titleX,
		TileY:     titleY,
		BuildId:   buildId,
	}

	//update hero status
	if heromd != nil {
		heromd.Status = bags.HeroStatus_Working
	}

	//update people status
	bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Working, buildPeopleIds...)

	//update explore
	m.exploreItems = append(m.exploreItems, exploreItem)
	return resultItems, exploreItem, heromd, peoples, err
}

func (m *MapObj) UnlockArea(user IMapUser, cityId, areaId int, isUseGold bool) ([]bag.IBagItem, *bags.HeroMD, []*bags.PeopleMD, error) {
	resultItems := make([]bag.IBagItem, 0)
	var heromd *bags.HeroMD
	peoples := make([]*bags.PeopleMD, 0)
	//check unlockedArea
	for _, v := range m.unlockedArea {
		if v == areaId {
			return nil, nil, nil, errorx.Building_UnLock_Err
		}
	}
	unlockTimes := m.getExploreTimes()
	areaCostCfg := configs.GetDb().GetAreaOpenCostCfg(unlockTimes)
	if areaCostCfg == nil {
		return nil, nil, nil, errorx.UnLock_Area_No_Err
	}
	//check clearbox
	if isUseGold {
		//check cost item
		err := bag.BagEx.CheckItem(user, areaCostCfg.Gold)
		if err != nil {
			return nil, nil, nil, err
		}
		//del item
		resultItems, err = bag.BagEx.DelItems(user, areaCostCfg.Gold, 0)
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		var exploreItem *common.ExploreItem
		var index int
		for k, v := range m.exploreItems {
			if v.Id == areaId {
				exploreItem = &v
				index = k
				break
			}
		}
		if exploreItem == nil {
			return nil, nil, nil, errorx.Explore_Area_No_Err
		}

		if exploreItem.EndTime > int(time.Now().Unix()) || exploreItem.EndTime == 0 {
			return nil, nil, nil, errorx.Explore_Area_Cd_Err
		}

		var err error

		//update hero
		if exploreItem.Hero != 0 {
			heromd, err = bags.GetHeroMDBycityID(user, cityId, exploreItem.Hero)
			if err != nil {
				return nil, nil, nil, err
			}
		}
		heromd.AddExp(areaCostCfg.PeopleNum)
		heromd.Status = bags.HeroStatus_Default

		//update peoloe
		bags.SetPeopleStatus(user.GetPeopleBag(), bags.HeroStatus_Default, exploreItem.WorkerIds...)

		peoples = user.GetPeopleBag().GetPeoplesByUID(exploreItem.WorkerIds)

		//update explore
		m.exploreItems = append(m.exploreItems[:index], m.exploreItems[index+1:]...)
	}

	//update  unlockedArea
	m.unlockedArea = append(m.unlockedArea, areaId)
	return resultItems, heromd, peoples, nil
}

func (m *MapObj) ClearRuins(user IMapUser, cityId, areaId int, isUseGold bool) ([]bag.IBagItem, *bags.HeroMD, []*bags.PeopleMD, error) {
	return nil, nil, nil, nil
}

func (m *MapObj) InitUnlockArea(user IMapUser, cityId, areaId int) error {
	//check unlockedArea
	for _, v := range m.unlockedArea {
		if v == areaId {
			return errorx.Building_UnLock_Err
		}
	}
	//update  unlockedArea
	m.unlockedArea = append(m.unlockedArea, areaId)
	m.setEmpityByUnlockedArea(areaId)
	return nil
}
