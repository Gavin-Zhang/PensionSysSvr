package controller

import (
	"encoding/json"
	"net/http"
	"time"
	"utils"
	"yanglao/gonet"

	"yanglao/constant"
	"yanglao/ees/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	err := checkAll(w, r, constant.Power_EES)
	if err != nil {
		seelog.Error("EES RegisterClientHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	client := structure.Client{ChinaId: r.FormValue("chinaid")}
	client.Name = r.FormValue("name")
	client.Phone = r.FormValue("phone")
	client.Addr = r.FormValue("addr")
	client.Healthy = r.FormValue("healthdescription")
	client.Remarks = r.FormValue("other")
	client.Handler = r.FormValue("registrant")
	temp, _ := json.Marshal(r.FormValue("contacts"))
	client.Contacts = string(temp)
	temp, _ = json.Marshal(r.FormValue("slow"))
	client.SlowIll = string(temp)
	client.RegisterTime, err = time.ParseInLocation("2006-01-02", r.FormValue("registrant_time"), time.Local)
	if err != nil {
		seelog.Error("EES RegisterClientHandler time.ParseInLocation err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "注册时间格式错误")
		return
	}

	ret, err := gonet.CallByName("EesMysqlSvr", "RegisterClient", &client)
	if err != nil {
		seelog.Error("EES RegisterClientHandler call EesMysqlSvr function RegisterClinet err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	result := ""
	utils.CheckError(goutils.ExpandResult(ret, &result))

	b := BackInfo{Code: constant.ResponseCode_Success}
	if result != "" {
		b.Code = constant.ResponseCode_Fail
		b.Meesage = result
	}
	sendback(w, b)
}
