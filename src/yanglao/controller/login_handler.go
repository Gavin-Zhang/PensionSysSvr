package controller

import (
	"net/http"

	"yanglao/base"
	"yanglao/constant"
	"yanglao/single"
)

//type LoginHandler struct {
//}

type BackPlayer struct {
	Account  string
	UserName string
	Phone    string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	account := r.FormValue("userName")
	passwd := /*base.MD5(*/ r.FormValue("pswd") /*)*/

	base.Cors(&w, r)

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
			Phone:    player.Phone}

		b.Code = constant.ResponseCode_Success
		b.Count = 1
		b.Data = make([]interface{}, b.Count)
		b.Data[0] = backplayer
	}

	sendback(w, b)
}
