package controller

import (
	"net/http"
	"strconv"

	"yanglao/constant"
	"yanglao/gonet"
	"yanglao/hcc/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func GetHouseKeepersHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("GetHouseKeeperHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"page", "limit"}) {
		seelog.Error("GetHouseKeepersHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		seelog.Error("GetHouseKeepersHandler page err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("GetHouseKeepersHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	ret, err := gonet.CallByName("HccMysqlSvr", "GetHouseKeepers", page, limit, GetHouseKeeperConditionMap(r))
	if err != nil {
		seelog.Error("GetHouseKeepersHandler call HccMysqlSvr function GetHouseKeepers err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var keepers []*structure.HouseKeeper
	err = goutils.ExpandResult(ret, &keepers)
	if err != nil {
		seelog.Error("GetHouseKeepersHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = len(keepers)
	back.Data = make([]interface{}, 1)
	back.Data[0] = keepers
	sendback(w, back)
}

func GetHouseKeeperConditionMap(r *http.Request) map[string]string {
	condition := make(map[string]string)
	if r.FormValue("service") != "" {
		condition["service"] = r.FormValue("service")
	}
	if r.FormValue("class") != "" {
		condition["class"] = r.FormValue("class")
	}
	return condition
}
