package configs

type GridMap struct {
	Type int `json:"type"`
	Area int `json:"area"` //区域
}

type AreaMapCfg struct {
	AreaMap []string `json:"areaMap"`
}

type MacpCfg struct {
	GridMap     map[string]GridMap `json:"gridMap"`
	AreaMap     map[int][]string   `json:"areaMap"`
	RoadW       float32            `json:"roadW"`
	RoadH       float32            `json:"roadH"`
	OriginY     float32            `json:"originY"`
	OriginX     float32            `json:"originX"`
	RoadMaxNumX int                `json:"roadMaxNumX"`
	RoadMaxNumY int                `json:"roadMaxNumY"`
	//RoadData    string              `json:"roadData"`
}

//----------------有关方法

// func (m *GameDb) GetMapModel(id int) *BuildingModelCfg {
// 	return m.BuildingModels[id]
// }
