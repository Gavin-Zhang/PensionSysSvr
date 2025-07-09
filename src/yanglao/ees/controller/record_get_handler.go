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

type Records struct {
	Count int
	Data  []*structure.Record
}

func GetRecordsHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	err := checkAll(w, r, constant.Power_EES)
	if err != nil {
		seelog.Error("EES GetRecordsHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		seelog.Error("EES GetRecordsHandler page err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		seelog.Error("EES GetRecordsHandler limit err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	condition := GetOrderConditionMap(r)
	ret, err := gonet.CallByName("EesMysqlSvr", "GetRecords", page, limit, condition)
	//ret, err := gonet.CallByName("MysqlSvr", "EES_GetRecords", page, limit, condition)
	if err != nil {
		seelog.Error("EES GetRecordsHandler call EesMysqlSvr function GetRecords err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var records Records
	err = goutils.ExpandResult(ret, &records)
	if err != nil {
		seelog.Error("EES GetRecordsHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = records.Count
	back.Data = make([]interface{}, 1)
	back.Data[0] = records.Data
	sendback(w, back)
}

func GetOrderConditionMap(r *http.Request) map[string]string {
	condition := make(map[string]string)
	if r.FormValue("name") != "" {
		condition["name__icontains"] = r.FormValue("name")
	}
	if r.FormValue("phone") != "" {
		condition["phone"] = r.FormValue("phone")
	}
	if r.FormValue("orderidx") != "" {
		condition["idx"] = r.FormValue("orderidx")
	}
	if r.FormValue("yearmonth") != "" {
		condition["begin_time__istartswith"] = r.FormValue("yearmonth")
	}
	seelog.Info("---> " + r.FormValue("yearmonthday"))
	if r.FormValue("yearmonthday") != "" {
		v := r.FormValue("yearmonthday")
		if len(v) > 2 {
			condition["begin_time__icontains"] = v[2:]
		}
	}
	return condition
}
