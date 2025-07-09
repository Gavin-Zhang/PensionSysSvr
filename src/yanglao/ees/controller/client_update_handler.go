package controller

import (
	"encoding/json"
	"net/http"
	"yanglao/utils"

	"yanglao/constant"
	"yanglao/ees/structure"
	"yanglao/gonet"

	"github.com/cihub/seelog"

	goutils "yanglao/gonet/utils"
)

func UpdateClientHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	err := checkAll(w, r, constant.Power_EES)
	if err != nil {
		seelog.Error("EES UpdateClientHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"idx", "chinaid", "name", "phone", "addr",
		"healthdescription", "other", "contacts", "slow", "changecid"}) {
		seelog.Error("EES UpdateClientHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	client := structure.Client{
		Idx:     r.FormValue("idx"),
		ChinaId: r.FormValue("chinaid"),
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		Addr:    r.FormValue("addr"),
		Healthy: r.FormValue("healthdescription"),
		Remarks: r.FormValue("other")}

	temp, _ := json.Marshal(r.FormValue("contacts"))
	client.Contacts = string(temp)
	temp, _ = json.Marshal(r.FormValue("slow"))
	client.SlowIll = string(temp)

	ret, err := gonet.CallByName("EesMysqlSvr", "UpdateClient", &client, r.FormValue("changecid") == "true")
	//ret, err := gonet.CallByName("MysqlSvr", "EES_UpdateClient", &client, r.FormValue("changecid") == "true")
	if err != nil {
		seelog.Error("UpdateClientHandler call EesMysqlSvr function UpdateClient err:", err)
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
