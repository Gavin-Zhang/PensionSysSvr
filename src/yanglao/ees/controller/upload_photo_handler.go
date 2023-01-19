package controller

import (
	"crypto/md5"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"

	"yanglao/constant"

	"github.com/cihub/seelog"
)

func UpdataPhotoHandler(w http.ResponseWriter, r *http.Request) {
	filecors(&w, r)

	err := checkAll(w, r, constant.Power_EES)
	if err != nil {
		seelog.Error("EES UpdataPhotoHandler err:", err)
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	mform := r.MultipartForm
	idxs, ok := mform.Value["idx"]
	if !ok {
		seelog.Error("EES UpdataPhotoHandler not find idx in param")
		sendErr(w, constant.ResponseCode_ParamErr, "缺少参数")
		return
	}
	if len(idxs) == 0 {
		seelog.Error("UpdataPhotoHandler idx array len == 0")
		sendErr(w, constant.ResponseCode_ParamErr, "参数错误")
		return
	}
	clientidxs, ok := mform.Value["clientidx"]
	if !ok {
		seelog.Error("EES UpdataPhotoHandler not find clientidx in param")
		sendErr(w, constant.ResponseCode_ParamErr, "缺少参数")
		return
	}
	if len(clientidxs) == 0 {
		seelog.Error("EES UpdataPhotoHandler clientidx array len == 0")
		sendErr(w, constant.ResponseCode_ParamErr, "参数错误")
		return
	}

	_, ok = mform.File["file"]
	if !ok {
		seelog.Error("EES UpdataPhotoHandler not find file in param")
		sendErr(w, constant.ResponseCode_ParamErr, "缺少参数")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	seelog.Info("file:", file)
	if err != nil {
		seelog.Error("EES UpdataPhotoHandler inovke FormFile err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	defer file.Close()

	ext := path.Ext(fileHeader.Filename)
	if ext != ".jpg" {
		seelog.Error("EES UpdataPhotoHandler image format error")
		sendErr(w, constant.ResponseCode_ParamErr, "照片格式错误")
		return
	}
	_, err = jpeg.Decode(file)
	if err != nil {
		seelog.Error("EES UpdataPhotoHandler Picture verification failed")
		sendErr(w, constant.ResponseCode_ParamErr, "文件校验失败")
		return
	}
	file.Seek(0, 0)

	dir := fmt.Sprintf("image/ees/photo/%x/%s", md5.Sum([]byte(clientidxs[0])), idxs[0])
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			seelog.Error("EES UpdataPhotoHandler create dir error:", err)
			sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
			return
		}
	}

	localfilename := dir + "/" + fileHeader.Filename
	out, err := os.Create(localfilename)
	if err != nil {
		seelog.Error("EES UpdataPhotoHandler create file error:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		seelog.Error("EES UpdataPhotoHandler copy file error:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	var back BackInfo
	back.Code = constant.ResponseCode_Success
	back.Count = 0
	sendback(w, back)
}
