package controller

import (
	"net/http"

	"yanglao/constant"

	"yanglao/hcc/structure"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func GetOrderEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("GetOrderEvaluationHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"orderidx"}) {
		seelog.Error("GetOrderEvaluationHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	ret, err := gonet.CallByName("HccMysqlSvr", "GetOrderEvaluation", r.FormValue("orderidx"))
	if err != nil {
		seelog.Error("GetOrderEvaluationHandler call HccMysqlSvr function GetOrderEvaluation err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	evaluation := new(structure.OrderEvaluation)
	err = goutils.ExpandResult(ret, &evaluation)
	if err != nil || evaluation == nil {
		seelog.Error("GetClientsHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = 1
	back.Data = make([]interface{}, 1)
	back.Data[0] = *evaluation
	sendback(w, back)
}
