package controller

import (
	"net/http"
	//"yanglao/utils"

	"yanglao/constant"
	"yanglao/gonet"
	//"yanglao/static"
	//"yanglao/hcc/structure"

	goutils "yanglao/gonet/utils"

	//"github.com/beego/beego/v2/client/orm"
	"yanglao/gonet/orm"

	"github.com/cihub/seelog"
)

func GetOrderWorkersHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("GetOrderWorkersHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"orderidx"}) {
		seelog.Error("GetOrderWorkersHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	ret, err := gonet.CallByName("HccMysqlSvr", "GetOrderAssign", r.FormValue("orderidx"))
	if err != nil {
		seelog.Error("GetOrderWorkersHandler call HccMysqlSvr function GetOrderAssign err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var workers []orm.Params
	err = goutils.ExpandResult(ret, &workers)
	if err != nil {
		seelog.Error("GetOrderWorkersHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = len(workers)
	back.Data = make([]interface{}, 1)
	back.Data[0] = workers
	sendback(w, back)
}
