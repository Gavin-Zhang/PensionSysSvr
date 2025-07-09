package controller

import (
	"fmt"
	"net/http"
	"strconv"
	//	"time"
	//	"yanglao/utils"

	"yanglao/constant"
	"yanglao/gonet"
	"yanglao/hcc/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

func AssignOrderHandler(w http.ResponseWriter, r *http.Request) {
	cors(&w, r)

	if checkSession(w, r) == nil {
		seelog.Error("AssignOrderHandler not find player by cookie")
		sendErr(w, constant.ResponseCode_CookieErr, "cookie error")
		return
	}

	if !checkNotEmptyParams(r, []string{"orderidx", "workerscount", "handler"}) {
		seelog.Error("AssignOrderHandler checkNotEmptyParams fail")
		sendErr(w, constant.ResponseCode_ParamErr, "信息不全")
		return
	}

	count, err := strconv.ParseInt(r.FormValue("workerscount"), 10, 16)
	if err != nil {
		seelog.Error("AssignOrderHandler strconv.ParseInt error:", err)
		sendErr(w, constant.ResponseCode_ParamErr, "参数错误")
		return
	}

	workers := make([]structure.OrderAssign, 0)
	for i := int64(0); i < count; i++ {
		seelog.Info(fmt.Sprintf("workers[%d][idx]", i))
		seelog.Info("workers[0][class]", r.FormValue("workers[0][class]"))
		if r.FormValue(fmt.Sprintf("workers[%d][idx]", i)) == "" {
			seelog.Error("AssignOrderHandler workers error")
			sendErr(w, constant.ResponseCode_ParamErr, "参数错误")
			return
		}

		worker := structure.OrderAssign{
			OrderIdx:  r.FormValue("orderidx"),
			WorkerIdx: r.FormValue(fmt.Sprintf("workers[%d][idx]", i)),
			Phone:     r.FormValue(fmt.Sprintf("workers[%d][phone]", i)),
			Status:    structure.ORDER_ASSIGN_SAVE,
			Handler:   r.FormValue("handler"),
		}
		workers = append(workers, worker)
	}

	ret, err := gonet.CallByName("HccMysqlSvr", "AssignOrder", r.FormValue("orderidx"), workers)
	if err != nil {
		seelog.Error("AssignOrderHandler call HccMysqlSvr function AssignOrder err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}

	result := ""
	if err = goutils.ExpandResult(ret, &result); err != nil {
		seelog.Error("AssignOrderHandler ExpandResult err:", err)
		sendErr(w, constant.ResponseCode_ProgramErr, "内部程序错误")
		return
	}
	if result != "" {
		seelog.Error("AssignOrderHandler mysql back:", result)
		sendErr(w, constant.ResponseCode_ProgramErr, result)
		return
	}

	b := BackInfo{Code: constant.ResponseCode_Success}
	sendback(w, b)
}
