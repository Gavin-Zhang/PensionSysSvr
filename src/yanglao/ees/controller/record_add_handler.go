package controller

import (
	"encoding/json"
	"net/http"
	//"strconv"
	"strings"
	"time"
	"yanglao/utils"

	"yanglao/constant"
	"yanglao/ees/structure"
	"yanglao/gonet"
	"yanglao/static"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func AddRecordHandler(w http.ResponseWriter, r *http.Request) {

	cors(&w, r)

	if err := checkAll(w, r, constant.Power_EES); err != nil {
		seelog.Error("EES AddRecordHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	record := structure.Record{
		Idx:       static.Db.EESRecord + strings.Replace(time.Now().Format("20060102150405.0000000"), ".", "", -1),
		ClientIdx: r.FormValue("clientidx"),
		Name:      r.FormValue("name"),
		Phone:     r.FormValue("phone"),
		Addr:      r.FormValue("addr"),
		ChinaId:   r.FormValue("chinaid"),
		Service:   r.FormValue("service"),
		Remarks:   r.FormValue("remarks"),
		Handler:   r.FormValue("handler")}

	value, err := time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("servicebegin"), time.Local)
	if err != nil {
		seelog.Error("EES AddRecordHandler time.ParseInLocation servicebegin err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "服务时间格式错误")
		return
	}
	record.BeginTime = value
	value, err = time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("serviceend"), time.Local)
	if err != nil {
		seelog.Error("EES AddRecordHandler time.ParseInLocation serviceend err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "服务时间格式错误")
		return
	}
	record.EndTime = value

	workers := make([]structure.RecordWorker, 0)
	if err := json.Unmarshal([]byte(r.FormValue("servers")), &workers); err != nil {
		seelog.Error("EES AddRecordHandler worker unmarshal err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "解析服务人员数据错误")
		return
	}
	for k, _ := range workers {
		workers[k].RecordIdx = record.Idx
	}

	ret, err := gonet.CallByName("EesMysqlSvr", "AddRecord", &record, workers)
	//ret, err := gonet.CallByName("MysqlSvr", "EES_AddRecord", &record, workers)
	if err != nil {
		seelog.Error("EES AddRecordHandler call EesMysqlSvr function AddRecord err:", err)
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
