package configs

import (
	"fmt"

	"github.com/TTsmall/wawaji_pub_hmhelper/log"
)

func (m *GameDb) Patch() {
	fmt.Println("====GameDb =Patch=====")
	m.genBuildNumLimits()
	m.genBuildLvel()
}

func (m *GameDb) genBuildNumLimits() {
	m.BuildNumLimits = make(map[int]int)
	m.BuildMinjuByNumType = make(map[int]*BuildNumConditionCfg)
	m.BuildMinjuWrokers = make(map[int][]int)
	for _, v := range m.BuildNumConditions {
		m.BuildNumLimits[v.BuildType] = v.MaxBulidNum
		if v.BuildType == 52 {
			m.BuildMinjuByNumType[v.BuildNum] = v
			if len(v.Npc) > 0 {
				m.BuildMinjuWrokers[v.Npc[0]] = v.Npc
			} else {
				log.Error("genBuildNumLimits building npc no")
			}

		}
	}
}

func (m *GameDb) genBuildLvel() {
	fmt.Println("======= GameDbgenBuildLvel")
	m.BuildLevelByItemType = make(map[int]map[int]*BuildCfg)
	for _, v := range m.Builds {
		if len(m.BuildLevelByItemType[v.Itemtype]) == 0 {
			m.BuildLevelByItemType[v.Itemtype] = make(map[int]*BuildCfg)
		}
		m.BuildLevelByItemType[v.Itemtype][v.Level] = v
	}
	fmt.Println("======= m.BuildLevelByItemType[v.Itemtype]=", m.BuildLevelByItemType[3300])
}
