package controller

import (
	"net/http"
	"strconv"
	"time"
	"utils"

	"yanglao/constant"
	"yanglao/gonet"
	"yanglao/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func FinishOrderHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("FinishOrderHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	params := []string{"idx", "consumptiontype", "servicebegin", "serviceend"}
	if !checkNotEmptyParams(r, params) {
		seelog.Error("FinishOrderHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	order := structure.Order{
		Idx:         r.FormValue("idx"),
		OrderStatus: structure.ORDER_STATUS_FINISHED,
		FinishTime:  time.Now()}

	value, err := time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("servicebegin"), time.Local)
	if err != nil {
		seelog.Error("FinishOrderHandler .ParseInLocation err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "服务时间格式错误")
		return
	}
	order.BeginTime = value
	value, err = time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("serviceend"), time.Local)
	if err != nil {
		seelog.Error("FinishOrderHandler time.ParseInLocation err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "服务时间格式错误")
		return
	}
	order.EndTime = value

	if !setCharge(&w, r, &order) {
		return
	}
	if !setPaymentStatus(&w, r, &order) {
		return
	}

	ret, err := DoFinish(&order, r.FormValue("consumptiontype"))
	if err != nil {
		seelog.Error("FinishOrderHandler call mysqlsvr function FinisOrder err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	result := false
	utils.CheckError(goutils.ExpandResult(ret, &result))

	b := BackInfo{Code: constant.ResponseCode_Success}
	if !result {
		b.Code = constant.ResponseCode_Fail
	}
	sendback(w, b)
}

func setCharge(w *http.ResponseWriter, r *http.Request, order *structure.Order) bool {
	if r.FormValue("consumptiontype") == "政府购买" {
		return true
	}

	params := []string{"charge", "fare", "highrise"}
	if !checkNotEmptyParams(r, params) {
		seelog.Error("FinishOrderHandler setCharge checkNotEmptyParams fail")
		sendErr(*w, constant.ResponseCode_ParamErr, "信息不全")
		return false
	}

	value, err := strconv.ParseInt(r.FormValue("charge"), 10, 16)
	if err != nil {
		seelog.Error("FinishOrderHandler strconv.ParseInt(charge) err:", err)
		sendErr(*w, constant.ResponseCode_ParamErr, "收费格式错误")
		return false
	}
	order.Charge = int16(value)

	value, err = strconv.ParseInt(r.FormValue("fare"), 10, 16)
	if err != nil {
		seelog.Error("FinishOrderHandler strconv.ParseInt(fare) err:", err)
		sendErr(*w, constant.ResponseCode_ParamErr, "车费格式错误")
		return false
	}
	order.Fare = int16(value)

	valuef, err := strconv.ParseFloat(r.FormValue("highrise"), 32)
	if err != nil {
		seelog.Error("FinishOrderHandler strconv.ParseFloat(highrise) err:", err)
		sendErr(*w, constant.ResponseCode_ParamErr, "高层费用格式错误")
		return false
	}
	order.HighRise = float32(valuef)
	return true
}

func setPaymentStatus(w *http.ResponseWriter, r *http.Request, order *structure.Order) bool {
	if r.FormValue("consumptiontype") == "政府购买" {
		order.PaymentStatus = structure.ORDER_PAY_STATUS_WAIT
		return true
	}

	params := []string{"paymenttype"}
	if !checkNotEmptyParams(r, params) {
		seelog.Error("FinishOrderHandler setPaymentStatus checkNotEmptyParams fail")
		sendErr(*w, constant.ResponseCode_ParamErr, "信息不全")
		return false
	}

	order.PaymentType = r.FormValue("paymenttype")
	order.PaymentStatus = structure.ORDER_PAY_STATUS_OVER
	return true
}

func DoFinish(order *structure.Order, consumptionType string) ([]interface{}, error) {
	if consumptionType == "政府购买" {
		return gonet.CallByName("mysqlsvr", "FinisOrder", order,
			"order_status", "begin_time", "end_time", "finish_time", "payment_status")
	}
	return gonet.CallByName("mysqlsvr", "FinisOrder", order,
		"order_status", "begin_time", "end_time", "finish_time", "payment_status",
		"payment_type", "charge", "fare", "high_rise")
}
