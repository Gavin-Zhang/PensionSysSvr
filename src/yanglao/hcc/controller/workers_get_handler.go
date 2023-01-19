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

type Workers struct {
	Count int
	Data  []*structure.Worker
}

func GetWorkersHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Cors(&w, r)

	player := single.PlayerMgr.GetByRequest(r)
	if player == nil {
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		seelog.Error("GetWorkersHandler not find player by cookie")
		return
	}
	single.SessionMgr.SetCookie(w, player.Session)

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		seelog.Error("GetWorkersHandler page err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("GetWorkersHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	condition := GetWorkerConditionMap(r)
	ret, err := gonet.CallByName("HccMysqlSvr", "GetWorkers", page, limit, condition)
	if err != nil {
		seelog.Error("GetWorkersHandler call HccMysqlSvr function GetWorkers err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var workers Workers
	err = goutils.ExpandResult(ret, &workers)
	if err != nil {
		seelog.Error("GetWorkersHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = workers.Count
	back.Data = make([]interface{}, 1)
	back.Data[0] = workers.Data
	sendback(w, back)
}

func GetWorkerConditionMap(r *http.Request) map[string]string {
	condition := make(map[string]string)
	if r.FormValue("service") != "" {
		condition["service"] = r.FormValue("service")
	}
	if r.FormValue("class") != "" {
		condition["class"] = r.FormValue("class")
	}
	return condition
}
