package helper

// import (
// 	"errors"
// 	"fmt"
// 	"math/rand"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/TTsmall/wawaji_pub_hmhelper/common"
// )

// const (
// 	SEMICOLON = ";"
// 	COMMA     = ","
// 	COLON     = ":"
// 	PIPE      = "|"
// 	SPACE     = " "
// 	HLINE     = "-"
// )

// type ItemInfo struct {
// 	ItemId int `client:"key"`
// 	Count  int `client:"value"`
// }

// type FloatItemInfo struct {
// 	ItemId int     `client:"key"`
// 	Count  float64 `client:"value"`
// }
// type ItemInfoProb struct {
// 	ItemId int `client:"itemId"`
// 	Count  int `client:"count"`
// 	Prob   int `client:"prob"`
// }

// type PropInfo struct {
// 	K int `client:"key"`
// 	V int `client:"value"`
// }

// // 属性增益，兼容浮点和整数
// type PropGain struct {
// 	Type  int     `client:"type"`
// 	Float float64 `client:"value"`
// 	Int   int
// }

// type WaveItem struct {
// 	Wave  int `client:"wave"`  // 波数
// 	Id    int `client:"id"`    // 物品ID
// 	Count int `client:"count"` // 数量
// }

// type Condition struct {
// 	K    int         `client:"key"`
// 	V    int         `client:"value"`
// 	Subs map[int]int `client:"subs"`
// }

// type DateTime struct {
// 	Year   int
// 	Month  int
// 	Day    int
// 	Hour   int
// 	Minute int
// 	Second int
// }

// type ItemInfoProbs []*ItemInfoProb
// type Conditions []*Condition
// type ItemInfos []*ItemInfo
// type FloatItemInfos []*FloatItemInfo

// type PropInfos []*PropInfo
// type PropGains []*PropGain

// type IntSlice []int
// type FloatSlice []float64

// func (s IntSlice) Len() int           { return len(s) }
// func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
// func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// type WaveItems []*WaveItem

// type IntMap map[int]int

// func (this *IntSlice) Decode(str string) error {
// 	ints, err := common.IntSliceFromString(str, ",")
// 	if err != nil {
// 		return err
// 	}
// 	*this = IntSlice(ints)
// 	return nil
// }

// func (this IntSlice) ToInt32Slice() []int32 {
// 	l := len(this)
// 	ret := make([]int32, l)
// 	if l == 0 {
// 		return ret
// 	}
// 	for i := 0; i < l; i++ {
// 		ret[i] = int32(this[i])
// 	}
// 	return ret
// }

// func (this *FloatSlice) Decode(str string) error {
// 	floats, err := common.FloatSliceFromString(str, ",")
// 	if err != nil {
// 		return err
// 	}
// 	*this = FloatSlice(floats)
// 	return nil
// }

// // func (this IntSlice) GetOrLast(index int) int {
// // 	l := len(this)
// // 	if index < l {
// // 		return this[index]
// // 	}
// // 	if l == 0 {
// // 		return 0
// // 	}
// // 	return this[l-1]
// // }

// // func (this *StringSlice) Decode(str string) error {
// // 	if len(strings.TrimSpace(str)) == 0 {
// // 		*this = make([]string, 0)
// // 		return nil
// // 	}
// // 	*this = strings.Split(str, ";")
// // 	return nil
// // }

// func (this ItemInfos) Times(times float64) ItemInfos {
// 	result := make(ItemInfos, 0, len(this))
// 	for _, itemInfo := range this {
// 		count := int(float64(itemInfo.Count) * times)
// 		if count > 0 {
// 			result = append(result, &ItemInfo{ItemId: itemInfo.ItemId, Count: count})
// 		}
// 	}
// 	return result
// }

// func (this ItemInfos) ToInt32Map() map[int32]int32 {
// 	result := make(map[int32]int32, len(this))
// 	for _, itemInfo := range this {
// 		result[int32(itemInfo.ItemId)] += int32(itemInfo.Count)
// 	}
// 	return result
// }

