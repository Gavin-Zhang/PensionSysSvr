package controller

import (
	"net/http"

	"yanglao/constant"

	"github.com/cihub/seelog"
)

func GetSubsidyTime(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("GetSubsidyTime not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"ClientIdx"}) {
		seelog.Error("GetSubsidyTime checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}
}
