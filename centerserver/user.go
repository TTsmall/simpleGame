package centerserver

import (
	"net/http"

	"github.com/TTsmall/wawaji_pub_hmhelper/log"
)

//用户数据
type UserCenterMD struct {
	Token    string
	Platform string
	Ver      string
	Udi      string
	Uid      string
	AppId    string
	Qid      string
}

func (md *UserCenterMD) SetReqHeader(req *http.Request) {
	log.Info("UserCenterMD SetReqHeader md=%v", md)
	req.Header.Add("Authorization", "Bearer "+md.Token)
	req.Header.Add("platform", md.Platform)
	req.Header.Add("ver", md.Ver)
	req.Header.Add("udi", md.Udi)
	req.Header.Add("uid", md.Uid)
	req.Header.Add("qid", md.Qid)
	req.Header.Add("zm-app-id", md.AppId)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
}

//用户接口
type IUser interface {
	//用户的key
	GetOpenID() string
	//拿用户的中台信息
	GetCenter() *UserCenterMD
}
