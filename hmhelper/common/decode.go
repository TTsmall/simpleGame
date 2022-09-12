package common

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	SEMICOLON = ";"
	COMMA     = ","
	COLON     = ":"
	PIPE      = "|"
	SPACE     = " "
	HLINE     = "-"
)

type ItemInfoPb struct {
	ItemId int `client:"key"`
	Count  int `client:"value"`
	Delta  int `client:"delta"`
}

type ItemInfo struct {
	ItemId int `client:"key"`
	Count  int `client:"value"`
}

type FloatItemInfo struct {
	ItemId int     `client:"key"`
	Count  float64 `client:"value"`
}

type FloatItemProbInfo struct {
	ItemId int     `client:"itemId"`
	Count  float64 `client:"count"`
	Prob   int     `client:"prob"`
}

type ItemInfoProb struct {
	ItemId int `client:"itemId"`
	Count  int `client:"count"`
	Prob   int `client:"prob"`
}

type PropInfo struct {
	K int `client:"key"`
	V int `client:"value"`
}

// 属性增益，兼容浮点和整数
type PropGain struct {
	Type  int     `client:"type"`
	Float float64 `client:"value"`
	Int   int
}

type WaveItem struct {
	Wave  int `client:"wave"`  // 波数
	Id    int `client:"id"`    // 物品ID
	Count int `client:"count"` // 数量
}

type Condition struct {
	K    int         `client:"key"`
	V    int         `client:"value"`
	Subs map[int]int `client:"subs"`
}

type DateTime struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

type ItemInfoProbs []*ItemInfoProb
type Conditions []*Condition
type ItemInfos []*ItemInfo
type FloatItemInfos []*FloatItemInfo
type FloatItemProbInfos []*FloatItemProbInfo

type PropInfos []*PropInfo
type PropGains []*PropGain
type ItemInfoPbs []*ItemInfoPb

// type IntSlice []int
// type FloatSlice []float64

// func (s IntSlice) Len() int           { return len(s) }
// func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
// func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type WaveItems []*WaveItem

type IntMap map[int]int

func (this *IntSlice) Decode(str string) error {
	ints, err := IntSliceFromString(str, ",")
	if err != nil {
		return err
	}
	*this = IntSlice(ints)
	return nil
}

func (this IntSlice) ToInt32Slice() []int32 {
	l := len(this)
	ret := make([]int32, l)
	if l == 0 {
		return ret
	}
	for i := 0; i < l; i++ {
		ret[i] = int32(this[i])
	}
	return ret
}

func (this *FloatSlice) Decode(str string) error {
	floats, err := FloatSliceFromString(str, ",")
	if err != nil {
		return err
	}
	*this = FloatSlice(floats)
	return nil
}

// func (this IntSlice) GetOrLast(index int) int {
// 	l := len(this)
// 	if index < l {
// 		return this[index]
// 	}
// 	if l == 0 {
// 		return 0
// 	}
// 	return this[l-1]
// }

// func (this *StringSlice) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		*this = make([]string, 0)
// 		return nil
// 	}
// 	*this = strings.Split(str, ";")
// 	return nil
// }

func (this ItemInfos) Times(times float64) ItemInfos {
	result := make(ItemInfos, 0, len(this))
	for _, itemInfo := range this {
		count := int(float64(itemInfo.Count) * times)
		if count > 0 {
			result = append(result, &ItemInfo{ItemId: itemInfo.ItemId, Count: count})
		}
	}
	return result
}

func (this ItemInfos) ToInt32Map() map[int32]int32 {
	result := make(map[int32]int32, len(this))
	for _, itemInfo := range this {
		result[int32(itemInfo.ItemId)] += int32(itemInfo.Count)
	}
	return result
}

func (this ItemInfos) Get(itemId int) *ItemInfo {
	for _, itemInfo := range this {
		if itemInfo.ItemId == itemId {
			return itemInfo
		}
	}
	return nil
}

