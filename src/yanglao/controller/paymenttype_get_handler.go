package controller

import (
	"net/http"

	"yanglao/constant"
	"yanglao/structure"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func GetPaymentTypeHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("GetPaymenttypeHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	ret, err := gonet.CallByName("mysqlsvr", "GetPaymentTypes")
	if err != nil {
		seelog.Error("GetPaymentTypeHandler call mysqlsvr function GetPaymentTypes err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var list []*structure.PaymentType
	err = goutils.ExpandResult(ret, &list)
	if err != nil {
		seelog.Error("GetPaymentTypeHandler ExpandResult err:", err)
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
