package controller

import (
	"net/http"
	"time"
	"utils"

	"yanglao/constant"
	"yanglao/gonet"
	"yanglao/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func OrderServiceFinishHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("OrderServiceFinishHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	params := []string{"idx", "servicebegin", "serviceend"}
	if !checkNotEmptyParams(r, params) {
		seelog.Error("OrderServiceFinishHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	order := structure.Order{
		Idx:         r.FormValue("idx"),
		OrderStatus: structure.ORDER_STATUS_SFINISHED}

	value, err := time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("servicebegin"), time.Local)
	if err != nil {
		seelog.Error("OrderServiceFinishHandler time.ParseInLocation servicebegin err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "服务时间格式错误")
		return
	}
	order.BeginTime = value
	value, err = time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("serviceend"), time.Local)
	if err != nil {
		seelog.Error("OrderServiceFinishHandler time.ParseInLocation serviceend err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "服务时间格式错误")
		return
	}
	order.EndTime = value

	ret, err := gonet.CallByName("mysqlsvr", "ServiceFinishOrder", order)
	if err != nil {
		seelog.Error("OrderServiceFinishHandler call mysqlsvr function ServiceFinishOrder err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	result := false
	utils.CheckError(goutils.ExpandResult(ret, &result))

	b := BackInfo{Code: constant.ResponseCode_Success}
	if !result {
		b.Code = constant.ResponseCode_Fail
	}
	sendback(w, b)
}