//GetOrLast 取index位置的或者最后一个。
// func (this ItemInfos) GetOrLast(index int) *ItemInfo {
// 	l := len(this)
// 	if l == 0 {
// 		return nil
// 	}
// 	if index < l {
// 		return this[index]
// 	}
// 	return this[l-1]
// }

func (this *ItemInfos) Decode(str string) error {
	*this = make(ItemInfos, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		list := strings.Split(strings.TrimSpace(v), COMMA)
		if len(list) < 2 {
			return errors.New(v + "物品信息格式错误")
		}
		itemId, _ := strconv.Atoi(list[0])
		var itemInfo ItemInfo
		itemInfo.ItemId = itemId
		itemCount, err := strconv.Atoi(list[1])
		if err != nil {
			return err
		}
		itemInfo.Count = itemCount
		*this = append(*this, &itemInfo)
	}
	return nil
}

func (this *FloatItemInfos) Decode(str string) error {
	*this = make(FloatItemInfos, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		list := strings.Split(strings.TrimSpace(v), COMMA)
		if len(list) < 2 {
			return errors.New(v + "物品信息格式错误")
		}
		itemId, _ := strconv.Atoi(list[0])
		var floatItemInfo FloatItemInfo
		floatItemInfo.ItemId = itemId
		itemCount, err := strconv.ParseFloat(list[1], 64)
		if err != nil {
			return err
		}
		floatItemInfo.Count = itemCount
		*this = append(*this, &floatItemInfo)
	}
	return nil
}

// func (this *ItemInfosSlice) Decode(str string) error {
// 	*this = make(ItemInfosSlice, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.TrimSpace(str), PIPE)
// 	if len(infoList) == 0 {
// 		return nil
// 	}
// 	for _, v := range infoList {
// 		itemInfos := &ItemInfos{}
// 		err := itemInfos.Decode(v)
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, *itemInfos)
// 	}
// 	return nil
// }

func (this *PropInfos) Decode(str string) error {
	*this = make(PropInfos, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		list := strings.Split(strings.TrimSpace(v), COMMA)
		if len(list) < 2 {
			return errors.New(v + "属性信息格式错误")
		}
		k, err := strconv.Atoi(list[0])
		if err != nil {
			return err
		}
		var propInfo PropInfo
		propInfo.K = k
		propInfo.V, err = strconv.Atoi(list[1])
		if err != nil {
			return err
		}
		*this = append(*this, &propInfo)
	}
	return nil
}

func (this PropInfos) Rand() *PropInfo {
	var sum = 0
	var weight = rand.Intn(100)
	for _, info := range this {
		sum += info.V
		if weight < sum {
			return info
		}
	}
	return this[rand.Intn(len(this))]
}

func (this *PropGains) Decode(str string) error {
	*this = make(PropGains, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		gain := &PropGain{}
		err := gain.Decode(v)
		if err != nil {
			return err
		}
		*this = append(*this, gain)
	}
	return nil
}

func (this *IntMap) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}
	*this = make(IntMap)
	for _, v := range infoList {
		list := strings.Split(strings.TrimSpace(v), COMMA)
		if len(list) != 2 {
			return errors.New(v + "IntMap 属性信息格式错误")
		}

		k, err := strconv.Atoi(list[0])
		if err != nil {
			return err
		}
		if _, ok := (*this)[k]; ok {
			return errors.New(v + "IntMap 属性重复")
		}
		v, err := strconv.Atoi(list[1])
		if err != nil {
			return err
		}
		(*this)[k] = v
	}
	return nil

}

func (this *PropInfo) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(str, COMMA)
	if len(infoList) < 2 {
		return errors.New(str + " PropInfo 属性信息格式错误")
	}
	var propInfo PropInfo
	propInfo.K, _ = strconv.Atoi(infoList[0])
	propInfo.V, _ = strconv.Atoi(infoList[1])
	*this = propInfo
	return nil
}

