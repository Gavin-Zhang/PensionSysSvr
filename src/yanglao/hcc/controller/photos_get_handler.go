package controller

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"yanglao/constant"
	"yanglao/hcc/structure"

	"github.com/cihub/seelog"
)

func GetPhotosHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("GetPhotosHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	params := []string{"idx", "clientidx"}
	if !checkNotEmptyParams(r, params) {
		seelog.Error("GetPhotosHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	dir := "./image/hcc/photo/"
	urlpath := fmt.Sprintf("%x/%s", md5.Sum([]byte(r.FormValue("clientidx"))), r.FormValue("idx"))
	dir += urlpath
	_, err := os.Stat(dir)
	if err != nil && !os.IsNotExist(err) {
		seelog.Error("GetPhotosHandler 获取文件夹异常 err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	} else if os.IsNotExist(err) {
		var back BackInfo
		back.Code = constant.ResponseCode_Success
		back.Count = 0
		back.Data = make([]interface{}, 1)
		sendback(w, back)
		return
	}

	photos := make([]structure.OrderPhoto, 0)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".jpg" {
			return nil
		}
		photo := structure.OrderPhoto{
			Name:       info.Name(),
			Path:       urlpath + "/" + info.Name(),
			Size:       info.Size(),
			CreateTime: info.ModTime()}
		photos = append(photos, photo)
		return nil
	})
	if err != nil {
		seelog.Error("GetPhotosHandler filepath.Walk err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = len(photos)
	back.Data = make([]interface{}, 1)
	back.Data[0] = photos
	sendback(w, back)
}
