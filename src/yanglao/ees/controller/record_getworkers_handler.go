package controller

import (
	"net/http"
	//"yanglao/utils"

	"yanglao/constant"
	"yanglao/gonet"
	//"yanglao/static"
	"yanglao/ees/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func GetRecordWorkersHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	err := checkAll(w, r, constant.Power_EES)
	if err != nil {
		seelog.Error("EES GetRecordsHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"orderidx"}) {
		seelog.Error("EES GetRecordWorkersHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	ret, err := gonet.CallByName("EesMysqlSvr", "GetRecordWorks", r.FormValue("orderidx"))
	//ret, err := gonet.CallByName("MysqlSvr", "EES_GetRecordWorks", r.FormValue("orderidx"))
	if err != nil {
		seelog.Error("EES GetRecordWorkersHandler call EesMysqlSvr function GetRecordWorks err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var workers []*structure.RecordWorker
	err = goutils.ExpandResult(ret, &workers)
	if err != nil {
		seelog.Error("EES GetRecordWorkersHandler ExpandResult err:", err)
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