func (this *PropGain) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(str, COMMA)
	if len(infoList) < 2 {
		return errors.New(str + " PropPlus 属性信息格式错误")
	}
	var err error
	this.Type, err = strconv.Atoi(infoList[0])
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(infoList[1], 64)
	if err != nil {
		return err
	}
	this.Float = float64(f)
	this.Int = int(f)
	return nil
}

func (this *ItemInfoProb) ToItemInfo() *ItemInfo {
	return &ItemInfo{ItemId: this.ItemId, Count: this.Count}
}

func (this *FloatItemProbInfo) ToItemInfo() *FloatItemProbInfo {
	return &FloatItemProbInfo{ItemId: this.ItemId, Count: this.Count, Prob: this.Prob}
}

func (this *ItemInfoProbs) Decode(str string) error {
	*this = make(ItemInfoProbs, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		list := strings.Split(strings.TrimSpace(v), COMMA)
		if len(list) != 3 {
			return errors.New(str + " ProbItems属性信息格式错误")
		}
		id, err := strconv.Atoi(list[0])
		if err != nil {
			return err
		}
		var probItem ItemInfoProb
		probItem.ItemId = id
		probItem.Count, err = strconv.Atoi(list[1])
		probItem.Prob, err = strconv.Atoi(list[2])
		if err != nil {
			return err
		}
		*this = append(*this, &probItem)
	}
	return nil
}

func (this ItemInfoProbs) Clone() ItemInfoProbs {
	results := make(ItemInfoProbs, 0, len(this))
	results = append(results, this...)
	return results
}

//这个转盘是每次抽走一个奖品就会重新计算weight
func (this ItemInfoProbs) RandUniqueViaWeight(count int) ItemInfos {
	var totalRate int
	for _, item := range this {
		totalRate += item.Prob
	}

	awardItemsList := this.Clone()

	results := make(ItemInfos, 0)
	for i := 1; i <= count; i++ {
		currentRandom := 0
		random := rand.Intn(totalRate)
		for index, item := range awardItemsList {
			currentRandom += item.Prob
			if currentRandom >= random {
				results = append(results, item.ToItemInfo())
				awardItemsList = append(awardItemsList[:index], awardItemsList[index+1:]...)
				totalRate = totalRate - item.Prob
				break
			}
		}
	}
	return results
}

//这个转盘是每次抽走一个奖品就会不重新计算weight
func (this ItemInfoProbs) RandUniqueViaAllWeight(count int) ItemInfos {
	if count == 0 {
		return nil
	}
	var totalRate int
	for _, item := range this {
		totalRate += item.Prob
	}
	awardItemsList := this.Clone()
	results := make(ItemInfos, 0)
	for i := 1; i <= count; i++ {
		currentRandom := 0
		random := rand.Intn(totalRate)
		for _, item := range awardItemsList {
			currentRandom += item.Prob
			if currentRandom >= random {
				results = append(results, item.ToItemInfo())
				break
			}
		}
	}
	return results
}

func (this *ItemInfoProb) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(str, COMMA)
	if len(infoList) < 3 {
		return errors.New(str + " 属性信息格式错误")
	}
	var itemInfo ItemInfoProb
	itemInfo.ItemId, _ = strconv.Atoi(infoList[0])
	itemInfo.Count, _ = strconv.Atoi(infoList[1])
	itemInfo.Prob, _ = strconv.Atoi(infoList[2])
	*this = itemInfo
	return nil
}

func (this *ItemInfo) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(str, COMMA)
	if len(infoList) < 2 {
		return errors.New(str + " 属性信息格式错误")
	}
	var itemInfo ItemInfo
	itemInfo.ItemId, _ = strconv.Atoi(infoList[0])
	itemInfo.Count, _ = strconv.Atoi(infoList[1])
	*this = itemInfo
	return nil
}

