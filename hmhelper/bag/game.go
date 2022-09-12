package bag

import "wawaji_pub/hmhelper/common"

//基础的全局配置表
type HelperConf struct {
	MaxGolds common.PropInfos `conf:"maxGold" default:"1,3000;2,1100;3,1100"` //每日最大金币数，不填则使用之前的数据，类似封印等级的设计次

}

//返回按时间天数拿到的最大每日金币获得上限
func (mg *HelperConf) GetMaxGoldByDayNum(daynum int) (result int) {
	for index := range BagEx.GameEx.MaxGolds {
		if daynum >= BagEx.GameEx.MaxGolds[index].K {
			result = BagEx.GameEx.MaxGolds[index].V
		}
	}
	return
}
