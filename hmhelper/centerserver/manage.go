package centerserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"wawaji_pub/hmhelper/common"
	"wawaji_pub/hmhelper/errorx"
	"wawaji_pub/hmhelper/log"
	"wawaji_pub/hmhelper/redisinfo"

	"github.com/buguang01/Logger"
	"github.com/buguang01/util"
	// zmutil "code.zm.shzhanmeng.com/zmpub/zmhelper/util"
)

const (
	//提现失败退回金币，中台回退的操作
	GOLDLOG_CCFAIRBACK = 100
)

type CenterOption func(md *CenterManager)

func SetCenterManager(url string) CenterOption {
	return func(md *CenterManager) {
		md.CenterSystemUrl = url
	}
}

func SetCenterManagerGoldID(id int) CenterOption {
	return func(md *CenterManager) {
		md.Item_GoldID = id
	}
}

func NewCenterManager(opts ...CenterOption) *CenterManager {
	CenterEx = new(CenterManager)
	CenterEx.Item_GoldID = 1
	for i := range opts {
		opts[i](CenterEx)
	}
	http.HandleFunc("/api/outer/withdraw/drawCoinToGame", CenterEx.HandleNotifyDrawGold)
	return CenterEx
}

//管理器
var (
	CenterEx *CenterManager
)

type CenterManager struct {
	// zmutil.DefaultModule
	CenterSystemUrl string
	Item_GoldID     int //金币的ID
}

func (this *CenterManager) AddGoldToCenterControl(user IUser, addGoldType, addGold int, centerkey string) error {
	// bagmd := user.GetBag(bag.BagKey_Gold).(*bag.UserGoldBag)

	client := &http.Client{}
	DataUrlVal := url.Values{}
	DataUrlVal.Add("add_coin", util.NewStringAny(addGold).ToString())
	DataUrlVal.Add("description", centerkey)
	DataUrlVal.Add("ext1", util.NewStringAny(addGoldType).ToString())

	//todo
	req, err := http.NewRequest("POST", this.CenterSystemUrl+"/api/outer/coin/addCoin", strings.NewReader(DataUrlVal.Encode())) //bytes.NewReader(bytesData))
	if err != nil {
		log.Error("CenterControlManager AddGoldToCenterControl GetUserKey=%s,addGoldType=%d,addGold=%d http err=%v", user.GetOpenID(), addGoldType, addGold, err)
		// bagmd.ErrAddGold += addGold
		return err
	}
	user.GetCenter().SetReqHeader(req)
	//log.Info("----CenterControlManager AddGoldToCenterControl addGoldType=%d,addGold=%d,user.CenterControlInfo=%v,qid=%s", addGoldType, addGold, user.CenterControlInfo, user.CenterControlInfo.Qid)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("CenterControlManager AddGoldToCenterControl GetUserKey=%s,addGoldType=%d,addGold=%d clientdo err=%v", user.GetOpenID(), addGoldType, addGold, err)
		// bagmd.ErrAddGold += addGold
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("CenterControlManager AddGoldToCenterControl GetUserKey=%s,addGoldType=%d,addGold=%d readall err=%v", user.GetOpenID(), addGoldType, addGold, err)
		// bagmd.ErrAddGold += addGold
		return err
	}
	log.Info(" GetUserKey=%s,addGold string(body)=%s,addGoldType=%d, addGold=%d", user.GetOpenID(), string(body), addGoldType, addGold)

	var ccGoldInfo CCAck

	if err := json.Unmarshal(body, &ccGoldInfo); err == nil {
		if ccGoldInfo.Code != 0 {
			// bagmd.ErrAddGold += addGold
			return err
		}
		// if addGoldType == MAKEUP {
		// 	user.ErrAddGold = 0
		// } else {
		// 	// m.Gold.AddGoldDetail(user, addGoldType, addGold, true)
		// }
		return nil
	} else {
		// bagmd.ErrAddGold += addGold
		Logger.PError(err, "CC addGold=%d, user.OpenId=%s", addGold, user.GetOpenID())
		return err
	}

}

func (this *CenterManager) HandleNotifyDrawGold(w http.ResponseWriter, r *http.Request) {
	log.Info("=====HandleNotifyDrawGold===========r.PostFormValue(draw_coin)=%s,uid=%v", r.PostFormValue("draw_coin"), r.PostFormValue("uid"))
	drawGoldStr := r.PostFormValue("draw_coin")
	drawGold, err := strconv.Atoi(drawGoldStr)
	if err != nil {
		var format = `{"code":%d, "message":"%s"}`
		w.Write([]byte(fmt.Sprintf(format, -1, "gold atoi err")))
	}
	openId := r.PostFormValue("uid")
	item := new(redisinfo.RedisItem)
	item.LogType = GOLDLOG_CCFAIRBACK
	item.RedisKey = openId
	item.Reward = make(common.ItemInfos, 0, 1)
	item.Reward = append(item.Reward, &common.ItemInfo{
		ItemId: this.Item_GoldID,
		Count:  drawGold,
	})
	redisinfo.RedisUserRewardSet(item)

}