// func (this *ProbValue) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.TrimSpace(str), COMMA)
// 	if len(infoList) < 2 {
// 		return errors.New(str + " ProbValue 属性信息格式错误")
// 	}
// 	var probInfo ProbValue
// 	probInfo.Value, _ = strconv.Atoi(infoList[0])
// 	probInfo.Prob, _ = strconv.Atoi(infoList[1])
// 	*this = probInfo
// 	return nil
// }

// func (this *HmsTime) String() string {
// 	return fmt.Sprintf("%d:%d:%d", this.Hour, this.Minute, this.Second)
// }

// func (this HmsTime) Today() time.Time {
// 	now := time.Now()
// 	return time.Date(now.Year(), now.Month(), now.Day(), this.Hour, this.Minute, this.Second, 0, time.Local)
// }

// func (this *HmsTimes) Decode(str string) error {
// 	*this = make(HmsTimes, 0)
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.TrimSpace(str), SEMICOLON)
// 	if len(infoList) < 1 {
// 		return errors.New(str + " HmsTime 属性信息格式错误")
// 	}
// 	for _, v := range infoList {
// 		one := &HmsTime{}
// 		err := one.Decode(v)
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, one)
// 	}
// 	return nil
// }

// func (this *DateTime) String() string {
// 	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", this.Year, this.Month, this.Day, this.Hour, this.Minute, this.Second)
// }

// func (this *DateTime) ToTime() time.Time {
// 	return time.Date(this.Year, time.Month(this.Month), this.Day, this.Hour, this.Minute, this.Second, 0, time.Local)
// }

// func (this *DateTime) Unix() int64 {
// 	if this.Month == 0 {
// 		return 0
// 	}
// 	return this.ToTime().Unix()
// }

// func (this *DateTime) Unix32() int32 {
// 	return int32(this.Unix())
// }

// func (this *DateTime) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.TrimSpace(str), SPACE)
// 	if len(infoList) < 2 {
// 		return errors.New(str + " DateTime 属性信息格式错误")
// 	}

// 	infoListLeft := strings.Split(strings.TrimSpace(infoList[0]), HLINE)
// 	if len(infoListLeft) < 3 {
// 		return errors.New(str + " DateTime 属性格式错误")
// 	}
// 	this.Year, _ = strconv.Atoi(infoListLeft[0])
// 	this.Month, _ = strconv.Atoi(infoListLeft[1])
// 	this.Day, _ = strconv.Atoi(infoListLeft[2])

// 	infoListRight := strings.Split(strings.TrimSpace(infoList[1]), COLON)
// 	if len(infoListRight) < 3 {
// 		return errors.New(str + " DateTime 属性信息格式错误")
// 	}
// 	this.Hour, _ = strconv.Atoi(infoListRight[0])
// 	this.Minute, _ = strconv.Atoi(infoListRight[1])
// 	this.Second, _ = strconv.Atoi(infoListRight[2])

// 	return nil
// }
// func (this *HmsTime) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.TrimSpace(str), COLON)
// 	if len(infoList) < 1 {
// 		return errors.New(str + " HmsTime属性信息格式错误")
// 	}
// 	var hms HmsTime
// 	hms.Hour, _ = strconv.Atoi(infoList[0])
// 	if len(infoList) > 1 {
// 		hms.Minute, _ = strconv.Atoi(infoList[1])
// 	}
// 	if len(infoList) > 2 {
// 		hms.Second, _ = strconv.Atoi(infoList[2])
// 	}
// 	if hms.Hour < 0 || hms.Hour > 23 || hms.Minute < 0 || hms.Minute > 59 || hms.Second < 0 || hms.Second > 59 {
// 		return errors.New(str + "时分秒不对")
// 	}
// 	*this = hms
// 	return nil
// }

