package controller

import (
	"net/http"
	"strconv"

	"yanglao/base"
	"yanglao/constant"
	"yanglao/single"
	"yanglao/structure"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

type Clients struct {
	Count int
	Data  []*structure.Client
}

func GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Cors(&w, r)

	player := single.PlayerMgr.GetByRequest(r)
	if player == nil {
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		seelog.Error("GetClientsHandler not find player by cookie")
		return
	}
	single.SessionMgr.SetCookie(w, player.Session)

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		seelog.Error("GetClientsHandler page err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("GetClientsHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
	}

	condition := GetConditionMap(r)
	ret, err := gonet.CallByName("mysqlsvr", "GetClients", page, limit, condition)
	if err != nil {
		seelog.Error("GetClientsHandler call mysqlsvr function GetClients err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var clients Clients
	err = goutils.ExpandResult(ret, &clients)
	if err != nil {
		seelog.Error("GetClientsHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = clients.Count
	back.Data = make([]interface{}, 1)
	back.Data[0] = clients.Data
	sendback(w, back)
}

func GetConditionMap(r *http.Request) map[string]string {
	condition := make(map[string]string)
	if r.FormValue("name") != "" {
		condition["name"] = r.FormValue("name")
	}
	if r.FormValue("phone") != "" {
		condition["phone"] = r.FormValue("phone")
	}
	if r.FormValue("chinaid") != "" {
		condition["chinaid"] = r.FormValue("chinaid")
	}
	seelog.Info("GetConditionMap:", r.FormValue("name")+"|"+r.FormValue("phone")+"|"+r.FormValue("chinaid"))
	return condition
}
