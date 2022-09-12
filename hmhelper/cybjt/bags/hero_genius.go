package bags

import (
	"fmt"

	"github.com/buguang01/util"
	"wawaji_pub/hmhelper/common"
	"wawaji_pub/hmhelper/cybjt/configs"
)

//天赋

func FactoryGenius(ids common.IntSlice) (result []Igenius) {
	result = make([]Igenius, 0, 10)
	for k := range ids {
		cfg := configs.GetDb().GetHeroGeniusInfo(ids[k])
		if item := NewGenius(cfg); item != nil {
			result = append(result, item)
		} else {
			fmt.Println("FactoryGenius is null,by id:", k)
		}
	}
	return result
}

type GEniusEnum int

const (
	//减少时间（秒）百分比
	GEniusEnum_Time GEniusEnum = iota + 1
	//改变NPC数量
	GEniusEnum_NpcNum
	//减少消耗物品，百分比
	GEniusEnum_Items
	//额外获得物品
	GEniusEnum_Reward
	//增加获取产出，百分比
	GEniusEnum_Resli
)

func NewGenius(cfg *configs.GeniusCfg) Igenius {
	switch GEniusEnum(cfg.EffectID) {
	case GEniusEnum_Time:
		result := new(GeniusByTime)
		result.WorkType = cfg.WorkType
		result.Cond = NewCond(cfg)
		result.Val = util.NewString(cfg.EffectValue).ToIntV()
		return result
	case GEniusEnum_NpcNum:
		result := new(GeniusByNpcNum)
		result.WorkType = cfg.WorkType
		result.Cond = NewCond(cfg)
		result.Val = util.NewString(cfg.EffectValue).ToIntV()
		return result
	case GEniusEnum_Items:
		result := new(GeniusByItems)
		result.WorkType = cfg.WorkType
		result.Cond = NewCond(cfg)
		result.Items = common.NewBaseDataString(cfg.EffectValue)
		return result
	case GEniusEnum_Reward:
		result := new(GeniusByReward)
		result.WorkType = cfg.WorkType
		result.Cond = NewCond(cfg)
		result.Items = common.NewBaseDataString(cfg.EffectValue)
		return result
	case GEniusEnum_Resli:
		result := new(GeniusByResli)
		result.WorkType = cfg.WorkType
		result.Cond = NewCond(cfg)
		result.Items = common.NewBaseDataString(cfg.EffectValue)
		return result
	}
	fmt.Println("NewGenius nil.")
	return nil
}

func NewGeniusByHero(tp GEniusEnum, wktp configs.WorkTypeEnum, val int) Igenius {
	switch GEniusEnum(tp) {
	case GEniusEnum_Time:
		result := new(GeniusByTime)
		result.WorkType = wktp
		result.Cond = new(Condition)
		result.Val = val
		return result
	case GEniusEnum_NpcNum:
		result := new(GeniusByNpcNum)
		result.WorkType = wktp
		result.Cond = new(Condition)
		result.Val = val
		return result
	case GEniusEnum_Items:
		result := new(GeniusByItems)
		result.WorkType = wktp
		result.Cond = new(Condition)
		result.Items = common.NewBaseDataString("")
		result.Items.UpData(0, val)
		return result
	case GEniusEnum_Reward:
		result := new(GeniusByReward)
		result.WorkType = wktp
		result.Cond = new(Condition)
		result.Items = common.NewBaseDataString("")
		result.Items.UpData(0, val)
		return result
	case GEniusEnum_Resli:
		result := new(GeniusByResli)
		result.WorkType = wktp
		result.Cond = new(Condition)
		result.Items = common.NewBaseDataString("")
		result.Items.UpData(0, val)
		return result
	}
	fmt.Println("NewGeniusByHero nil.")
	return nil
}

//天赋接口
type Igenius interface {
	//运行天赋
	RunGenius(args *MakeResult)
}

type GeniusBase struct {
	Cond     ICondition           //条件
	WorkType configs.WorkTypeEnum //工作类型
}
type GeniusByTime struct {
	GeniusBase
	Val int //减少的百分比
}

//运行天赋
func (md *GeniusByTime) RunGenius(args *MakeResult) {
	if md.WorkType != args.WorkType && md.WorkType != 0 {
		return
	}
	if !md.Cond.IsRun(args) {
		return
	}
	args.BaseUsest += md.Val

}

type GeniusByNpcNum struct {
	GeniusBase
	Val int //变化的人数
}

//运行天赋
func (md *GeniusByNpcNum) RunGenius(args *MakeResult) {
	if md.WorkType != args.WorkType && md.WorkType != 0 {
		return
	}
	if !md.Cond.IsRun(args) {
		return
	}
	args.Npcnum += md.Val

}

type GeniusByItems struct {
	GeniusBase
	Items *common.BaseData
}

//运行天赋
func (md *GeniusByItems) RunGenius(args *MakeResult) {
	if md.WorkType != args.WorkType && md.WorkType != 0 {
		return
	}
	if !md.Cond.IsRun(args) {
		return
	}
	args.BaseItems.UpDataBc(md.Items, nil)

}

type GeniusByReward struct {
	GeniusBase
	Items *common.BaseData
}

//运行天赋
func (md *GeniusByReward) RunGenius(args *MakeResult) {
	if md.WorkType != args.WorkType && md.WorkType != 0 {
		return
	}
	if !md.Cond.IsRun(args) {
		return
	}
	args.Reward.UpDataBc(md.Items, nil)
}

type GeniusByResli struct {
	GeniusBase
	Items *common.BaseData
}

//运行天赋
func (md *GeniusByResli) RunGenius(args *MakeResult) {
	if md.WorkType != args.WorkType && md.WorkType != 0 {
		return
	}
	if !md.Cond.IsRun(args) {
		return
	}
	args.BaseResLi.UpDataBc(md.Items, nil)
}