// func (this *SkillIdInfos) Decode(str string) error {
// 	*this = make(SkillIdInfos, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	skillIds := strings.Split(strings.TrimSpace(str), COMMA)
// 	if len(skillIds) == 0 {
// 		return nil
// 	}
// 	for _, v := range skillIds {
// 		if len(v) > 0 {
// 			// partnerT := gameDb.GetPartnerByName(v)
// 			// if partnerT == nil {
// 			// 	return errors.New("怪物名称输入错误")
// 			// }
// 			val, _ := strconv.Atoi(v)
// 			*this = append(*this, val)
// 		} else {
// 			*this = append(*this, 0)
// 		}
// 	}
// 	return nil
// }
// func (this *ProbValues) Decode(str string) error {
// 	*this = make(ProbValues, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		list := strings.Split(strings.TrimSpace(v), COMMA)
// 		if len(list) < 2 {
// 			continue
// 		}
// 		k, err := strconv.Atoi(list[0])
// 		if err != nil {
// 			return err
// 		}
// 		var probValue ProbValue
// 		probValue.Value = k
// 		probValue.Prob, err = strconv.Atoi(list[1])
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, &probValue)
// 	}
// 	return nil
// }

// func (this *ItemIndexInfos) Decode(str string) error {
// 	*this = make(ItemIndexInfos, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	itemIndexes := strings.Split(strings.TrimSpace(str), COMMA)
// 	if len(itemIndexes) == 0 {
// 		return nil
// 	}
// 	for _, v := range itemIndexes {
// 		if len(v) > 0 {
// 			// partnerT := gameDb.GetPartnerByName(v)
// 			// if partnerT == nil {
// 			// 	return errors.New("怪物名称输入错误")
// 			// }
// 			val, _ := strconv.Atoi(v)
// 			*this = append(*this, val)
// 		} else {
// 			*this = append(*this, 0)
// 		}
// 	}
// 	return nil
// }
// func (this *ProbItems) Decode(str string) error {
// 	*this = make(ProbItems, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		list := strings.Split(strings.TrimSpace(v), COMMA)
// 		if len(list) != 3 {
// 			return errors.New(str + " ProbItems属性信息格式错误")
// 		}
// 		id, err := strconv.Atoi(list[0])
// 		if err != nil {
// 			return err
// 		}
// 		var probItem ProbItem
// 		probItem.Id = id
// 		probItem.Count, err = strconv.Atoi(list[1])
// 		probItem.Prob, err = strconv.Atoi(list[2])
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, &probItem)
// 	}
// 	return nil
// }

// func (this *DisPlayIds) Decode(str string) error {
// 	*this = make(DisPlayIds, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		list := strings.Split(strings.TrimSpace(v), COMMA)
// 		if len(list) < 3 {
// 			//return errors.New(v + "属性信息格式错误")
// 			continue
// 		}
// 		id, err := strconv.Atoi(list[0])
// 		if err != nil {
// 			return err
// 		}
// 		var disPlayId DisPlayId
// 		disPlayId.WeaponId = id
// 		disPlayId.ClothesId, err = strconv.Atoi(list[1])
// 		disPlayId.WingId, err = strconv.Atoi(list[2])
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, &disPlayId)
// 	}
// 	return nil
// }

// func (this *RangeNums) Decode(str string) error {
// 	*this = make(RangeNums, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(str, SEMICOLON)
// 	for _, v := range infoList {
// 		ranges := strings.Split(v, COLON)
// 		if len(ranges) < 2 {
// 			return errors.New(str + " RangeNums 属性信息格式错误")
// 		}
// 		var rangeNum RangeNum
// 		rangeNum.Min, _ = strconv.Atoi(ranges[0])
// 		rangeNum.Max, _ = strconv.Atoi(ranges[1])
// 		*this = append(*this, &rangeNum)
// 	}
// 	return nil
// }

// func (this *SignItem) Decode(str string) error {
// 	infoList := strings.Split(str, COMMA)
// 	if len(infoList) < 2 {
// 		return errors.New(str + " SignItem 属性信息格式错误")
// 	}
// 	var signItem SignItem
// 	signItem.Id, _ = strconv.Atoi(infoList[0])
// 	signItem.Count, _ = strconv.Atoi(infoList[1])
// 	*this = signItem
// 	return nil
// }

