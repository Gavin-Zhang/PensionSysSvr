package service

import (
	//"fmt"
	"gonet"
	"net/http"
	"yanglao/controller"
)

type Httpsvr struct {
	gonet.ActorModel
}

func (p *Httpsvr) Init(args string) {
	http.HandleFunc("/login", controller.LoginHandler)
	http.HandleFunc("/getclients", controller.GetClientsHandler)
	http.HandleFunc("/register", controller.RegisterClientHandler)

	http.ListenAndServe("0.0.0.0:8001", nil)
}

func init() {
	gonet.RegService(&Httpsvr{})
}