// func (this ItemInfos) Get(itemId int) *ItemInfo {
// 	for _, itemInfo := range this {
// 		if itemInfo.ItemId == itemId {
// 			return itemInfo
// 		}
// 	}
// 	return nil
// }

// //GetOrLast 取index位置的或者最后一个。
// // func (this ItemInfos) GetOrLast(index int) *ItemInfo {
// // 	l := len(this)
// // 	if l == 0 {
// // 		return nil
// // 	}
// // 	if index < l {
// // 		return this[index]
// // 	}
// // 	return this[l-1]
// // }

// func (this *ItemInfos) Decode(str string) error {
// 	*this = make(ItemInfos, 0)
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
// 			return errors.New(v + "物品信息格式错误")
// 		}
// 		itemId, _ := strconv.Atoi(list[0])
// 		var itemInfo ItemInfo
// 		itemInfo.ItemId = itemId
// 		itemCount, err := strconv.Atoi(list[1])
// 		if err != nil {
// 			return err
// 		}
// 		itemInfo.Count = itemCount
// 		*this = append(*this, &itemInfo)
// 	}
// 	return nil
// }

// func (this *FloatItemInfos) Decode(str string) error {
// 	*this = make(FloatItemInfos, 0)
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
// 			return errors.New(v + "物品信息格式错误")
// 		}
// 		itemId, _ := strconv.Atoi(list[0])
// 		var floatItemInfo FloatItemInfo
// 		floatItemInfo.ItemId = itemId
// 		itemCount, err := strconv.ParseFloat(list[1], 64)
// 		if err != nil {
// 			return err
// 		}
// 		floatItemInfo.Count = itemCount
// 		*this = append(*this, &floatItemInfo)
// 	}
// 	return nil
// }

// // func (this *ItemInfosSlice) Decode(str string) error {
// // 	*this = make(ItemInfosSlice, 0)
// // 	if len(str) == 0 {
// // 		return nil
// // 	}
// // 	infoList := strings.Split(strings.TrimSpace(str), PIPE)
// // 	if len(infoList) == 0 {
// // 		return nil
// // 	}
// // 	for _, v := range infoList {
// // 		itemInfos := &ItemInfos{}
// // 		err := itemInfos.Decode(v)
// // 		if err != nil {
// // 			return err
// // 		}
// // 		*this = append(*this, *itemInfos)
// // 	}
// // 	return nil
// // }

// func (this *PropInfos) Decode(str string) error {
// 	*this = make(PropInfos, 0)
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
// 			return errors.New(v + "属性信息格式错误")
// 		}
// 		k, err := strconv.Atoi(list[0])
// 		if err != nil {
// 			return err
// 		}
// 		var propInfo PropInfo
// 		propInfo.K = k
// 		propInfo.V, err = strconv.Atoi(list[1])
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, &propInfo)
// 	}
// 	return nil
// }

// func (this PropInfos) Rand() *PropInfo {
// 	var sum = 0
// 	var weight = rand.Intn(100)
// 	for _, info := range this {
// 		sum += info.V
// 		if weight < sum {
// 			return info
// 		}
// 	}
// 	return this[rand.Intn(len(this))]
// }

// func (this *PropGains) Decode(str string) error {
// 	*this = make(PropGains, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		gain := &PropGain{}
// 		err := gain.Decode(v)
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, gain)
// 	}
// 	return nil
// }

// func (this *IntMap) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}
// 	*this = make(IntMap)
// 	for _, v := range infoList {
// 		list := strings.Split(strings.TrimSpace(v), COMMA)
// 		if len(list) != 2 {
// 			return errors.New(v + "IntMap 属性信息格式错误")
// 		}