// func (this *RangeNum) Decode(str string) error {
// 	infoList := strings.Split(str, COLON)
// 	if len(infoList) < 2 {
// 		return errors.New(str + " Rangenum属性信息格式错误")
// 	}
// 	var rangeNum RangeNum
// 	rangeNum.Min, _ = strconv.Atoi(infoList[0])
// 	rangeNum.Max, _ = strconv.Atoi(infoList[1])
// 	*this = rangeNum
// 	return nil
// }

// func (this *RangeNum) Rand() int {
// 	return common.RandNum(this.Min, this.Max)
// }

// // 复用ProbItems的Decode
// func (this *ProbDiamonds) Decode(str string) error {
// 	*this = make(ProbDiamonds, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		list := strings.Split(strings.TrimSpace(v), COMMA)
// 		if len(list) < 3 {
// 			// return errors.New(v + "属性信息格式错误")
// 			continue
// 		}
// 		var err error
// 		var values [3]int
// 		for i := 0; i < 3; i++ {
// 			values[i], err = strconv.Atoi(list[i])
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		var probDiamond ProbDiamond
// 		probDiamond.MinDiamond = values[0]
// 		probDiamond.MaxDiamond = values[1]
// 		probDiamond.Prob = values[2]
// 		*this = append(*this, &probDiamond)
// 	}
// 	return nil
// }

// func (this *IntSlice2) Decode(str string) error {
// 	*this = make(IntSlice2, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		ints, err := common.IntSliceFromString(v, ",")
// 		if err != nil {
// 			return fmt.Errorf("IntSlice2,Decode,err = %s", err.Error())
// 		}
// 		*this = append(*this, ints)
// 	}
// 	return nil
// }

// func (this *IntMap) Add(delta IntMap) {
// 	for k, v := range delta {
// 		(*this)[k] += v
// 	}
// }

// //Times,只保留>0的数据。
// func (this IntMap) Times(times float64) IntMap {
// 	newMap := make(IntMap, len(this))
// 	for k, v := range this {
// 		newValue := int(float64(v) * times)
// 		if newValue > 0 {
// 			newMap[k] = newValue
// 		}
// 	}
// 	return newMap
// }

// func (this IntMap) ToItemInfos() ItemInfos {
// 	itemInfos := make(ItemInfos, len(this))
// 	i := 0
// 	for k, v := range this {
// 		itemInfos[i] = &ItemInfo{ItemId: k, Count: v}
// 		i++
// 	}
// 	return itemInfos
// }

// func (this IntMap) Clone() IntMap {
// 	ret := make(IntMap, len(this))
// 	for k, v := range this {
// 		ret[k] = v
// 	}
// 	return ret
// }

func (this *WaveItem) Decode(str string) error {

	values, err := IntSliceFromString(strings.TrimSpace(str), COMMA)
	if err != nil || len(values) < 3 {
		return fmt.Errorf("WaveItem:Decode bad str:%s,err:%v", str, err)
	}
	var waveItem WaveItem
	waveItem.Wave = values[0]
	waveItem.Id = values[1]
	waveItem.Count = values[2]
	*this = waveItem
	return nil
}

// 复用ProbItems的Decode
func (this *WaveItems) Decode(str string) error {
	*this = make(WaveItems, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		var waveItem WaveItem
		err := waveItem.Decode(v)
		if err != nil {
			return err
		}
		*this = append(*this, &waveItem)
	}
	return nil
}

