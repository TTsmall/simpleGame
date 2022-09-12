package common

import "sort"

type sortedMap struct {
	m map[int]int64
	s []int
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func SortedInt64Keys(m map[int]int64) IntSlice {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]int, len(m))
	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

// type sortedMap struct {
// 	m map[string]int
// 	s []string
// }

// func (sm *sortedMap) Len() int {
// 	return len(sm.m)
// }

// func (sm *sortedMap) Less(i, j int) bool {
// 	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
// }

// func (sm *sortedMap) Swap(i, j int) {
// 	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
// }

// func sortedKeys(m map[string]int) []string {
// 	sm := new(sortedMap)
// 	sm.m = m
// 	sm.s = make([]string, len(m))
// 	i := 0
// 	for key, _ := range m {
// 		sm.s[i] = key
// 		i++
// 	}
// 	sort.Sort(sm)
// 	return sm.s
// }

type KVInt struct {
	K int
	V int
}

type KVInts []*KVInt

func (s KVInts) Len() int      { return len(s) }
func (s KVInts) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByName implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Organs value.
type ByK struct{ KVInts }

func (s ByK) Less(i, j int) bool { return s.KVInts[i].K < s.KVInts[j].K }

// ByWeight implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Organs value.
type ByV struct{ KVInts }

func (s ByV) Less(i, j int) bool { return s.KVInts[i].V < s.KVInts[j].V }

type ByVDes struct{ KVInts }

func (s ByVDes) Less(i, j int) bool { return s.KVInts[i].V > s.KVInts[j].V }

func SortKvIntMap(mapIns map[int]int) KVInts {
	var s KVInts
	s = make(KVInts, len(mapIns))
	var i int
	for k, v := range mapIns {
		s[i] = &KVInt{k, v}
		i++
	}
	sort.Sort(ByV{s})
	return s
}

func SortKvIntSlice(s KVInts) KVInts {
	sort.Sort(ByV{s})
	return s
}

func SortKvIntSliceDes(s KVInts) KVInts {
	sort.Sort(ByVDes{s})
	return s
}
func SortKvIntMapDes(mapIns map[int]int) KVInts {
	var s KVInts
	s = make(KVInts, len(mapIns))
	var i int
	for k, v := range mapIns {
		s[i] = &KVInt{k, v}
		i++
	}
	sort.Sort(ByVDes{s})
	return s
}

func SortKvIntMapk(mapIns map[int]int) KVInts {
	var s KVInts
	s = make(KVInts, len(mapIns))
	var i int
	for k, v := range mapIns {
		s[i] = &KVInt{k, v}
		i++
	}
	sort.Sort(ByK{s})
	return s
}
