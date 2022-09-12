package configs

import "github.com/TTsmall/wawaji_pub_hmhelper/common"

//基础的全局配置表
type HelperConf struct {
	BuildingBagNumMax int             `conf:"buildingbagnummax" default:"10"` //建筑收起背包上限
	LoveItemEffect    int             `conf:"loveItemEffect" default:"1"`     //英雄使用喜好道具时经验的倍率
	AddSpeedCardTimes int             `conf:"addSpeedCardTimes" default:"3"`  //加速卡默认减少15分钟
	CoinSpeedTimes    int             `conf:"coinSpeedTimes" default:"3"`     //一个金币可减少1分钟，向上取整
	VideoReduceTimes  int             `conf:"videoReduceTimes" default:"3"`   //看一个视频广告加速15分钟
	FirstArea         common.IntSlice `conf:"firstArea" default:"3"`          //建角后的初始区域
}

// //返回按时间天数拿到的最大每日金币获得上限
// func (mg *HelperConf) GetMaxGoldByDayNum(daynum int) (result int) {
// 	for index := range GetDb().GameEx.MaxGolds {
// 		if daynum >= GetDb().GameEx.MaxGolds[index].K {
// 			result = GetDb().GameEx.MaxGolds[index].V
// 		}
// 	}
// 	return
// }