// func (this *BuffEffectInfos) Decode(str string) error {
// 	*this = make(BuffEffectInfos, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		list := strings.Split(strings.TrimSpace(v), COMMA)
// 		if len(list) < 3 {
// 			// return errors.New(v + "属性信息格式错误")
// 			continue
// 		}
// 		var err error
// 		var values [3]int
// 		for i := 0; i < 3; i++ {
// 			values[i], err = strconv.Atoi(list[i])
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		var buffEffectInfo BuffEffectInfo
// 		buffEffectInfo.Id = values[0]
// 		buffEffectInfo.Rate = values[1]
// 		buffEffectInfo.EffectType = values[2]
// 		*this = append(*this, &buffEffectInfo)
// 	}
// 	return nil
// }

// func (this *ItemFloatInfos) Decode(str string) error {
// 	*this = make(ItemFloatInfos, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(str, SEMICOLON)
// 	for _, v := range infoList {
// 		ranges := strings.Split(v, COMMA)
// 		if len(ranges) < 2 {
// 			return errors.New(str + " ItemFloatInfos 属性信息格式错误")
// 		}
// 		var itemFloatInfo ItemFloatInfo
// 		itemFloatInfo.ItemId, _ = strconv.Atoi(ranges[0])
// 		itemFloatInfo.Count, _ = strconv.ParseFloat(ranges[1], 64)
// 		*this = append(*this, &itemFloatInfo)
// 	}
// 	return nil
// }

// func (this IntSlice) String(sep string) string {
// 	var arrStr = make([]string, len(this))
// 	for i, v := range this {
// 		arrStr[i] = strconv.Itoa(v)
// 	}
// 	return strings.Join(arrStr, sep)
// }

func (this *Condition) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList, err := IntSliceFromString(str, COMMA)
	if err != nil {
		return err
	}
	c, err := NewCondition(infoList)
	if err != nil {
		return err
	}

	*this = c
	return nil
}

func NewCondition(infoList []int) (Condition, error) {
	l := len(infoList)
	var c Condition
	if l < 2 {
		return c, nil
		//return c, errors.New("condition 属性信息格式错误")
	}
	c.K, c.V = infoList[0], infoList[1]
	if l > 2 {
		if l%2 != 0 {
			return c, errors.New("condition 长度必须是偶数")
		}
		subs := make(map[int]int, l/2)
		for i := 2; i < l; i++ {
			subKey, subValue := infoList[i], infoList[i+1]
			subs[subKey] = subValue
			i++
		}
		c.Subs = subs
	}
	return c, nil
}