// 		k, err := strconv.Atoi(list[0])
// 		if err != nil {
// 			return err
// 		}
// 		if _, ok := (*this)[k]; ok {
// 			return errors.New(v + "IntMap 属性重复")
// 		}
// 		v, err := strconv.Atoi(list[1])
// 		if err != nil {
// 			return err
// 		}
// 		(*this)[k] = v
// 	}
// 	return nil

// }

// func (this *PropInfo) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(str, COMMA)
// 	if len(infoList) < 2 {
// 		return errors.New(str + " PropInfo 属性信息格式错误")
// 	}
// 	var propInfo PropInfo
// 	propInfo.K, _ = strconv.Atoi(infoList[0])
// 	propInfo.V, _ = strconv.Atoi(infoList[1])
// 	*this = propInfo
// 	return nil
// }

// func (this *PropGain) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(str, COMMA)
// 	if len(infoList) < 2 {
// 		return errors.New(str + " PropPlus 属性信息格式错误")
// 	}
// 	var err error
// 	this.Type, err = strconv.Atoi(infoList[0])
// 	if err != nil {
// 		return err
// 	}
// 	f, err := strconv.ParseFloat(infoList[1], 64)
// 	if err != nil {
// 		return err
// 	}
// 	this.Float = float64(f)
// 	this.Int = int(f)
// 	return nil
// }

// func (this *ItemInfoProb) ToItemInfo() *ItemInfo {
// 	return &ItemInfo{ItemId: this.ItemId, Count: this.Count}
// }

// func (this *ItemInfoProbs) Decode(str string) error {
// 	*this = make(ItemInfoProbs, 0)
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
// 		var probItem ItemInfoProb
// 		probItem.ItemId = id
// 		probItem.Count, err = strconv.Atoi(list[1])
// 		probItem.Prob, err = strconv.Atoi(list[2])
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, &probItem)
// 	}
// 	return nil
// }

// func (this ItemInfoProbs) Clone() ItemInfoProbs {
// 	results := make(ItemInfoProbs, 0, len(this))
// 	results = append(results, this...)
// 	return results
// }

// //这个转盘是每次抽走一个奖品就会重新计算weight
// func (this ItemInfoProbs) RandUniqueViaWeight(count int) ItemInfos {
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

// //这个转盘是每次抽走一个奖品就会不重新计算weight
// func (this ItemInfoProbs) RandUniqueViaAllWeight(count int) ItemInfos {
// 	if count == 0 {
// 		return nil
// 	}
// 	var totalRate int
// 	for _, item := range this {
// 		totalRate += item.Prob
// 	}
// 	awardItemsList := this.Clone()
// 	results := make(ItemInfos, 0)
// 	for i := 1; i <= count; i++ {
// 		currentRandom := 0
// 		random := rand.Intn(totalRate)
// 		for _, item := range awardItemsList {
// 			currentRandom += item.Prob
// 			if currentRandom >= random {
// 				results = append(results, item.ToItemInfo())
// 				break
// 			}
// 		}
// 	}
// 	return results
// }

// func (this *ItemInfoProb) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(str, COMMA)
// 	if len(infoList) < 3 {
// 		return errors.New(str + " 属性信息格式错误")
// 	}
// 	var itemInfo ItemInfoProb
// 	itemInfo.ItemId, _ = strconv.Atoi(infoList[0])
// 	itemInfo.Count, _ = strconv.Atoi(infoList[1])
// 	itemInfo.Prob, _ = strconv.Atoi(infoList[2])
// 	*this = itemInfo
// 	return nil
// }

// func (this *ItemInfo) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(str, COMMA)
// 	if len(infoList) < 2 {
// 		return errors.New(str + " 属性信息格式错误")
// 	}
// 	var itemInfo ItemInfo
// 	itemInfo.ItemId, _ = strconv.Atoi(infoList[0])
// 	itemInfo.Count, _ = strconv.Atoi(infoList[1])
// 	*this = itemInfo
// 	return nil
// }

