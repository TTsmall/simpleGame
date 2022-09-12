package bag

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"wawaji_pub/hmhelper/common"
)

/*
基础仓库
ID，数量的仓库
*/
type BagMDOption func(md *BagMD)

//设置基础背包的总上限
func BagSetTotalMax(max int) BagMDOption {
	return func(md *BagMD) {
		md.TotalMax = max
		if md.TotalMax > 0 {
			md.CurrNum = 0
			for _, v := range md.Items.Data {
				md.CurrNum += v
			}
		}
	}
}

func NewBagMD(opts ...BagMDOption) InitBag {
	return func() IBag {
		md := new(BagMD)
		md.Items = common.NewBaseDataString("")
		md.TotalMax = -1
		for index := range opts {
			opts[index](md)
		}
		return md
	}
}

func (md *BagMD) SetBagMD(opts ...BagMDOption) {
	for index := range opts {
		opts[index](md)
	}
}

//在User上的数据
type BagMD struct {
	Items    *common.BaseData
	TotalMax int `json:"-"` //所有物品的总上限
	CurrNum  int `json:"-"` //当前总数
}

// 实现driver.Valuer接口
func (md BagMD) Value() (driver.Value, error) {
	if md.Items == nil {
		md.Items = common.NewBaseDataString("")
	}
	return []byte(md.Items.ToString()), nil
}

func (md *BagMD) String() string {
	return fmt.Sprintf("%+v", *md)
}

// 实现sql.Scanner接口
func (md *BagMD) Scan(val interface{}) (err error) {
	if buf, ok := val.([]byte); ok {
		if len(buf) == 0 {
			buf = []byte("")
		}
		md.Items = common.NewBaseDataString(string(buf))
	}
	if md.TotalMax > 0 {
		md.CurrNum = 0
		for _, v := range md.Items.Data {
			md.CurrNum += v
		}
	}
	return
}
func (md BagMD) MarshalJSON() ([]byte, error) {
	fmt.Println("BagMD.MarshalJSON=", md.Items.ToString())
	return []byte(fmt.Sprint("\"", md.Items.ToString(), "\"")), nil
}

func (md *BagMD) UnmarshalJSON(buf []byte) error {
	fmt.Println("BagMD.UnmarshalJSON", string(buf))
	data := strings.ReplaceAll(string(buf), "\"", "")
	md.Items = common.NewBaseDataString(data)
	if md.TotalMax > 0 {
		md.CurrNum = 0
		for _, v := range md.Items.Data {
			md.CurrNum += v
		}
	}
	return nil
}

func (md *BagMD) AddItem(cf *ItemCfg, itemid, num, logtype int) (result IBagItem, nitems common.ItemInfos) {
	if cf == nil {
		cf = BagEx.ItemConfEx[itemid]
	}
	res := new(ResultItem)
	res.ItemID = itemid
	res.Num = md.Items.GetNumByKey(itemid)
	if md.TotalMax > 0 {
		if md.CurrNum+num > md.TotalMax {
			num = cf.MaxNum - res.Num
		}
	}
	if cf.MaxNum > 0 {
		if res.Num+num > cf.MaxNum {
			num = cf.MaxNum - res.Num
		}
	}
	res.Delta = num
	md.Items.UpData(itemid, num)
	res.Num += num
	return res, nil
}

//删除物品,以ItemID
func (md *BagMD) DelItemByID(cf *ItemCfg, itemid, num, logtype int) (result IBagItem) {
	if cf == nil {
		cf = BagEx.ItemConfEx[itemid]
	}
	res := new(ResultItem)
	res.ItemID = itemid
	res.Num = md.Items.GetNumByKey(itemid)
	res.Delta = -num
	md.Items.UpData(itemid, -num)
	res.Num += -num
	return res
}

//检查道具数量，返回true表示满足，返回false表示不满足
func (md *BagMD) CheckItem(itemid, num int) bool {
	if md.Items.GetNumByKey(itemid) < num {
		return false
	}
	return true
}

//乘余总数量
func (md *BagMD) GetSurplusNum() int {
	return md.TotalMax - md.CurrNum
}

//道具ID对应的数量
func (md *BagMD) GetItemNum(itemid int) int {
	return md.Items.GetNumByKey(itemid)
}
