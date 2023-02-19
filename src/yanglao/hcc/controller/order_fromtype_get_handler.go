package controller

import (
	"net/http"

	"yanglao/base"
	"yanglao/constant"
	"yanglao/hcc/structure"
	"yanglao/single"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func GetOrderFromTypeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Cors(&w, r)

	player := single.PlayerMgr.GetByRequest(r)
	if player == nil {
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		seelog.Error("GetOrderFromTypeHandler not find player by cookie")
		return
	}
	single.SessionMgr.SetCookie(w, player.Session)

	ret, err := gonet.CallByName("HccMysqlSvr", "GetOrderFormTypeList")
	if err != nil {
		seelog.Error("GetOrderFromTypeHandler call HccMysqlSvr function GetOrderFormTypeList err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var list []*structure.OrderFromType
	err = goutils.ExpandResult(ret, &list)
	if err != nil {
		seelog.Error("GetOrderFromTypeHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = len(list)
	back.Data = make([]interface{}, 1)
	back.Data[0] = list
	sendback(w, back)
}
