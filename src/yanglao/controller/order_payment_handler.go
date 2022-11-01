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

func OrderPaymentHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("OrderPaymentHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"idx", "consumptiontype", "paymenttime"}) {
		seelog.Error("OrderPaymentHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	order := structure.Order{
		Idx:           r.FormValue("idx"),
		PaymentStatus: structure.ORDER_PAY_STATUS_OVER,
		OrderStatus:   structure.ORDER_STATUS_FINISHED,
		FinishTime:    time.Now()}

	value, err := time.ParseInLocation("2006-01-02", r.FormValue("paymenttime"), time.Local)
	if err != nil {
		seelog.Error("OrderPaymentHandler time.ParseInLocation err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "支付日期格式错误")
		return
	}
	order.PaymentTime = value

	consumptiontype := r.FormValue("consumptiontype")
	if consumptiontype != "政府购买" {
		params := []string{"fare", "charge", "highrise", "paymenttype"}
		if !checkNotEmptyParams(r, params) {
			seelog.Error("OrderPaymentHandler checkNotEmptyParams fail")
			sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
			return
		}

		value, err := strconv.ParseInt(r.FormValue("charge"), 10, 16)
		if err != nil {
			seelog.Error("FinishOrderHandler strconv.ParseInt(charge) err:", err)
			sendErr(w, constant.ResponseCode_ParamErr, "收费格式错误")
			return
		}
		order.Charge = int16(value)

		value, err = strconv.ParseInt(r.FormValue("fare"), 10, 16)
		if err != nil {
			seelog.Error("FinishOrderHandler strconv.ParseInt(fare) err:", err)
			sendErr(w, constant.ResponseCode_ParamErr, "车费格式错误")
			return
		}
		order.Fare = int16(value)

		valuef, err := strconv.ParseFloat(r.FormValue("highrise"), 32)
		if err != nil {
			seelog.Error("FinishOrderHandler strconv.ParseFloat(highrise) err:", err)
			sendErr(w, constant.ResponseCode_ParamErr, "高层费用格式错误")
			return
		}
		order.HighRise = float32(valuef)
		order.PaymentType = r.FormValue("paymenttype")
	}

	ret, err := DoFinish(&order, consumptiontype)
	if err != nil {
		seelog.Error("OrderPaymentHandler call mysqlsvr function FinisOrder err:", err)
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

func DoFinish(order *structure.Order, consumptionType string) ([]interface{}, error) {
	if consumptionType == "政府购买" {
		return gonet.CallByName("mysqlsvr", "FinisOrder", order,
			"order_status", "finish_time", "payment_status", "payment_time")
	}
	return gonet.CallByName("mysqlsvr", "FinisOrder", order,
		"order_status", "finish_time", "payment_status",
		"payment_type", "charge", "fare", "high_rise")
}
