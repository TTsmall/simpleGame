package bags

import (
	"time"

	"wawaji_pub/hmhelper/common"
	"wawaji_pub/hmhelper/cybjt/configs"
)

//英雄的操作

type HeroStatusEnum int

const (
	//默认状态，空闲
	HeroStatus_Default HeroStatusEnum = iota
	//工作状态
	HeroStatus_Working
)

//百景图-英雄模块
type HeroMD struct {
	HeroID       int            `json:"uid"`              //英雄ID
	CityID       int            `json:"ctid"`             //城市ID
	IsHome       bool           `json:"ishome,omitempty"` //是否有家
	EquipUID     int            `json:"eid,omitempty"`    //装备UID
	HeroLv       int            `json:"hlv"`              //英雄等级
	HeroExp      int            `json:"hexp,omitempty"`   //英雄经验
	StarID       int            `json:"sid"`              //星级
	Status       HeroStatusEnum `json:"stid,omitempty"`   //状态
	Strength     int            `json:"sth"`              //体力
	StrengthTime time.Time      `json:"sthtime"`          //最后计算体力时间
}

//背包接口

//拿到物品ID
func (md *HeroMD) GetItemID() int {
	return md.HeroID
}

//拿实例UID,因为英雄是每个都是单例，所以UID=HeroID
func (md *HeroMD) GetUid() int {
	return md.HeroID
}

//操作接口

//添加经验,是否升级
func (md *HeroMD) AddExp(exp int) (result bool) {
	result = false
	md.HeroExp += exp
	hcfg := configs.GetDb().GetHeroInfo(md.HeroID)
	if hcfg.MaxLevel <= md.HeroLv {
		md.HeroLv = hcfg.MaxLevel
		md.HeroExp = 0
		return
	}
	lvcfg := configs.GetDb().GetHeroLvInfo(md.HeroID, md.HeroLv)
	for lvcfg.Exp <= md.HeroExp {
		result = true
		md.HeroLv++
		md.HeroExp -= lvcfg.Exp
		if hcfg.MaxLevel <= md.HeroLv {
			md.HeroLv = hcfg.MaxLevel
			md.HeroExp = 0
			return
		}
		lvcfg = configs.GetDb().GetHeroLvInfo(md.HeroID, md.HeroLv)

	}
	return
}

//是否满级
func (md *HeroMD) IsMaxLv() bool {
	hcfg := configs.GetDb().GetHeroInfo(md.HeroID)
	return hcfg.MaxLevel <= md.HeroLv
}

//是否满星
func (md *HeroMD) IsMaxStar() bool {
	hcfg := configs.GetDb().GetHeroInfo(md.HeroID)
	return hcfg.MaxStar <= md.StarID
}

//是否喜欢这个道具
func (md *HeroMD) IsLikesItem(itemid int) bool {
	hcfg := configs.GetDb().GetHeroInfo(md.HeroID)
	for k := range hcfg.LoveItem {
		if hcfg.LoveItem[k] == itemid {
			return true
		}
	}
	return false
}

//升星
func (md *HeroMD) UpStar() int {
	hcfg := configs.GetDb().GetHeroInfo(md.HeroID)
	if md.StarID < hcfg.MaxStar {
		md.StarID++
	}
	return md.StarID
}

//修改城市ID
func (md *HeroMD) EditCityID(cid int) int {
	md.CityID = cid
	return md.CityID
}

//修改是否有家
func (md *HeroMD) EditIsHome(v bool) bool {
	md.IsHome = v
	return md.IsHome
}

//获取英雄天赋
// func (md *HeroMD) GetGeniusVal() (geniusinfo *configs.GeniusCfg, val string) {
// 	starcfg := configs.GetDb().GetHeroStarInfo(md.HeroID, md.StarID)
// 	return configs.GetDb().GetHeroGeniusInfo(starcfg.GeniusID), starcfg.GeniusValue
// }

