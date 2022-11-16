package controller

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"

	"yanglao/constant"

	"github.com/cihub/seelog"
)

func DeletePhotosHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("DeletePhotosHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	params := []string{"idx", "clientidx", "name"}
	if !checkNotEmptyParams(r, params) {
		seelog.Error("DeletePhotosHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	dir := "./photo/" + fmt.Sprintf("%x/%s", md5.Sum([]byte(r.FormValue("clientidx"))), r.FormValue("idx"))
	file := dir + "/" + r.FormValue("name")

	if err := os.Remove(file); err != nil {
		seelog.Error("DeletePhotosHandler os.Remove error:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "删除失败")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = 0
	sendback(w, back)
}