func (this *Conditions) Decode(str string) error {
	*this = make(Conditions, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
	if len(infoList) == 0 {
		return nil
	}
	for _, one := range infoList {
		var c Condition
		err := c.Decode(one)
		if err != nil {
			return err
		}
		*this = append(*this, &c)
	}
	return nil
}

// //Decode
// func (this *Weight) Decode(str string) error {
// 	var m IntMap
// 	m.Decode(str)
// 	weightLen := len(m)
// 	if weightLen == 0 {
// 		return nil
// 	}
// 	v := Weight{OriginMap: m}
// 	v.WeightArr = make([]int, weightLen)
// 	v.KeyArr = make([]int, weightLen)
// 	j := 0
// 	sum := 0
// 	for k, weight := range m {
// 		if weight < 0 {
// 			return fmt.Errorf("Weight:Decode:wrong weight:%d", weight)
// 		}
// 		v.WeightArr[j] = weight
// 		v.KeyArr[j] = k
// 		j += 1
// 		sum += weight
// 	}
// 	v.Sum = sum
// 	*this = v
// 	return nil
// }

// //Rand 随机一个权重中的值,(值,是否成功得到)
// func (this Weight) Rand() (int, bool) {
// 	if this.Sum == 0 {
// 		return 0, false
// 	}
// 	randomIndex := util.RouletteSelectWithSum(this.WeightArr, this.Sum)
// 	if randomIndex < 0 {
// 		return 0, false
// 	}
// 	return this.KeyArr[randomIndex], true
// }

// //TruncateByKey
// //从有序（小-大)的propInfos中，取出比key小的最大值。
// func (this PropInfos) TruncateByKey(currentKey int) (key, value int) {
// 	l := len(this)
// 	if l == 0 {
// 		return
// 	}
// 	last := this[0]
// 	for _, one := range this {
// 		if one.K > currentKey {
// 			break
// 		}
// 		last = one
// 	}
// 	return last.K, last.V
// }

// func (this *ProbItem) ToItemInfo() *ItemInfo {
// 	return &ItemInfo{ItemId: this.Id, Count: this.Count}
// }

// func (this ProbItems) Clone() ProbItems {
// 	results := make(ProbItems, 0, len(this))
// 	results = append(results, this...)
// 	return results
// }

// //这个转盘是每次抽走一个奖品就会重新计算weight
// func (this ProbItems) RandUniqueViaWeight(count int) ItemInfos {
// 	var totalRate int
// 	for _, item := range this {
// 		totalRate += item.Prob
// 	}

// 	awardItemsList := this.Clone()

// 	results := make(ItemInfos, 0)
// 	for i := 1; i <= count; i++ {
// 		currentRandom := 0
// 		random := rand.Intn(totalRate)
// 		for index, item := range awardItemsList {
// 			currentRandom += item.Prob
// 			if currentRandom >= random {
// 				results = append(results, item.ToItemInfo())
// 				awardItemsList = append(awardItemsList[:index], awardItemsList[index+1:]...)
// 				totalRate = totalRate - item.Prob
// 				break
// 			}
// 		}
// 	}
// 	return results
// }

func (this *Condition) Clone() *Condition {
	c := &Condition{}
	c.K = this.K
	c.V = this.V
	if len(this.Subs) > 0 {
		c.Subs = make(map[int]int, len(this.Subs))
		for k, v := range this.Subs {
			c.Subs[k] = v
		}
	}
	return c
}

func (this Conditions) Clone() Conditions {
	results := make(Conditions, 0, len(this))
	for _, v := range this {
		results = append(results, v.Clone())
	}
	return results
}

//CloneAndAdd
func (this Conditions) CloneAndAdd(m map[int]int) Conditions {
	conditions := this.Clone()
	for _, v := range conditions {
		v.V += m[v.K]
		v.V--
	}
	return conditions
}

func (this *DateTime) String() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", this.Year, this.Month, this.Day, this.Hour, this.Minute, this.Second)
}

func (this *DateTime) ToTime() time.Time {
	return time.Date(this.Year, time.Month(this.Month), this.Day, this.Hour, this.Minute, this.Second, 0, time.Local)
}

func (this *DateTime) Unix() int64 {
	if this.Month == 0 {
		return 0
	}
	return this.ToTime().Unix()
}

func (this *DateTime) Unix32() int32 {
	return int32(this.Unix())
}

func (this *DateTime) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(strings.TrimSpace(str), SPACE)
	if len(infoList) < 2 {
		return errors.New(str + " DateTime 属性信息格式错误")
	}

	infoListLeft := strings.Split(strings.TrimSpace(infoList[0]), HLINE)
	if len(infoListLeft) < 3 {
		return errors.New(str + " DateTime 属性格式错误")
	}
	this.Year, _ = strconv.Atoi(infoListLeft[0])
	this.Month, _ = strconv.Atoi(infoListLeft[1])
	this.Day, _ = strconv.Atoi(infoListLeft[2])

	infoListRight := strings.Split(strings.TrimSpace(infoList[1]), COLON)
	if len(infoListRight) < 3 {
		return errors.New(str + " DateTime 属性信息格式错误")
	}
	this.Hour, _ = strconv.Atoi(infoListRight[0])
	this.Minute, _ = strconv.Atoi(infoListRight[1])
	this.Second, _ = strconv.Atoi(infoListRight[2])

	return nil
}

//仅方便客户端获取高级战机的配置
type DragonEnhanceCost struct {
	Ratio float64
	Plus  int
	Cost  ItemInfo
}
