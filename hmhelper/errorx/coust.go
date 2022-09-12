package errorx

//给客户端的返回,-99为未知错误,0表示为成功
var (
	LOGIC_Unknown_Error  = &ErrorItem{-2, "逻辑处理错误"}
	Redis_Get_Error      = &ErrorItem{-100, "Redis信息获取失败"}
	Json_Unmarshal_Error = &ErrorItem{-200, "Json解析失败"}
	Center_ErrWithDraw   = &ErrorItem{-158, "申请提现异常"}

	Server_Not_Exist_Error = &ErrorItem{1000, "Server不存在"}
	User_Other_Login_Error = &ErrorItem{1100, "用户在其他地方登录"}
	User_Login_Invalid     = &ErrorItem{1101, "用户登录信息失效"}

	Friend_Not_Exist_Error      = &ErrorItem{1150, "好友不存在"}
	Friend_Not_Add_Error        = &ErrorItem{1151, "没有好友添加信息"}
	Friend_Exist_Error          = &ErrorItem{1152, "好友已存在"}
	Friend_Not_Add_MySelf_Error = &ErrorItem{1153, "不能加自己"}

	Bag_Item_Not_Enough = &ErrorItem{1200, "%s数量不足"}
	Bag_Item_Illegal    = &ErrorItem{1201, "物品非法"}

	Build_Item_Not_Error         = &ErrorItem{1250, "没有此建筑"}
	Build_Item_No_Space          = &ErrorItem{1251, "请挪到其他地方建造"}
	Build_Item_Num_Error         = &ErrorItem{1252, "%s建筑数量已达到上限"}
	Build_People_Not_Error       = &ErrorItem{1253, "建筑居民已达到上限"}
	Build_Hero_Not_Error         = &ErrorItem{1254, "没有居民"}
	Build_Hero_No_Use_Error      = &ErrorItem{1255, "当前不可使用此英雄"}
	Build_Item_User_Not_Error    = &ErrorItem{1256, "地图没有此建筑"}
	Build_Item_Time_Not_Error    = &ErrorItem{1257, "建筑冷却CD中"}
	Produce_Time_Not_Error       = &ErrorItem{1258, "建筑生产冷却CD中"}
	Building_Speed_up_Type_Error = &ErrorItem{1259, "建筑加速消耗类型非法"}
	Building_Minju_Error         = &ErrorItem{1260, "民居生产非法"}
	Occup_Minju_Not_Error        = &ErrorItem{1261, "非居民屋不能安排英雄"}
	Release_Building_No_Err      = &ErrorItem{1262, "当前建筑不能释放建筑工人"}
	LevelUp_Building_No_Err      = &ErrorItem{1263, "当前建筑不能LevelUp"}
	Building_Operate_Err         = &ErrorItem{1264, "当前建筑操作中"}
	Building_Explore_Err         = &ErrorItem{1265, "当前建筑解锁中"}
	Building_UnLock_Err          = &ErrorItem{1266, "当前建筑已经解锁"}
	UnLock_Area_No_Err           = &ErrorItem{1267, "解锁区域不存在"}
	Explore_Area_No_Err          = &ErrorItem{1268, "尚未探索"}
	Explore_Area_Cd_Err          = &ErrorItem{1269, "探索CD"}
	Build_Type_Occup_Error       = &ErrorItem{1270, "建筑类型不能住人"}
)
