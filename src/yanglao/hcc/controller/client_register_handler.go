package controller

import (
	"encoding/json"
	"net/http"
	"time"
	"utils"
	"yanglao/gonet"

	"yanglao/base"
	"yanglao/constant"
	"yanglao/hcc/structure"
	"yanglao/single"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Cors(&w, r)

	player := single.PlayerMgr.GetByRequest(r)
	if player == nil {
		seelog.Error("RegisterClientHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}
	single.SessionMgr.SetCookie(w, player.Session)

	client := structure.Client{ChinaId: r.FormValue("chinaid")}
	ret, err := gonet.CallByName("HccMysqlSvr", "CheckChinaID", client.ChinaId)
	if err != nil {
		seelog.Error("RegisterClientHandler call HccMysqlSvr function CheckChinaID err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	check_result := constant.RegisterResult_Success
	err = goutils.ExpandResult(ret, &check_result)
	if err != nil {
		seelog.Error("RegisterClientHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	if check_result != constant.RegisterResult_Success {
		b := BackInfo{
			Code:    constant.RegisterResult_RepeatID,
			Message: "已注册过!"}
		sendback(w, b)
		return
	}

	client.Name = r.FormValue("name")
	client.Phone = r.FormValue("phone")
	client.Addr = r.FormValue("addr")
	client.Community = r.FormValue("ascription")
	client.Type = r.FormValue("type")
	client.Healthy = r.FormValue("healthdescription")
	client.Remarks = r.FormValue("other")
	client.Handler = r.FormValue("registrant")
	client.InDbTime = time.Now()
	temp, _ := json.Marshal(r.FormValue("contacts"))
	client.Contacts = string(temp)
	temp, _ = json.Marshal(r.FormValue("slow"))
	client.SlowIll = string(temp)
	client.RegisterTime, err = time.ParseInLocation("2006-01-02", r.FormValue("registrant_time"), time.Local)
	if err != nil {
		seelog.Error("RegisterClientHandler time.ParseInLocation err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "注册时间格式错误")
		return
	}

	ret, err = gonet.CallByName("HccMysqlSvr", "RegisterClient", &client)
	if err != nil {
		seelog.Error("RegisterClientHandler call HccMysqlSvr function RegisterClinet err:", err)
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
