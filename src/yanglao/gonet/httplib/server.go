package httplib

import (
	"gonet"
	"gonet/utils"
	"io"
	"log"
	"net/http"
)

type httpHandler struct {
	handle   uint32
	funcName string
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()

		args := map[string]string{}

		for k, v := range r.Form {
			args[k] = v[0]
		}

		resultSet, err := gonet.Call(h.handle, h.funcName, args)
		if err != nil {
			log.Print("ERROR")
			return
		}

		var result string
		utils.ExpandResult(resultSet, &result)

		io.WriteString(w, result)
	}
}

func GetFunc(pattern string, handle uint32, funcName string) {
	handler := &httpHandler{}
	handler.handle = handle
	handler.funcName = funcName
	http.Handle(pattern, handler)
}

func httpServeLoop(addr string) {
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return
	}

	log.Println("httpServeLoop stop.")
}

func Serve(addr string) error {
	go httpServeLoop(addr)
	return nil
}
