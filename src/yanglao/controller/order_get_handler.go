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

type Orders struct {
	Count int
	Data  []*structure.Order
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("GetClientsHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	condition := GetOrderConditionMap(r)
	ret, err := gonet.CallByName("mysqlsvr", "GetOrders", page, limit, condition)
	if err != nil {
		seelog.Error("GetClientsHandler call mysqlsvr function GetClients err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var orders Orders
	err = goutils.ExpandResult(ret, &orders)
	if err != nil {
		seelog.Error("GetClientsHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = orders.Count
	back.Data = make([]interface{}, 1)
	back.Data[0] = orders.Data
	sendback(w, back)
}

func GetOrderConditionMap(r *http.Request) map[string]string {
	condition := make(map[string]string)
	if r.FormValue("name") != "" {
		condition["name"] = r.FormValue("name")
	}
	if r.FormValue("phone") != "" {
		condition["phone"] = r.FormValue("phone")
	}
	if r.FormValue("orderidx") != "" {
		condition["idx"] = r.FormValue("orderidx")
	}
	return condition
}
