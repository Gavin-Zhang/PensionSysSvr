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

func AssignOrderHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("AssignOrderHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	order := structure.Order{
		Idx:         r.FormValue("idx"),
		Waiter:      r.FormValue("waiter"),
		WaiterPhone: r.FormValue("waiterphone"),
		OrderStatus: structure.ORDER_STATUS_DISPENSE,
		AssignTime:  time.Now()}

	ret, err := gonet.CallByName("mysqlsvr", "AssignOrder", &order)
	if err != nil {
		seelog.Error("AssignOrderHandler call mysqlsvr function AssignOrder err:", err)
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
