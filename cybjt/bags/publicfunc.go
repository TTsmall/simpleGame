package bags

import (
	"fmt"

	"github.com/TTsmall/wawaji_pub_hmhelper/bag"
	"github.com/TTsmall/wawaji_pub_hmhelper/errorx"
)

//设置NPC状态
func SetPeopleStatus(bagmd *PeopleBagMD, sts HeroStatusEnum, ids ...int) {
	for i := range ids {
		if md, ok := bagmd.GetPeopleByUID(ids[i]); ok {
			md.Status = sts
		}
	}
}

//拿英雄
func GetHeroMDBycityID(user bag.IUser, cityId, heroId int) (heromd *HeroMD, err error) {
	defer func() {
		fmt.Printf("GetHeroMD:%+v,%+v", heromd, err)
	}()
	heroItem, ok := user.GetBag(BagKey_Hero).(*bag.BagExMD).GetItemByUID(heroId)
	if !ok {
		return nil, errorx.Build_Hero_Not_Error
	}
	heromd = heroItem.(*HeroMD)
	if !heromd.GetCanUse(cityId) {
		return nil, errorx.Build_Hero_No_Use_Error
	}
	return
}

func GetHeroMD(user bag.IUser, heroId int) (heromd *HeroMD, err error) {
	heroItem, ok := user.GetBag(BagKey_Hero).(*bag.BagExMD).GetItemByUID(heroId)
	if !ok {
		return nil, errorx.Build_Hero_Not_Error
	}
	heromd = heroItem.(*HeroMD)
	return heromd, nil
}
