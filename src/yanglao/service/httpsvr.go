package service

import (
	//"fmt"
	"net/http"
	"yanglao/controller"
	"yanglao/gonet"
)

type Httpsvr struct {
	gonet.ActorModel
}

func (p *Httpsvr) Init(args string) {
	http.HandleFunc("/login", controller.LoginHandler)
	http.HandleFunc("/getclients", controller.GetClientsHandler)
	http.HandleFunc("/register", controller.RegisterClientHandler)
	http.HandleFunc("/getorders", controller.GetOrdersHandler)

	http.ListenAndServe("0.0.0.0:8001", nil)
}

func init() {
	gonet.RegService(&Httpsvr{})
}