func (this *CenterManager) WithdrawGold(user IUser, cashType, withDrawType int) error {
	// log.Info("----CenterControlManager WithdrawGold cashType=%d,withDrawType=%d,user.CenterControlInfo=%v", cashType, withDrawType, user.CenterControlInfo)
	client := &http.Client{}
	data := make(map[string]string)
	data["cash_type"] = strconv.Itoa(cashType)
	data["withdraw_type"] = strconv.Itoa(withDrawType)
	DataUrlVal := url.Values{}
	for key, val := range data {
		DataUrlVal.Add(key, val)
	}
	req, err := http.NewRequest("POST", this.CenterSystemUrl+"/api/outer/withdraw/addWithdraw", strings.NewReader(DataUrlVal.Encode())) //bytes.NewReader(bytesData))

	if err != nil {
		log.Error("CenterControlManager WithdrawGold userId=%d,cashType=%d,withDrawType=%d http err=%v", user.GetOpenID(), cashType, withDrawType, err)
		return err
	}
	user.GetCenter().SetReqHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("CenterControlManager WithdrawGold userId=%d,cashType=%d,withDrawType=%d clientdo err=%v", user.GetOpenID(), cashType, withDrawType, err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("CenterControlManager WithdrawGold userId=%d,cashType=%d,withDrawType=%d readall err=%v", user.GetOpenID(), cashType, withDrawType, err)
		return err
	}
	log.Info("WithdrawGold string(body)=%s", string(body))

	//fmt.Println(string(body))
	var ccWithDrawInfo CCAck
	if err := json.Unmarshal(body, &ccWithDrawInfo); err == nil {
		log.Info("-------------ccWithDrawInfo=%v", ccWithDrawInfo)
		if ccWithDrawInfo.Code != 0 {
			return errorx.Center_ErrWithDraw
		}
		// drawGoldT := gameDb().GetDrawGold(cashType)
		// if drawGoldT == nil {
		// 	return base.ErrWithDrawType
		// }
		// m.UserManager.ReduceGold(user, drawGoldT.Gold, DRAE_GOLD)
		// m.Gold.AddGoldDetail(user, DRAE_GOLD, drawGoldT.Gold, true)
		return nil
	} else {
		log.Error("CC addGold err=%v", err)
		return err
	}
	//fmt.Println(string(body))
}

func (this *CenterManager) InviteWithdrawGold(user IUser, cashType int) error {
	// log.Info("----CenterControlManager WithdrawGold cashType=%d,withDrawType=%d,user.CenterControlInfo=%v", cashType, withDrawType, user.CenterControlInfo)
	client := &http.Client{}
	data := make(map[string]string)
	data["cash_type"] = strconv.Itoa(cashType)
	DataUrlVal := url.Values{}
	for key, val := range data {
		DataUrlVal.Add(key, val)
	}
	req, err := http.NewRequest("POST", this.CenterSystemUrl+"/api/outer/withdraw/addInviteWithdraw", strings.NewReader(DataUrlVal.Encode())) //bytes.NewReader(bytesData))

	if err != nil {
		log.Error("CenterControlManager InviteWithdraw userId=%d,cashType=%d,http err=%v", user.GetOpenID(), cashType, err)
		return err
	}
	user.GetCenter().SetReqHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("CenterControlManager InviteWithdraw userId=%d,cashType=%d, clientdo err=%v", user.GetOpenID(), cashType, err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("CenterControlManager InviteWithdraw userId=%d,cashType=%d, readall err=%v", user.GetOpenID(), cashType, err)
		return err
	}
	log.Info("InviteWithdraw string(body)=%s", string(body))

	//fmt.Println(string(body))
	var ccWithDrawInfo CCAck
	if err := json.Unmarshal(body, &ccWithDrawInfo); err == nil {
		log.Info("-------------ccInviteWithdraw=%v", ccWithDrawInfo)
		if ccWithDrawInfo.Code != 0 {
			return errorx.Center_ErrWithDraw
		}
		return nil
	} else {
		log.Error("CC ccInviteWithdraw json.Unmarshal err=%v", err)
		return err
	}
	//fmt.Println(string(body))
}
