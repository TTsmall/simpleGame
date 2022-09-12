package bag

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/TTsmall/wawaji_pub_hmhelper/common"
)

/*
带UID的背包
每一个成员都是一个实例，会有自己的UID
*/

type BagExMDOption func(md *BagExMD)

func BagSetNewItemF(f NewBagExItem) BagExMDOption {
	return func(md *BagExMD) {
		// fmt.Println("BagSetNewItemF")
		md.newItemF = f
	}
}
func BagSetExistChange(f ExistChange) BagExMDOption {
	return func(md *BagExMD) {
		// fmt.Println("BagSetExistChange")
		md.existChangeF = f
	}
}

//加载数据
func BagSetUnmarshalF(f func(buf []byte) (result map[int]IBagExItem)) BagExMDOption {
	return func(md *BagExMD) {
		md.Items = f(md.tmpjson)
		md.tmpjson = nil
	}
}

func NewBagExMD(opts ...BagExMDOption) InitBag {
	return func() IBag {
		md := new(BagExMD)
		md.Items = make(map[int]IBagExItem)
		md.MaxUid = 0
		// fmt.Println("NewBagExMD", len(opts))
		for index := range opts {
			opts[index](md)
		}
		// fmt.Printf("NewBagExMD:%+v\n", md)
		return md
	}
}
func (md *BagExMD) SetBag(opts ...BagExMDOption) {
	for index := range opts {
		opts[index](md)
	}
}

//生成拿成员的方法
type NewBagExItem func(uid int, conf *ItemCfg) IBagExItem

//存在时，调用的修改方法
type ExistChange func(cf *ItemCfg, item IBagExItem) common.ItemInfos

type BagExMD struct {
	Items  map[int]IBagExItem //放实例的地方
	MaxUid int                //用来管理自增长的ID

	newItemF     NewBagExItem //拿成员的方法
	existChangeF ExistChange  //UID存在时的修改方法
	tmpjson      []byte       //临时从DB中读出来的数据
}

// 实现driver.Valuer接口
func (md BagExMD) Value() (driver.Value, error) {
	if md.Items == nil {
		md.Items = make(map[int]IBagExItem)
	}
	return json.Marshal(md.Items)
}

func (md *BagExMD) String() string {
	return fmt.Sprintf("%+v", *md)
}

// 实现sql.Scanner接口
func (md *BagExMD) Scan(val interface{}) (err error) {
	// fmt.Printf("BagExMD.Scan:%+v", md)
	if buf, ok := val.([]byte); ok {
		if len(buf) == 0 {
			buf = []byte("{}")
		}
		// err = json.Unmarshal(buf, &md.Items)
		md.tmpjson = buf
	}
	// for k := range md.Items {
	// 	if md.MaxUid < k {
	// 		md.MaxUid = k
	// 	}
	// }
	return
}

//添加物品，背包内需要实现logtype的逻辑
func (md *BagExMD) AddItem(cf *ItemCfg, itemid, num, logtype int) (result IBagItem, nitems common.ItemInfos) {
	if cf == nil {
		cf = BagEx.ItemConfEx[itemid]
	}
	fmt.Printf("%+v,%+v", md, cf)
	item := md.newItemF(md.MaxUid+1, cf)
	if _, ok := md.Items[item.GetUid()]; ok {
		//东西存在了，做存在的处理
		nitems = md.existChangeF(cf, item)
	} else {
		md.Items[item.GetUid()] = item
		md.MaxUid = item.GetUid()
		result = item
	}
	return
}

//删除物品,以UID
func (md *BagExMD) DelItemByUID(uid, logtype int) (result IBagExItem) {
	if item, ok := md.Items[uid]; ok {
		delete(md.Items, uid)
		return item
	}
	return nil
}

//拿对应的实例
func (md *BagExMD) GetItemByUID(uid int) (IBagExItem, bool) {
	item, ok := md.Items[uid]
	return item, ok
}
