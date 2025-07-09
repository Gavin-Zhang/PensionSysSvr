package service

import (
	"net/http"
	. "yanglao/ees/controller"
	"yanglao/gonet"
	"yanglao/static"
)

type EesHttpSvr struct {
	gonet.ActorModel
}

func (p *EesHttpSvr) Init(args string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/getclient", GetClientHandler)
	mux.HandleFunc("/getclients", GetClientsHandler)
	mux.HandleFunc("/updateclient", UpdateClientHandler)
	mux.HandleFunc("/register", RegisterClientHandler)
	//	mux.HandleFunc("/getserviceclass", GetServiceClassHandler)
	//	mux.HandleFunc("/getservices", GetServicesHandler)
	//	mux.HandleFunc("/getworkerclass", GetWorkerClassHandler)
	mux.HandleFunc("/getworkers", GetWorkersHandler)

	//	mux.HandleFunc("/getconsumptiontypes", GetConsumptionTypeHandler)
	mux.HandleFunc("/getrecords", GetRecordsHandler)
	mux.HandleFunc("/addrecord", AddRecordHandler)
	mux.HandleFunc("/getrecordworkers", GetRecordWorkersHandler)
	//	mux.HandleFunc("/getorderworkers", GetOrderWorkersHandler)
	//	mux.HandleFunc("/assignorder", AssignOrderHandler)
	//	mux.HandleFunc("/getpaymenttype", GetPaymentTypeHandler)
	//	mux.HandleFunc("/orderserviced", OrderServiceFinishHandler)
	//	mux.HandleFunc("/orderpayment", OrderPaymentHandler)
	//	mux.HandleFunc("/orderevaluation", GetOrderEvaluationHandler)

	mux.HandleFunc("/uploadavatar", UploadAvatarHandler)
	mux.HandleFunc("/uploadphoto", UpdataPhotoHandler)
	mux.HandleFunc("/getphotos", GetPhotosHandler)
	mux.HandleFunc("/deletephoto", DeletePhotosHandler)

	mux.Handle("/avatarphoto/", http.StripPrefix("/avatarphoto/", http.FileServer(http.Dir("image/ees/avatar"))))
	mux.Handle("/imagephoto/", http.StripPrefix("/imagephoto/", http.FileServer(http.Dir("image/ees/photo"))))

	http.ListenAndServe(static.HttpConfig.ESSPort, mux)
}

func init() {
	gonet.RegService(&EesHttpSvr{})
}