//获取英雄属性
func (md *HeroMD) GetProp() *HeroProbModel {
	result := new(HeroProbModel)
	result.prob = make(map[int]int)
	lvcfg := configs.GetDb().GetHeroLvInfo(md.HeroID, md.HeroLv)
	for _, v := range lvcfg.Prop {
		result.prob[v.K] = v.V
	}
	return result
}

//获取英雄是否在指定城市可用
func (md *HeroMD) GetCanUse(cid int) bool {
	// return true //md.Status == HeroStatus_Default
	// return md.CityID == cid && md.IsHome && md.Status == HeroStatus_Default
	return md.IsHome && md.Status == HeroStatus_Default
}

//1初始建造
func (md *HeroMD) CreateBuilding(cfg *configs.BuildCfg) (result *MakeResult) {
	probmd := md.GetProp()
	val := probmd.GetBuild()
	val = val / configs.GetDb().HeroPropGeniusCfgs[1021].ProNum
	if val > 100 {
		val = 100
	}

	if cfg.NeedHero == 1 {
		result = NewMakeResult(nil, cfg, configs.WorkTypeEnum_CreateBuilding, cfg.PeopleNum, cfg.BuildTime,
			cfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	} else {
		result = NewMakeResult(nil, cfg, 1, cfg.PeopleNum-1, cfg.BuildTime,
			cfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	}
	starcfg := configs.GetDb().GetHeroStarInfo(md.HeroID, md.StarID)
	genli := FactoryGenius(starcfg.GeniusID)
	genli = append(genli, NewGeniusByHero(GEniusEnum_Time, configs.WorkTypeEnum_CreateBuilding, val))
	for index := range genli {
		genli[index].RunGenius(result)
	}
	return
}

//2建筑升级
func (md *HeroMD) UpBuilding(cfg *configs.BuildCfg) (result *MakeResult) {
	probmd := md.GetProp()
	val := probmd.GetBuild()
	val = val / configs.GetDb().HeroPropGeniusCfgs[1021].ProNum
	if val > 100 {
		val = 100
	}
	if cfg.NeedHero == 1 {
		result = NewMakeResult(nil, cfg, configs.WorkTypeEnum_UpBuilding, cfg.PeopleNum, cfg.BuildTime,
			cfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	} else {
		result = NewMakeResult(nil, cfg, configs.WorkTypeEnum_UpBuilding, cfg.PeopleNum-1, cfg.BuildTime,
			cfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	}
	starcfg := configs.GetDb().GetHeroStarInfo(md.HeroID, md.StarID)
	genli := FactoryGenius(starcfg.GeniusID)
	genli = append(genli, NewGeniusByHero(GEniusEnum_Time, configs.WorkTypeEnum_UpBuilding, val))
	for index := range genli {
		genli[index].RunGenius(result)
	}
	return
}

//3生产工作
func (md *HeroMD) Working(cfg *configs.BuildCfg) (result *MakeResult) {
	probmd := md.GetProp()
	val := 0 // probmd.GetBuild()
	/*
			不填=没有工作
		1=农牧
		2=制作
		3=理财
		4=收税
		5=产出英雄经验书
		6=民宅产出
	*/
	switch cfg.WorkType {
	case 1:
		val = probmd.GetAgriculture()
		val = val / configs.GetDb().HeroPropGeniusCfgs[1031].ProNum
	case 2:
		val = probmd.GetMaking()
		val = val / configs.GetDb().HeroPropGeniusCfgs[1041].ProNum
	case 3:
		val = probmd.GetFinancial()
		val = val / configs.GetDb().HeroPropGeniusCfgs[1051].ProNum
	}
	if val > 100 {
		val = 100
	}
	if cfg.NeedHero == 1 {
		result = NewMakeResult(nil, cfg, configs.WorkTypeEnum_Working, cfg.WorkPeopleNum, cfg.WorkTime,
			cfg.WorkCost, common.ItemInfos{}, cfg.WorkOut)
	} else {
		result = NewMakeResult(nil, cfg, configs.WorkTypeEnum_Working, cfg.WorkPeopleNum-1, cfg.WorkTime,
			cfg.WorkCost, common.ItemInfos{}, cfg.WorkOut)
	}
	starcfg := configs.GetDb().GetHeroStarInfo(md.HeroID, md.StarID)
	genli := FactoryGenius(starcfg.GeniusID)
	genli = append(genli, NewGeniusByHero(GEniusEnum_Resli, configs.WorkTypeEnum_Working, val))
	for index := range genli {
		genli[index].RunGenius(result)
	}
	return
}

//4奇遇

//5探索新区域
func (md *HeroMD) ExploreNew(cfg *configs.AreaOpenCostCfg) (result *MakeResult) {
	probmd := md.GetProp()
	val := probmd.GetBuild()
	val = val / configs.GetDb().HeroPropGeniusCfgs[1021].ProNum
	if val > 100 {
		val = 100
	}
	if cfg.NeedHero == 1 {
		result = NewMakeResult(cfg, nil, configs.WorkTypeEnum_ExploreNew, cfg.PeopleNum, cfg.UseTime,
			cfg.Cost, common.ItemInfos{}, common.ItemInfos{})
	} else {
		result = NewMakeResult(cfg, nil, configs.WorkTypeEnum_ExploreNew, cfg.PeopleNum-1, cfg.UseTime,
			cfg.Cost, common.ItemInfos{}, common.ItemInfos{})
	}
	starcfg := configs.GetDb().GetHeroStarInfo(md.HeroID, md.StarID)
	genli := FactoryGenius(starcfg.GeniusID)
	genli = append(genli, NewGeniusByHero(GEniusEnum_Time, configs.WorkTypeEnum_UpBuilding, val))
	for index := range genli {
		genli[index].RunGenius(result)
	}
	return
}

//6修复
func (md *HeroMD) RepairBuilding(cfg *configs.BuildCfg) (result *MakeResult) {
	probmd := md.GetProp()
	val := probmd.GetBuild()
	val = val / configs.GetDb().HeroPropGeniusCfgs[1021].ProNum
	if val > 100 {
		val = 100
	}
	if cfg.NeedHero == 1 {
		result = NewMakeResult(nil, cfg, configs.WorkTypeEnum_RepairBuilding, cfg.PeopleNum, cfg.BuildTime,
			cfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	} else {
		result = NewMakeResult(nil, cfg, configs.WorkTypeEnum_RepairBuilding, cfg.PeopleNum-1, cfg.BuildTime,
			cfg.BuildCost, common.ItemInfos{}, common.ItemInfos{})
	}
	starcfg := configs.GetDb().GetHeroStarInfo(md.HeroID, md.StarID)
	genli := FactoryGenius(starcfg.GeniusID)
	genli = append(genli, NewGeniusByHero(GEniusEnum_Time, configs.WorkTypeEnum_RepairBuilding, val))
	for index := range genli {
		genli[index].RunGenius(result)
	}
	return
}

//英雄属性
type HeroProbModel struct {
	prob map[int]int
}

//建造
func (md *HeroProbModel) GetBuild() int {
	val := 0
	val = (md.prob[1022]*(md.prob[1023]+100)/100 +
		md.prob[1024]*(md.prob[1025]+100)/100) *
		(md.prob[1026] + 100) / 100
	return val
}

//农牧
func (md *HeroProbModel) GetAgriculture() int {
	val := 0
	val = (md.prob[1032]*(md.prob[1033]+100)/100 +
		md.prob[1034]*(md.prob[1035]+100)/100) *
		(md.prob[1036] + 100) / 100
	return val
}

//制造
func (md *HeroProbModel) GetMaking() int {
	val := 0
	val = (md.prob[1042]*(md.prob[1043]+100)/100 +
		md.prob[1044]*(md.prob[1045]+100)/100) *
		(md.prob[1046] + 100) / 100
	return val
}

//理财
func (md *HeroProbModel) GetFinancial() int {
	val := 0
	val = (md.prob[1052]*(md.prob[1053]+100)/100 +
		md.prob[1054]*(md.prob[1055]+100)/100) *
		(md.prob[1056] + 100) / 100
	return val
}

//冒险
func (md *HeroProbModel) GetAdventure() int {
	val := 0
	val = (md.prob[1062]*(md.prob[1063]+100)/100 +
		md.prob[1064]*(md.prob[1065]+100)/100) *
		(md.prob[1066] + 100) / 100
	return val
}

//英雄工作返回属性
type MakeResult struct {
	AreaCfg  *configs.AreaOpenCostCfg //区域
	BuildCfg *configs.BuildCfg        //建筑
	WorkType configs.WorkTypeEnum     //工作类型

	BaseUsest int              //缩减时间比例，0
	BaseItems *common.BaseData //消耗的资源减少比例，0
	BaseResLi *common.BaseData //获得的资源增加比例，0
	Reward    *common.BaseData //额外获得资源

	Npcnum int              //使用NPC的数量
	Usest  int              //使用的时间（秒）
	Items  common.ItemInfos //消耗的资源
	ResLi  common.ItemInfos //获得的资源

}

func NewMakeResult(area *configs.AreaOpenCostCfg, build *configs.BuildCfg,
	wktype configs.WorkTypeEnum, npcnum, usest int,
	items, reward, resli common.ItemInfos) *MakeResult {
	result := new(MakeResult)
	result.AreaCfg = area
	result.BuildCfg = build
	result.WorkType = wktype
	result.Npcnum = npcnum
	result.BaseUsest = 0
	result.Usest = usest
	result.BaseItems = common.NewBaseDataString("")
	result.Items = make(common.ItemInfos, 0)
	result.Items = append(result.Items, items...)
	result.Reward = common.NewBaseDataString("")

	result.BaseResLi = common.NewBaseDataString("")
	result.ResLi = make(common.ItemInfos, 0)
	result.ResLi = append(result.ResLi, resli...)
	return result
}

//需要的NPC数
func (md *MakeResult) GetNpcNum() int {
	if md.Npcnum < 0 {
		return 0
	}
	return md.Npcnum
}

//需要的时间（秒）
func (md *MakeResult) GetUseTime() int {
	if md.BaseUsest > 100 {
		md.BaseUsest = 100
	}
	return md.Usest * (100 - md.BaseUsest) / 100
}

//需要的物品
func (md *MakeResult) GetItems() (result common.ItemInfos) {
	result = make(common.ItemInfos, 0, 10)
	for _, v := range md.Items {
		item := &common.ItemInfo{
			ItemId: v.ItemId,
			Count:  v.Count,
		}
		item.Count = item.Count *
			(100 - md.BaseItems.GetNumByKey(item.ItemId) - md.BaseItems.GetNumByKey(0)) / 100
		if item.Count < 0 {
			item.Count = 0
		}
		result = append(result, item)
	}
	return result
}

//可能的额外获得的奖励
func (md *MakeResult) GetReward() (result common.ItemInfos) {
	result = make(common.ItemInfos, 0, 10)
	for k, v := range md.Reward.Data {
		item := &common.ItemInfo{
			ItemId: k,
			Count:  v,
		}
		result = append(result, item)
	}
	return result
}

//生产时的结果
func (md *MakeResult) GetResLi() (result common.ItemInfos) {
	result = make(common.ItemInfos, 0, 10)
	for _, v := range md.ResLi {
		item := &common.ItemInfo{
			ItemId: v.ItemId,
			Count:  v.Count,
		}
		item.Count = item.Count *
			(100 + md.BaseResLi.GetNumByKey(item.ItemId) + md.BaseResLi.GetNumByKey(0)) / 100
		if item.Count < 0 {
			item.Count = 0
		}
		result = append(result, item)
	}
	return result
}
