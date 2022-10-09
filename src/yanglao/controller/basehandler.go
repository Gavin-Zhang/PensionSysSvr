package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
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
