package building

import (
	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
	"github.com/TTsmall/wawaji_pub_hmhelper/common"
)

const (
	DECORATE_BUILDING   = 1 //装饰建筑
	AMUSEMENT_BUILDING  = 2 //娱乐建筑
	PRODUCTION_BUILDING = 3 //生产建筑
	OCCUP_BUILDING      = 4 //居住建筑

	BUILD_ITEMTYPE_MINJU = 52 //
	BUILD_ITEMTYPE_CAP   = 54

	SPEEDUPCOSTTYPE_WATCH_VIDEO = 1 //看视频加速
	SPEEDUPCOSTTYPE_GOLD        = 2 //金币加速
	SPEEDUPCOSTTYPE_CARD        = 3 //加速卡加速

	RELEASETYPE_BUILDING = 1
	RELEASETYPE_PRODUCE  = 2
	RELEASETYPE_LEVELUP  = 3
	RELEASETYPE_EXPLORE  = 4
)

var (
	BUILDING_TYPE = [4]int{DECORATE_BUILDING, AMUSEMENT_BUILDING, PRODUCTION_BUILDING, OCCUP_BUILDING}
)

type IUser interface {
	//用户对像需要实现的接口，使用KEY，拿到对应的仓库
	bag.IUser
	//GetBuilding(tp int) (result IBuilding)
	GetBuildIncrId() int
}

//move(forward),clear,package,(lvlup,award)st art,end;
type IBuilding interface {
	//开始建造
	StartBuild(user IUser, id, heroId, useBuildingTime int, peopleIds []int, pos common.PosInfo) (*BuildingInfo, error)
	//建造完成
	EndBuild(user IUser, uid int) (*BuildingInfo, error)
	//检查上限
	CheckLimit(user IUser, id int) bool
	//从建筑背包放下
	//
	GetPos(uid int) common.PosInfo
	//
	GetAllPos() []*BuildingPos

	SetBuildingInfo(building *BuildingInfo)

	ReleasePeople(uid, releaseType int) (*BuildingInfo, int, []int, error)

	DelteBuild(user IUser, uid int) (int, error)
}

type BuildingPos struct {
	Id  int
	Pos common.PosInfo
}

//基础建筑(装饰类)
type BuildingInfo struct {
	Id  int            //建筑表id
	Uid int            //唯一ID
	Pos common.PosInfo //位置
	//IsBuild              bool            `json:"isbuild,omitempty"`   //是否在建建筑   todo
	//UpEndTime            int             `json:"upendtime,omitempty"` //todo
	Builder              int              `json:"builders,omitempty"` //建筑用人id
	AwardEndTime         int              `json:"awardendtime,omitempty"`
	Occupiers            common.IntSlice  `json:"occupies,omitempty"`             //占用人
	BuilderPeople        common.IntSlice  `json:"builderPeople,omitempty"`        //工作中的工人
	BuildCompletedTime   int              `json:"buildCompletedTime,omitempty"`   //建造完成事件
	ProduceCompletedTime int              `json:"produceCompletedTime,omitempty"` //生产完成时间
	UpLevelCompletedTime int              `json:"upLevelCompletedTime,omitempty"` //升级完成时间
	Awards               common.ItemInfos `json:"awards,omitempty"`               //奖励
}

func (md *BuildingInfo) SetBuilder(hid int) {
	// fmt.Println("SetBuilder:", hid)
	// fmt.Println(string(debug.Stack()))
	md.Builder = hid
}

type Building struct {
	Buildings map[int]*BuildingInfo
}

func NewBuildingInfo(id, uid int, pos common.PosInfo, buildCompletedTime, builder int, builderPeople common.IntSlice) *BuildingInfo {
	return &BuildingInfo{
		Id:                 id,
		Uid:                uid,
		Pos:                pos,
		BuildCompletedTime: buildCompletedTime,
		Builder:            builder,
		BuilderPeople:      builderPeople,
		Occupiers:          make(common.IntSlice, 5),
	}
}

func NewBuilding() *Building {
	return &Building{
		Buildings: make(map[int]*BuildingInfo),
	}
}

//生产类
type IAwardBuilding interface {
	IBuilding
	LvlUpStartBuild(uid int, newPos common.PosInfo) error
	LvlUpEndBuild(uid int) error
	AwardStart(uid int) error
	AwardEnd(uid int) (resultItem *bag.ResultItem, err error)
}

type HeroSites []*HeroSite
type HeroSite struct {
	HeroId int
	Site   int
}

//居住类
type IOccupBuilding interface {
	IAwardBuilding
	OccupBuilding(uid, heroId int) error
}
