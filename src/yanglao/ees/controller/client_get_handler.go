package controller

import (
	"net/http"
	"strconv"

	"yanglao/constant"
	"yanglao/ees/structure"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

type Clients struct {
	Count int
	Data  []*structure.Client
}

func GetClientHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("EES GetClientHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"idx"}) {
		seelog.Error("EES GetClientHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	ret, err := gonet.CallByName("EesMysqlSvr", "GetClient", r.FormValue("idx"))
	if err != nil {
		seelog.Error("EES GetClientHandler call EesMysqlSvr function GetClient err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	client := new(structure.Client)
	if err = goutils.ExpandResult(ret, &client); err != nil {
		seelog.Error("EES GetClientHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = 1
	back.Data = make([]interface{}, 1)
	back.Data[0] = client
	sendback(w, back)
}

func GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	seelog.Info("GetClientsHandler")
	cors(&w, r)

	if err := checkAll(w, r, constant.Power_EES); err != nil {
		seelog.Error("EES RegisterClientHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		seelog.Error("EES GetClientsHandler page err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("EES GetClientsHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
	}

	condition := GetClientConditionMap(r)
	ret, err := gonet.CallByName("EesMysqlSvr", "GetClients", page, limit, condition)
	if err != nil {
		seelog.Error("EES GetClientsHandler call EesMysqlSvr function GetClients err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var clients Clients
	err = goutils.ExpandResult(ret, &clients)
	if err != nil {
		seelog.Error("EES GetClientsHandler ExpandResult err:", err)
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

func GetClientConditionMap(r *http.Request) map[string]string {
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
	return condition
}
