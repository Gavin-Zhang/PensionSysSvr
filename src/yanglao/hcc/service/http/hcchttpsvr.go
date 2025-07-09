package service

import (
	"net/http"
	"yanglao/gonet"
	"yanglao/hcc/controller"
	"yanglao/static"
)

type HccHttpSvr struct {
	gonet.ActorModel
}

func (p *HccHttpSvr) Init(args string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", controller.LoginHandler)
	mux.HandleFunc("/getclient", controller.GetClientHandler)
	mux.HandleFunc("/getclients", controller.GetClientsHandler)
	mux.HandleFunc("/register", controller.RegisterClientHandler)
	mux.HandleFunc("/updateclient", controller.UpdateClientHandler)
	mux.HandleFunc("/getserviceclass", controller.GetServiceClassHandler)
	mux.HandleFunc("/getservices", controller.GetServicesHandler)
	mux.HandleFunc("/getworkerclass", controller.GetWorkerClassHandler)
	mux.HandleFunc("/getworkers", controller.GetWorkersHandler)

	mux.HandleFunc("/addworker", controller.AddWorkerHandler)

	mux.HandleFunc("/getconsumptiontypes", controller.GetConsumptionTypeHandler)
	mux.HandleFunc("/getorders", controller.GetOrdersHandler)
	mux.HandleFunc("/addorder", controller.AddOrderHandler)
	mux.HandleFunc("/getorderworkers", controller.GetOrderWorkersHandler)
	mux.HandleFunc("/assignorder", controller.AssignOrderHandler)
	mux.HandleFunc("/getpaymenttype", controller.GetPaymentTypeHandler)
	mux.HandleFunc("/orderserviced", controller.OrderServiceFinishHandler)
	mux.HandleFunc("/orderpayment", controller.OrderPaymentHandler)
	mux.HandleFunc("/orderevaluation", controller.GetOrderEvaluationHandler)

	mux.HandleFunc("/orderformtype", controller.GetOrderFromTypeHandler)

	mux.HandleFunc("/gethousekeepers", controller.GetHouseKeepersHandler)

	mux.HandleFunc("/uploadavatar", controller.UploadAvatarHandler)
	mux.HandleFunc("/uploadphoto", controller.UpdataPhotoHandler)
	mux.HandleFunc("/getphotos", controller.GetPhotosHandler)
	mux.HandleFunc("/deletephoto", controller.DeletePhotosHandler)

	mux.HandleFunc("/getsubsidytime", controller.GetSubsidyTime)

	mux.Handle("/avatarphoto/", http.StripPrefix("/avatarphoto/", http.FileServer(http.Dir("image/hcc/avatar"))))
	mux.Handle("/imagephoto/", http.StripPrefix("/imagephoto/", http.FileServer(http.Dir("image/hcc/photo"))))

	http.ListenAndServe(static.HttpConfig.HCCPort, mux)
}

func init() {
	gonet.RegService(&HccHttpSvr{})
}
