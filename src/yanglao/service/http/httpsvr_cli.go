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
	mux.HandleFunc("/getclient", controller.GetClientHandler)
	mux.HandleFunc("/getclients", controller.GetClientsHandler)
	mux.HandleFunc("/register", controller.RegisterClientHandler)
	mux.HandleFunc("/getserviceclass", controller.GetServiceClassHandler)
	mux.HandleFunc("/getservices", controller.GetServicesHandler)
	mux.HandleFunc("/getworkerclass", controller.GetWorkerClassHandler)
	mux.HandleFunc("/getworkers", controller.GetWorkersHandler)

	mux.HandleFunc("/getconsumptiontypes", controller.GetConsumptionTypeHandler)
	mux.HandleFunc("/getorders", controller.GetOrdersHandler)
	mux.HandleFunc("/addorder", controller.AddOrderHandler)
	mux.HandleFunc("/getorderworkers", controller.GetOrderWorkersHandler)
	mux.HandleFunc("/assignorder", controller.AssignOrderHandler)
	mux.HandleFunc("/getpaymenttype", controller.GetPaymentTypeHandler)
	mux.HandleFunc("/orderserviced", controller.OrderServiceFinishHandler)
	mux.HandleFunc("/orderpayment", controller.OrderPaymentHandler)
	mux.HandleFunc("/orderevaluation", controller.GetOrderEvaluationHandler)

	mux.HandleFunc("/uploadavatar", controller.UploadAvatarHandler)
	mux.HandleFunc("/uploadphoto", controller.UpdataPhotoHandler)

	mux.Handle("/avatarphoto/", http.StripPrefix("/avatarphoto/", http.FileServer(http.Dir("avatar"))))
	mux.Handle("/imagephoto/", http.StripPrefix("/imagephoto/", http.FileServer(http.Dir("photo"))))

	http.ListenAndServe("0.0.0.0:8001", mux)
}

func init() {
	gonet.RegService(&Httpsvr_cli{})
}
