package bag

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TTsmall/wawaji_pub_hmhelper/centerserver"
	"github.com/TTsmall/wawaji_pub_hmhelper/common"
)

const (
	//金币日志记录数
	MAX_GOLDDETAIL_RECORDS = 50
)
const (
	//提现
	GOLDLOG_DRAE = 3
	//提现失败退回金币，中台回退的操作
	GOLDLOG_CCFAIRBACK = 100
	//补偿金币，errorgold的操作
	GOLDLOG_MAKEUP = 101
)

type GoldOption func(md *UserGoldBag)

func GoldBagByUser(user IGoldUser) GoldOption {
	return func(md *UserGoldBag) {
		md.user = user
	}
}

func NewGoldBag(opts ...GoldOption) InitBag {
	return func() IBag {
		md := new(UserGoldBag)
		md.Gold = 0
		md.DailyGold = 0
		md.ErrAddGold = 0
		md.GoldLogs = make(GoldDetails, 0, 50)
		for _, opt := range opts {
			opt(md)
		}
		return md
	}
}
func (md *UserGoldBag) SetBagMD(opts ...GoldOption) {
	for index := range opts {
		opts[index](md)
	}
}

//金币背包
type UserGoldBag struct {
	Gold       int         `db:"gold" gorm:"column:gold"`                       //金币数量
	DailyGold  int         `db:"dailygold" gorm:"column:dailygold"`             //每日金币
	ErrAddGold int         `db:"errAddGold" gorm:"column:errAddGold"`           //没有加上的金币
	GoldLogs   GoldDetails `db:"goldDetailValue" gorm:"column:goldDetailValue"` //金币日志
	user       IGoldUser   `db:"-" gorm:"-"`                                    //用户的注册时间
}

//添加物品，背包内需要实现logtype的逻辑
func (bag *UserGoldBag) AddItem(cf *ItemCfg, itemid, num, logtype int) (result IBagItem, nitems common.ItemInfos) {
	if cf == nil {
		cf = BagEx.ItemConfEx[itemid]
	}
	res := new(ResultItem)
	res.ItemID = itemid
	res.Num = bag.Gold

	cfg := BagEx.GoldConfEx[logtype]
	if cfg.IsDaily {

		goldmax := BagEx.GameEx.GetMaxGoldByDayNum(
			common.TimeSubOfDay(bag.user.GetCreateTime()),
		)
		//金币的数据需要读用户的注册时间
		if bag.DailyGold+num > goldmax {
			num = goldmax - bag.DailyGold
		}
		bag.DailyGold += num
		bag.Gold += num
	} else {
		bag.Gold += num
	}
	res.Delta = num
	res.Num += num
	if num > 0 {
		lg := GoldDetail{
			Coin:      int32(num),
			Module:    int32(logtype),
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		bag.GoldLogs = append(bag.GoldLogs, lg)
		if len(bag.GoldLogs) > MAX_GOLDDETAIL_RECORDS {
			st := len(bag.GoldLogs) - MAX_GOLDDETAIL_RECORDS
			bag.GoldLogs = bag.GoldLogs[st:]
		}
		if cfg, ok := BagEx.GoldConfEx[logtype]; ok && logtype != GOLDLOG_CCFAIRBACK {
			if centerserver.CenterEx.AddGoldToCenterControl(bag.user, logtype, num, cfg.CenterKey) != nil {
				bag.ErrAddGold += num
			}
		}

	}

	return res, nil
}

//删除物品
func (bag *UserGoldBag) DelItemByID(cf *ItemCfg, itemid, num, logtype int) (result IBagItem) {
	// cfg := GoldConfEx[logtype]
	if cf == nil {
		cf = BagEx.ItemConfEx[itemid]
	}
	res := new(ResultItem)
	res.ItemID = itemid
	res.Num = bag.Gold
	res.Delta = -num

	bag.Gold += -num
	res.Num += -num
	if num > 0 {
		lg := GoldDetail{
			Coin:      int32(-num),
			Module:    int32(logtype),
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		bag.GoldLogs = append(bag.GoldLogs, lg)
		if len(bag.GoldLogs) > MAX_GOLDDETAIL_RECORDS {
			st := len(bag.GoldLogs) - MAX_GOLDDETAIL_RECORDS
			bag.GoldLogs = bag.GoldLogs[st:]
		}
		if cfg, ok := BagEx.GoldConfEx[logtype]; ok && logtype != GOLDLOG_DRAE {
			if centerserver.CenterEx.AddGoldToCenterControl(bag.user, logtype, num, cfg.CenterKey) != nil {
				bag.ErrAddGold += num
			}
		}
	}
	return res
}

//检查道具数量，返回true表示满足，返回false表示不满足
func (md *UserGoldBag) CheckItem(itemid, num int) bool {
	if md.Gold < num {
		return false
	}
	return true
}

//道具ID对应的数量
func (md *UserGoldBag) GetItemNum(itemid int) int {
	return md.Gold
}

//同步错误金币给中台
func (bag *UserGoldBag) SyncErrGold() {
	if bag.ErrAddGold > 0 {
		if cfg, ok := BagEx.GoldConfEx[GOLDLOG_MAKEUP]; ok {
			if centerserver.CenterEx.AddGoldToCenterControl(bag.user, GOLDLOG_MAKEUP, bag.ErrAddGold, cfg.CenterKey) == nil {
				bag.ErrAddGold = 0
			}
		}
	}
}

type GoldDetail struct {
	Coin      int32  `json:"coin"`
	Module    int32  `json:"module"`
	CreatedAt string `json:"createdAt"`
}

type GoldDetails []GoldDetail

// 实现driver.Valuer接口
func (md GoldDetails) Value() (driver.Value, error) {
	buf, err := json.Marshal(md)
	return buf, err
}

func (md *GoldDetails) String() string {
	return fmt.Sprintf("%+v", *md)
}

// 实现sql.Scanner接口
func (md *GoldDetails) Scan(val interface{}) (err error) {
	if buf, ok := val.([]byte); ok {
		if len(buf) == 0 {
			buf = []byte("[]")
		}
		err = json.Unmarshal(buf, md)
	}
	return
}

//------------------外部接口
type IGoldUser interface {
	//用户的key
	GetOpenID() string
	//拿用户的中台信息
	GetCenter() *centerserver.UserCenterMD
	//获取用户创建字段
	GetCreateTime() time.Time
}
