package controller

import (
	"crypto/md5"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"
	"utils"

	"yanglao/constant"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	filecors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("UploadAvatarHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	mform := r.MultipartForm
	idxs, ok := mform.Value["idx"]
	if !ok {
		seelog.Error("UploadAvatarHandler not find idx in param")
		sendErr(w, constant.ResponseCode_ParamErr, "缺少参数")
		return
	}
	if len(idxs) == 0 {
		seelog.Error("UploadAvatarHandler idx array len == 0")
		sendErr(w, constant.ResponseCode_ParamErr, "参数错误")
		return
	}

	_, ok = mform.File["file"]
	if !ok {
		seelog.Error("UploadAvatarHandler not find file in param")
		sendErr(w, constant.ResponseCode_ParamErr, "缺少参数")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	seelog.Info("file:", file)
	if err != nil {
		seelog.Error("UploadAvatarHandler inovke FormFile err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	defer file.Close()

	ext := path.Ext(fileHeader.Filename)
	if ext != ".jpg" {
		seelog.Error("UploadAvatarHandler image format error")
		sendErr(w, constant.ResponseCode_ParamErr, "头像格式错误")
		return
	}
	_, err = jpeg.Decode(file)
	if err != nil {
		seelog.Error("UploadAvatarHandler Picture verification failed")
		sendErr(w, constant.ResponseCode_ParamErr, "文件校验失败")
		return
	}
	file.Seek(0, 0)

	outname := fmt.Sprintf("%x", md5.Sum([]byte(idxs[0])))
	ret, err := gonet.CallByName("mysqlsvr", "SetAvatar", idxs[0], outname)
	if err != nil {
		seelog.Error("UploadAvatarHandler call mysqlsvr function SetAvatar err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	result := false
	utils.CheckError(goutils.ExpandResult(ret, &result))
	if !result {
		seelog.Error("UploadAvatarHandler mysqlsvr function SetAvatar update fail")
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	localfilename := "avatar/" + outname + ".jpg"
	out, err := os.Create(localfilename)
	if err != nil {
		seelog.Error("UploadAvatarHandler create file error:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		seelog.Error("UploadAvatarHandler copy file error:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = 0
	sendback(w, back)
}
