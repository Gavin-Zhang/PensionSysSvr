package controller

import (
	"net/http"
	"strconv"

	//"yanglao/base"
	"yanglao/constant"
	"yanglao/ees/structure"
	//"yanglao/single"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

type Workers struct {
	Count int
	Data  []*structure.Worker
}

func GetWorkersHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if err := checkAll(w, r, constant.Power_EES); err != nil {
		seelog.Error("EES GetWorkersHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		seelog.Error("EES GetWorkersHandler page err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("EES GetWorkersHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	condition := GetWorkerConditionMap(r)
	ret, err := gonet.CallByName("EesMysqlSvr", "GetWorkers", page, limit, condition)
	//ret, err := gonet.CallByName("MysqlSvr", "EES_GetWorkers", page, limit, condition)
	if err != nil {
		seelog.Error("EES GetWorkersHandler call EesMysqlSvr function GetWorkers err:", err)
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
	if r.FormValue("name") != "" {
		condition["name__icontains"] = r.FormValue("name")
	}
	if r.FormValue("phone") != "" {
		condition["phone"] = r.FormValue("phone")
	}
	if r.FormValue("class") != "" {
		condition["class"] = r.FormValue("class")
	}
	return condition
}
