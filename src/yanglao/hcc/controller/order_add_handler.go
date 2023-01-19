package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"utils"

	"yanglao/constant"
	"yanglao/gonet"
	"yanglao/static"
	"yanglao/hcc/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {

	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("AddOrderHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	order := structure.Order{
		Idx:             static.Db.OrderHead + strings.Replace(time.Now().Format("20060102150405.0000000"), ".", "", -1),
		ClientIdx:       r.FormValue("clientidx"),
		Name:            r.FormValue("name"),
		Phone:           r.FormValue("phone"),
		Community:       r.FormValue("ascription"),
		Addr:            r.FormValue("addr"),
		Service:         r.FormValue("service"),
		Remarks:         r.FormValue("remarks"),
		Handler:         r.FormValue("handler"),
		ConsumptionType: r.FormValue("consumptiontype"),
		ChinaId:         r.FormValue("chinaid"),
		OrderStatus:     structure.ORDER_STATUS_WAIT_ASSIGN,
		PaymentStatus:   structure.ORDER_PAY_STATUS_WAIT}

	value, err := time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("servicetime"), time.Local)
	if err != nil {
		seelog.Error("AddOrderHandler time.ParseInLocation err:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "注册时间格式错误")
		return
	}
	order.ServiceTime = value

	order.Fare = 0
	if r.FormValue("fare") != "" {
		value, err := strconv.ParseInt(r.FormValue("fare"), 10, 16)
		if err != nil {
			seelog.Error("AddOrderHandler strconv.ParseInt(fare) err:", err)
			sendErr(w, constant.ResponseCode_ParamErr, "车费不合法")
			return
		}
		order.Fare = int16(value)
	}

	order.HighRise = 0.0
	if r.FormValue("highrise") != "" {
		value, err := strconv.ParseFloat(r.FormValue("highrise"), 32)
		if err != nil {
			seelog.Error("AddOrderHandler strconv.ParseFloat(highrise) err:", err)
			sendErr(w, constant.ResponseCode_ParamErr, "高层费用不合法")
			return
		}
		order.HighRise = float32(value)
	}

	ret, err := gonet.CallByName("HccMysqlSvr", "AddOrder", &order)
	if err != nil {
		seelog.Error("AddOrderHandler call HccMysqlSvr function AddOrder err:", err)
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
