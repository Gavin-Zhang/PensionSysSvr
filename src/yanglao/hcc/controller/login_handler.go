package controller

import (
	"net/http"

	"yanglao/base"
	"yanglao/constant"
	"yanglao/single"

	"github.com/cihub/seelog"
)

//type LoginHandler struct {
//}

type BackPlayer struct {
	Account  string
	UserName string
	Phone    string
	Power    uint32
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	corsWithoutCredentials(&w, r)

	//r.ParseForm()
	account := r.FormValue("userName")
	passwd := /*base.MD5(*/ r.FormValue("pswd") /*)*/

	base.Cors(&w, r)
	seelog.Info("=======================================")
	seelog.Info("w:", w.Header())
	seelog.Info("=======================================")
	seelog.Info("r:", r.Header)
	seelog.Info("=======================================")
	seelog.Info("user:", account, "passwd:", passwd)
	seelog.Info("=======================================")
	player := single.PlayerMgr.Load(account, passwd)
	var b BackInfo
	b.Code = constant.ResponseCode_Fail

	if player != nil && player.Account == account {
		if player.Session != "" {
			single.SessionMgr.EndSession(player.Session)
		}

		session := single.SessionMgr.BeginSession(w, r)
		session.Set(single.SessionKey_Acc, account)
		player.Session = session.SessionID()

		backplayer := BackPlayer{
			Account:  player.Account,
			UserName: player.UserName,
			Phone:    player.Phone,
			Power:    player.Power}

		b.Code = constant.ResponseCode_Success
		b.Count = 1
		b.Data = make([]interface{}, b.Count)
		b.Data[0] = backplayer
	}

	sendback(w, b)
}
