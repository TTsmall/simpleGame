package common

import (
	"fmt"
	"strconv"
	"strings"
)

type MapIntSlice map[int]IntSlice
type IntSlice []int
type StringSlice []string
type FloatSlice []float64

func FloatSliceFromString(str, sep string) (FloatSlice, error) {
	if len(str) == 0 {
		return FloatSlice(make([]float64, 0)), nil
	}
	strs := strings.Split(str, sep)
	length := len(strs)
	var err error
	var res float64
	var result = make(FloatSlice, length)
	for i := 0; i < length; i++ {
		if len(strs[i]) == 0 {
			continue
		}
		res, err = strconv.ParseFloat(strs[i], 64)
		result[i] = float64(res)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func IntSliceFromString(str string, sep string) (IntSlice, error) {
	if len(str) == 0 {
		return IntSlice(make([]int, 0)), nil
	}
	strs := strings.Split(str, sep)
	var err error
	var result = make(IntSlice, len(strs))
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) == 0 {
			continue
		}
		result[i], err = strconv.Atoi(strs[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this IntSlice) Index(element int) int {
	for i, v := range this {
		if v == element {
			return i
		}
	}
	return -1
}

func Int32IndexOf(arr []int32, element int32) int {
	for i, v := range arr {
		if v == element {
			return i
		}
	}
	return -1
}

func (this IntSlice) RemoveIndex(index int) IntSlice {
	if index < 0 || index >= len(this) {
		return this
	}
	return append(this[:index], this[index+1:]...)
}

func (this IntSlice) RemoveElement(element int) IntSlice {
	for i, v := range this {
		if v == element {
			return append(this[:i], this[i+1:]...)
		}
	}
	return this
}

func (this IntSlice) Add(element int) IntSlice {
	return append(this, element)
}

func (this IntSlice) AddUnique(element int) IntSlice {
	if this.Index(element) < 0 {
		return this
	}
	return append(this, element)
}

func (this IntSlice) String(sep string) string {
	var arrStr = make([]string, len(this))
	for i, v := range this {
		arrStr[i] = strconv.Itoa(v)
	}
	return strings.Join(arrStr, sep)
}

func ConvertIntSlice2Int32Slice(origin []int) []int32 {
	ret := make([]int32, len(origin))
	for i, v := range origin {
		ret[i] = int32(v)
	}
	return ret
}

func ConvertInt32Slice2IntSlice(origin []int32) []int {
	ret := make([]int, len(origin))
	for i, v := range origin {
		ret[i] = int(v)
	}
	return ret
}

func SliceIntUnique(origin []int) []int {
	ret := make([]int, 0, len(origin))
	tempMap := make(map[int]struct{}, len(origin))
	for _, v := range origin {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}

func SliceInt32Unique(origin []int32) []int32 {
	ret := make([]int32, 0, len(origin))
	tempMap := make(map[int32]struct{}, len(origin))
	for _, v := range origin {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}

func SliceStringUnique(origin []string) []string {
	ret := make([]string, 0, len(origin))
	tempMap := make(map[string]struct{}, len(origin))
	for _, v := range origin {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}

func JoinStringSlice(a []string, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = v
	}
	return strings.Join(b, sep)
}

func JoinIntSlice(a []int, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func JoinInt32Slice(a []int32, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, sep)
}

func InterfaceSlice2StringSlice(arr []interface{}) []string {
	strArr := make([]string, len(arr))
	for index, v := range arr {
		strArr[index] = fmt.Sprint(v)
	}
	return strArr
}

func MapInt32Keys(m map[int]int) []int32 {
	r := make([]int32, len(m))
	i := 0
	for k := range m {
		r[i] = int32(k)
		i++
	}
	return r
}

func SliceToMap(m []int) map[int]int {
	r := make(map[int]int)
	for k, v := range m {
		r[v] = k
	}
	return r
}

func MapToMapInt32(m map[int]int) map[int32]int32 {
	r := make(map[int32]int32)
	for k, v := range m {
		r[int32(k)] = int32(v)
	}
	return r
}

func Map32ToMapInt(m map[int32]int32) map[int]int {
	r := make(map[int]int)
	for k, v := range m {
		r[int(k)] = int(v)
	}
	return r
}

func CombatMap(sourM, destM map[int]int) {
	for k, v := range sourM {
		destM[k] = v
	}
}

func SliceStringToInt(strs []string) []int {
	destS := make([]int, len(strs))
	for k, v := range strs {
		valueInt, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		destS[k] = valueInt
	}
	return destS
}
