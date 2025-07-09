package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"yanglao/base"
	"yanglao/constant"
	"yanglao/single"

	"github.com/cihub/seelog"
)

type BackInfo struct {
	Code    int           `json:"code"`
	Message string        `json:"msg"`
	Count   int           `json:"count"`
	Data    []interface{} `json:"data"`
}

func sendback(w http.ResponseWriter, info BackInfo) {
	m, _ := json.Marshal(info)
	fmt.Println(string(m))
	w.Write(m)
}

func sendErr(w http.ResponseWriter, code int, msg string) {
	var backinfo BackInfo
	backinfo.Code = code
	backinfo.Message = msg
	sendback(w, backinfo)
}

func cors(w *http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.Cors(w, r)
}

func corsWithoutCredentials(w *http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	base.CorsWithoutCredentials(w, r)
}

func filecors(w *http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(128)
	base.Cors(w, r)
}

func checkSession(w http.ResponseWriter, r *http.Request) *single.Player {
	player := single.PlayerMgr.GetByRequest(r)
	if player == nil {
		return nil
	}
	// 刷新SESSION
	single.SessionMgr.SetCookie(w, player.Session)
	return player
}

func checkAll(w http.ResponseWriter, r *http.Request, power uint32) error {
	player := single.PlayerMgr.GetByRequest(r)
	if player == nil {
		return errors.New("not found player")
	}

	if player.Power != constant.Power_ALL {
		if player.Power != power {
			return errors.New("no power")
		}
	}

	// 刷新SESSION
	single.SessionMgr.SetCookie(w, player.Session)
	return nil
}

func checkNotEmptyParams(r *http.Request, params []string) bool {
	for _, param := range params {
		if r.FormValue(param) == "" {
			if vs := r.Form[param]; len(vs) == 0 {
				seelog.Error("HCC checkNotEmptyParams not found param:", param)
				return false
			}
		}
	}
	return true
}
