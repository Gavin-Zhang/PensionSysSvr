package controller

import (
	"crypto/md5"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"
	"yanglao/utils"

	"yanglao/constant"

	"yanglao/gonet"
	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	filecors(&w, r)

	err := checkAll(w, r, constant.Power_EES)
	if err != nil {
		seelog.Error("EES UploadAvatarHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	mform := r.MultipartForm
	idxs, ok := mform.Value["idx"]
	if !ok {
		seelog.Error("EES UploadAvatarHandler not find idx in param")
		sendErr(w, constant.ResponseCode_ParamErr, "缺少参数")
		return
	}
	if len(idxs) == 0 {
		seelog.Error("EES UploadAvatarHandler idx array len == 0")
		sendErr(w, constant.ResponseCode_ParamErr, "参数错误")
		return
	}

	_, ok = mform.File["file"]
	if !ok {
		seelog.Error("EES UploadAvatarHandler not find file in param")
		sendErr(w, constant.ResponseCode_ParamErr, "缺少参数")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	seelog.Info("file:", file)
	if err != nil {
		seelog.Error("EES UploadAvatarHandler inovke FormFile err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	defer file.Close()

	ext := path.Ext(fileHeader.Filename)
	if ext != ".jpg" {
		seelog.Error("EES UploadAvatarHandler image format error")
		sendErr(w, constant.ResponseCode_ParamErr, "头像格式错误")
		return
	}
	_, err = jpeg.Decode(file)
	if err != nil {
		seelog.Error("EES UploadAvatarHandler Picture verification failed")
		sendErr(w, constant.ResponseCode_ParamErr, "文件校验失败")
		return
	}
	file.Seek(0, 0)

	outname := fmt.Sprintf("%x", md5.Sum([]byte(idxs[0])))
	ret, err := gonet.CallByName("EesMysqlSvr", "SetAvatar", idxs[0], outname)
	//ret, err := gonet.CallByName("MysqlSvr", "EES_SetAvatar", idxs[0], outname)
	if err != nil {
		seelog.Error("EES UploadAvatarHandler call EesMysqlSvr function SetAvatar err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	result := false
	utils.CheckError(goutils.ExpandResult(ret, &result))
	if !result {
		seelog.Error("EES UploadAvatarHandler EesMysqlSvr function SetAvatar update fail")
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	dir := "image/ees/avatar"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			seelog.Error("EES UpdataPhotoHandler create dir error:", err)
			sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
			return
		}
	}

	localfilename := dir + "/" + outname + ".jpg"
	out, err := os.Create(localfilename)
	if err != nil {
		seelog.Error("EES UploadAvatarHandler create file error:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		seelog.Error("EES UploadAvatarHandler copy file error:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = 0
	sendback(w, back)
}