// func (this *WaveItem) Decode(str string) error {

// 	values, err := common.IntSliceFromString(strings.TrimSpace(str), COMMA)
// 	if err != nil || len(values) < 3 {
// 		return fmt.Errorf("WaveItem:Decode bad str:%s,err:%v", str, err)
// 	}
// 	var waveItem WaveItem
// 	waveItem.Wave = values[0]
// 	waveItem.Id = values[1]
// 	waveItem.Count = values[2]
// 	*this = waveItem
// 	return nil
// }

// // 复用ProbItems的Decode
// func (this *WaveItems) Decode(str string) error {
// 	*this = make(WaveItems, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}

// 	for _, v := range infoList {
// 		var waveItem WaveItem
// 		err := waveItem.Decode(v)
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, &waveItem)
// 	}
// 	return nil
// }

// func (this *Condition) Decode(str string) error {
// 	if len(strings.TrimSpace(str)) == 0 {
// 		return nil
// 	}
// 	infoList, err := common.IntSliceFromString(str, COMMA)
// 	if err != nil {
// 		return err
// 	}
// 	c, err := NewCondition(infoList)
// 	if err != nil {
// 		return err
// 	}

// 	*this = c
// 	return nil
// }

// func NewCondition(infoList []int) (Condition, error) {
// 	l := len(infoList)
// 	var c Condition
// 	if l < 2 {
// 		return c, nil
// 		//return c, errors.New("condition 属性信息格式错误")
// 	}
// 	c.K, c.V = infoList[0], infoList[1]
// 	if l > 2 {
// 		if l%2 != 0 {
// 			return c, errors.New("condition 长度必须是偶数")
// 		}
// 		subs := make(map[int]int, l/2)
// 		for i := 2; i < l; i++ {
// 			subKey, subValue := infoList[i], infoList[i+1]
// 			subs[subKey] = subValue
// 			i++
// 		}
// 		c.Subs = subs
// 	}
// 	return c, nil
// }

// func (this *Conditions) Decode(str string) error {
// 	*this = make(Conditions, 0)
// 	if len(str) == 0 {
// 		return nil
// 	}
// 	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), SEMICOLON), SEMICOLON)
// 	if len(infoList) == 0 {
// 		return nil
// 	}
// 	for _, one := range infoList {
// 		var c Condition
// 		err := c.Decode(one)
// 		if err != nil {
// 			return err
// 		}
// 		*this = append(*this, &c)
// 	}
// 	return nil
// }

// func (this *Condition) Clone() *Condition {
// 	c := &Condition{}
// 	c.K = this.K
// 	c.V = this.V
// 	if len(this.Subs) > 0 {
// 		c.Subs = make(map[int]int, len(this.Subs))
// 		for k, v := range this.Subs {
// 			c.Subs[k] = v
// 		}
// 	}
// 	return c
// }

// func (this Conditions) Clone() Conditions {
// 	results := make(Conditions, 0, len(this))
// 	for _, v := range this {
// 		results = append(results, v.Clone())
// 	}
// 	return results
// }

// //CloneAndAdd
// func (this Conditions) CloneAndAdd(m map[int]int) Conditions {
// 	conditions := this.Clone()
// 	for _, v := range conditions {
// 		v.V += m[v.K]
// 		v.V--
// 	}
// 	return conditions
// }

// //取出要求等级
// func (this Condition) PickLevel() int {
// 	if this.K == 2 {
// 		return this.V
// 	}
// 	return 0
// }

// //取出要求等级
// func (this Conditions) PickLevel() int {
// 	for _, c := range this {
// 		lv := c.PickLevel()
// 		if lv > 0 {
// 			return lv
// 		}
// 	}
// 	return 0
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

// //仅方便客户端获取高级战机的配置
// type DragonEnhanceCost struct {
// 	Ratio float64
// 	Plus  int
// 	Cost  ItemInfo
// }
