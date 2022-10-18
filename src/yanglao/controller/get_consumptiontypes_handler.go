package controller

import (
	"net/http"

	"yanglao/constant"
	"yanglao/structure"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func GetConsumptionTypesHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("GetConsumptionTypesHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	ret, err := gonet.CallByName("mysqlsvr", "GetConsumptionTypes")
	if err != nil {
		seelog.Error("GetConsumptionTypesHandler call mysqlsvr function GetConsumptionTypes err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var list []*structure.ConsumptionType
	err = goutils.ExpandResult(ret, &list)
	if err != nil {
		seelog.Error("GetConsumptionTypesHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = len(list)
	back.Data = make([]interface{}, 1)
	back.Data[0] = list
	sendback(w, back)
}
