package configs

import "github.com/TTsmall/wawaji_pub_hmhelper/common"

//NPC有关

//NPC，平民表
type NpcCfg struct {
	Id           int             `col:"id" client:"id"`     //ID
	Name         string          `col:"name" client:"name"` //名字
	Note         string          `col:"note" client:"note"`
	Icon         string          `col:"icon" client:"icon"`
	DisplayID    int             `col:"displayID" client:"displayID"`
	Type         int             `col:"type" client:"type"`
	LinkTo       common.IntSlice `col:"linkTo" client:"linkTo"`
	Plot         string          `col:"plot" client:"plot"`
	Talk         common.IntSlice `col:"talk" client:"talk"`
	TalkInterval string          `col:"talkInterval" client:"talkInterval"`
	Title        int             `col:"title" client:"title"`
	SecretNpc    common.IntSlice `col:"secretNpc" client:"secretNpc"`
	HouseInfo    string          `col:"houseInfo" client:"houseInfo"`
}

//----------------有关方法
