package controller

import (
	"net/http"
	"strconv"

	"yanglao/base"
	"yanglao/constant"
	"yanglao/single"
	"yanglao/hcc/structure"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

type Services struct {
	Count int
	Data  []*structure.Service
}

func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Cors(&w, r)

	player := single.PlayerMgr.GetByRequest(r)
	if player == nil {
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		seelog.Error("GetServicesHandler not find player by cookie")
		return
	}
	single.SessionMgr.SetCookie(w, player.Session)

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		seelog.Error("GetServicesHandler page err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("GetServicesHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	condition := GetServiceConditionMap(r)
	ret, err := gonet.CallByName("HccMysqlSvr", "GetServices", page, limit, condition)
	if err != nil {
		seelog.Error("GetServicesHandler call HccMysqlSvr function GetServices err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var services Services
	err = goutils.ExpandResult(ret, &services)
	if err != nil {
		seelog.Error("GetServicesHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = services.Count
	back.Data = make([]interface{}, 1)
	back.Data[0] = services.Data
	sendback(w, back)
}

func GetServiceConditionMap(r *http.Request) map[string]string {
	condition := make(map[string]string)
	if r.FormValue("service") != "" {
		condition["service"] = r.FormValue("service")
	}
	if r.FormValue("class") != "" {
		condition["class"] = r.FormValue("class")
	}
	return condition
}
