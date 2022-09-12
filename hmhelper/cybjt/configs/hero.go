package configs

import "wawaji_pub/hmhelper/common"

//英雄有关

//英雄
type HeroCfg struct {
	Id             int              `col:"heroId" client:"heroId"`                 //ID
	Name           string           `col:"name" client:"name"`                     //名字
	Head           int              `col:"head" client:"head"`                     //头像
	DrawSpine      string           `col:"picture" client:"picture"`               //角色立绘
	Model          int              `col:"model" client:"model"`                   //走路模型
	Scale          float64          `col:"scale" client:"scale"`                   //缩放比
	Quality        int              `col:"type" client:"type"`                     //品质
	MaxLevel       int              `col:"maxLevel" client:"maxLevel"`             //等级上限
	MinStar        int              `col:"minStar" client:"minStar"`               //初始星级
	MaxStar        int              `col:"maxStar" client:"maxStar"`               //最大星级
	RecoveryReward common.ItemInfos `col:"recoveryReward" client:"recoveryReward"` //抽卡自动分解碎片
	HeroText       string           `col:"heroText" client:"heroText"`             //文本描述
	Prop1          string           `col:"prop1" client:"prop1"`                   //建造评级
	Prop2          string           `col:"prop2" client:"prop2"`                   //农牧评级
	Prop3          string           `col:"prop3" client:"prop3"`                   //制作评级
	Prop4          string           `col:"prop4" client:"prop4"`                   //理财评级
	Prop5          string           `col:"prop5" client:"prop5"`                   //奇遇评级
	LoveItem       common.IntSlice  `col:"loveItem" client:"loveItem"`             //喜好道具
	HouseInfo      string           `col:"houseInfo" client:"houseInfo"`           //居住描述
}

//英雄星级
type HeroStarCfg struct {
	Id         int              `col:"id" client:"id"`                 //双主键并出来的ID，HeroID*10+starID
	HeroID     int              `col:"heroId" client:"heroId"`         //英雄ID
	Name       string           `col:"name" client:"name"`             //英雄名字
	StarID     int              `col:"star" client:"star"`             //星级
	StarUpCost common.ItemInfos `col:"starUpCost" client:"starUpCost"` //升星消耗
	GeniusName string           `col:"geniusName" client:"geniusName"` //天赋名字
	GeniusInfo string           `col:"info" client:"info"`             //天赋说明
	GeniusID   common.IntSlice  `col:"geniusId" client:"geniusId"`     //天赋ID
	//GeniusValue string           `col:"geniusValue" client:"geniusValue"` //天赋值
}
type WorkTypeEnum int

const (
	//任意工作
	WorkTypeEnum_Default WorkTypeEnum = iota
	//初始建造
	WorkTypeEnum_CreateBuilding
	//建筑升级
	WorkTypeEnum_UpBuilding
	//生产工作
	WorkTypeEnum_Working
	//奇遇
	WorkTypeEnum_ExplorePve
	//探索新区域
	WorkTypeEnum_ExploreNew
	//修复
	WorkTypeEnum_RepairBuilding
)

//天赋
type GeniusCfg struct {
	Id            int          `col:"id" client:"id"`                         //天赋ID
	GeniusExplain string       `col:"geniusExplain" client:"geniusExplain"`   //天赋说明
	WorkType      WorkTypeEnum `col:"workType" client:"workType"`             //工作类型
	WorkCondID    int          `col:"effectiveType" client:"effectiveType"`   //天赋条件类型
	WorkCond      string       `col:"effectiveValue" client:"effectiveValue"` //天赋条件参数
	EffectID      int          `col:"geniusType" client:"geniusType"`         //天赋效果ID
	EffectValue   string       `col:"value" client:"value"`                   //效果值

}

//英雄等级
type HeroLevelProp struct {
	Id     int              `col:"id" client:"id"`         //双主键，HeroID*100+Level
	HeroID int              `col:"heroId" client:"heroId"` //英雄ID
	Level  int              `col:"level" client:"level"`   //等级
	Exp    int              `col:"exp" client:"exp"`       //升级需要的经验
	Prop   common.PropInfos `col:"prop" client:"prop"`     //等级属性
}

//品质信息
type HeroMaxCfg struct {
	QualiityID int `col:"quality" client:"quality"`     //品质
	MaxLevel   int `col:"maxLevel" client:"maxLevel"`   //最高等级限制
	Condition  int `col:"condition" client:"condition"` //条件？
}

//英雄的属性
type HeroPropertyCfg struct {
	Id          int    `col:"id" client:"id"`                   //属性ID
	DisplayName string `col:"displayName" client:"displayName"` //属性名字
	DisplayDes  string `col:"displayDes" client:"displayDes"`   //属性描述
	Type        int    `col:"type" client:"type"`               //类型(1:百分比,2固定值)
	ValType     int    `col:"valType" client:"valType"`         //数据类型(1:浮点数,2整数)
}

//英雄属性对应的天赋的折算
type HeroPropGeniusCfg struct {
	Id     int `col:"id" client:"id"`         //属性ID
	ProNum int `col:"proNum" client:"proNum"` //每多少点折算百分比
}

//----------------有关方法

//英雄信息
func (m *GameDb) GetHeroInfo(hid int) *HeroCfg {
	return m.HeroCfgs[hid]
}

//拿英雄星级
func (m *GameDb) GetHeroStarInfo(hid, star int) *HeroStarCfg {
	return m.HeroStarCfgs[hid*100+star]
}

//英雄天赋信息
func (m *GameDb) GetHeroGeniusInfo(gid int) *GeniusCfg {
	return m.GeniusCfgs[gid]
}

//英雄等级信息
func (m *GameDb) GetHeroLvInfo(hid, lv int) *HeroLevelProp {
	return m.HeroLevelProps[hid*100+lv]
}

// //英雄品质信息
// func (m *GameDb) GetHeroQuality(qid int) *HeroMaxCfg {
// 	return m.HeroMaxCfgs[qid]
// }
