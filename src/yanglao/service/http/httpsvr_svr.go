package service

import (
	"net/http"
	"yanglao/controller"
	"yanglao/gonet"
)

type Httpsvr_svr struct {
	gonet.ActorModel
}

func (p *Httpsvr_svr) Init(args string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", controller.LoginHandler)
	mux.HandleFunc("/getclients", controller.GetClientsHandler)
	mux.HandleFunc("/register", controller.RegisterClientHandler)
	mux.HandleFunc("/getorders", controller.GetOrdersHandler)

	http.ListenAndServe("0.0.0.0:8002", mux)
}

func init() {
	gonet.RegService(&Httpsvr_svr{})
}
