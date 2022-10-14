package service

import (
	"net/http"
	"yanglao/controller"
	"yanglao/gonet"
)

type Httpsvr_cli struct {
	gonet.ActorModel
}

func (p *Httpsvr_cli) Init(args string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", controller.LoginHandler)
	mux.HandleFunc("/getclients", controller.GetClientsHandler)
	mux.HandleFunc("/register", controller.RegisterClientHandler)
	mux.HandleFunc("/getorders", controller.GetOrdersHandler)
	mux.HandleFunc("/getserviceclass", controller.GetServiceClassHandler)
	mux.HandleFunc("/getservices", controller.GetServicesHandler)
	mux.HandleFunc("/getworkerclass", controller.GetWorkerClassHandler)
	mux.HandleFunc("/getworkers", controller.GetWorkersHandler)

	http.ListenAndServe("0.0.0.0:8001", mux)
}

func init() {
	gonet.RegService(&Httpsvr_cli{})
}
