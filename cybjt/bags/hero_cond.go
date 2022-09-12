package bags

import (
	"github.com/TTsmall/wawaji_pub_hmhelper/cybjt/configs"
	"github.com/buguang01/util"
)

type WorkCondEnum int

const (
	//无条件
	WorkCondEnum_Default WorkCondEnum = iota + 1
	//TitleType
	WorkCondEnum_Title
	//WorkType
	WorkCondEnum_Work
	//ItemType
	WorkCondEnum_Item
)

func NewCond(cfg *configs.GeniusCfg) ICondition {
	switch WorkCondEnum(cfg.WorkCondID) {
	case WorkCondEnum_Default:
		return new(Condition)
	case WorkCondEnum_Title:
		result := new(CondByBuildTitleType)
		result.val = util.NewString(cfg.WorkCond).ToIntV()
		return result
	case WorkCondEnum_Work:
		result := new(CondBYBuildWorkType)
		result.val = util.NewString(cfg.WorkCond).ToIntV()
		return result
	case WorkCondEnum_Item:
		result := new(CondBYBuildItemType)
		result.val = util.NewString(cfg.WorkCond).ToIntV()
		return result
	}
	return new(Condition)
}

//条件
type ICondition interface {
	IsRun(args *MakeResult) bool
}

//type 1
type Condition struct {
}

func (cond *Condition) IsRun(args *MakeResult) bool {
	return true
}

//type 2
type CondByBuildTitleType struct {
	val int //类型
}

func (cond *CondByBuildTitleType) IsRun(args *MakeResult) bool {
	if args.BuildCfg == nil {
		return false
	}
	return cond.val == args.BuildCfg.TitleType
}

//type 3
type CondBYBuildWorkType struct {
	val int //类型

}

func (cond *CondBYBuildWorkType) IsRun(args *MakeResult) bool {
	if args.BuildCfg == nil {
		return false
	}
	return cond.val == args.BuildCfg.WorkType
}

//type 4
type CondBYBuildItemType struct {
	val int //类型

}

func (cond *CondBYBuildItemType) IsRun(args *MakeResult) bool {
	if args.BuildCfg == nil {
		return false
	}
	return cond.val == args.BuildCfg.Itemtype
}
