package base

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"yanglao/static"

	"github.com/cihub/seelog"
)

func init_allow_origin(w *http.ResponseWriter, r *http.Request) {
	for _, v := range static.HttpConfig.Info {
		if v.AllowOrigin == r.Header.Get("Origin") {
			seelog.Info("r:Origin", r.Header.Get("Origin"), " set allow origion:", v.AllowOrigin)
			(*w).Header().Set("Access-Control-Allow-Origin", v.AllowOrigin)
			return
		}
	}
}

func Cors(w *http.ResponseWriter, r *http.Request) {
	//(*w).Header().Set("Access-Control-Allow-Origin", static.HttpConfig.AllowOrigin /*"http://localhost"*/)
	init_allow_origin(w, r)
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Control-Allow-Methods, Access-Control-Allow-Origin") //header的类型
	(*w).Header().Add("Access-Control-Expose-Headers", "*")
	(*w).Header().Add("Access-Control-Allow-Credentials", "true")
	r.Header.Set("Content-Type", "application/json;charset=utf-8")
}

func CorsWithoutCredentials(w *http.ResponseWriter, r *http.Request) {
	init_allow_origin(w, r)
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Control-Allow-Methods, Access-Control-Allow-Origin") //header的类型
	(*w).Header().Add("Access-Control-Expose-Headers", "*")
	r.Header.Set("Content-Type", "application/json;charset=utf-8")
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
