package db

import (
	"time"
	"yanglao/hcc/controller"
	"yanglao/hcc/structure"

	//"github.com/beego/beego/v2/client/orm"
	"yanglao/gonet/orm"

	"github.com/cihub/seelog"
)

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) GetConsumptionTypes() []*structure.ConsumptionType {
	var list []*structure.ConsumptionType
	_, err := p.o.QueryTable("consumption_type").All(&list)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetConsumptionTypes  err:", err)
		return make([]*structure.ConsumptionType, 0)
	}
	return list
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) GetPaymentTypes() []*structure.PaymentType {
	var list []*structure.PaymentType
	_, err := p.o.QueryTable("payment_type").All(&list)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetPaymentTypes  err:", err)
		return make([]*structure.PaymentType, 0)
	}
	return list
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) GetOrder(idx string) *structure.Order {
	order := new(structure.Order)
	err := p.o.QueryTable("order").Filter("idx", idx).One(order)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetOrder  err:", err)
		return nil
	}
	return order
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) GetOrders(page int, limit int, condition map[string]string) controller.Orders {
	var orders []*structure.Order
	qs := p.o.QueryTable("order")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).OrderBy("-created").All(&orders)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetOrders 1 err:", err)
		return controller.Orders{}
	}

	count, err := qs.Count()
	if err != nil {
		seelog.Error("HccMysqlSvr::GetOrders 2 err:", err)
		return controller.Orders{}
	}

	back := controller.Orders{
		Count: int(count),
		Data:  orders,
	}
	return back
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) AddOrder(order *structure.Order) bool {
	_, err := p.o.Insert(order)
	if err != nil {
		seelog.Error("HccMysqlSvr::AddOrder save order err: ", err)
		return false
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) AssignOrder(orderidx string, workers []structure.OrderAssign) string {
	err := p.o.Begin()
	if err != nil {
		seelog.Error("HccMysqlSvr::AssignOrder Begin error: ", err)
		return "程序错误"
	}
	order := structure.Order{Idx: orderidx}
	if err = p.o.Read(&order); err != nil {
		seelog.Error("HccMysqlSvr::AssignOrder Read error: ", err)
		p.o.Rollback()
		return "程序错误"
	}
	if order.OrderStatus == structure.ORDER_STATUS_WAIT_PAY || order.OrderStatus == structure.ORDER_STATUS_FINISHED {
		seelog.Error("HccMysqlSvr::AssignOrder status error")
		p.o.Rollback()
		return "状态已变更无法修改"
	}
	if order.AssignTime.IsZero() {
		order.AssignTime = time.Now()
		order.OrderStatus = structure.ORDER_STATUS_WAIT_SERVICE
		if _, err = p.o.Update(&order, "assign_time", "order_status"); err != nil {
			seelog.Error("HccMysqlSvr::AssignOrder order update error:", err)
			p.o.Rollback()
			return "状态已变更无法修改"
		}
	}

	_, err = p.o.Raw("update order_assign set status=? where order_idx=?", structure.ORDER_ASSIGN_CANCEL, orderidx).Exec()
	if err != nil {
		p.o.Rollback()
		seelog.Error("HccMysqlSvr::AssignOrder order_assign update  error:", err)
		return "程序错误"
	}

	for _, v := range workers {
		if _, err = p.o.Insert(&v); err != nil {
			p.o.Rollback()
			seelog.Error("HccMysqlSvr::AssignOrder order_assign insert error:", err)
			return "程序错误"
		}
	}
	p.o.Commit()
	return ""
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) GetOrderAssign(idx string) []orm.Params {
	var maps []orm.Params
	sql := "select `order_assign`.`phone`, `worker`.`idx`, `worker`.`name`, `worker`.`class` from `order_assign` join `worker` on `order_assign`.`order_idx`=? and `order_assign`.`worker_idx`=`worker`.`idx` and `order_assign`.`status`=?;"
	_, err := p.o.Raw(sql, idx, structure.ORDER_ASSIGN_SAVE).Values(&maps)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetOrderAssign  err:", err)
	}
	return maps
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) ServiceFinishOrder(order *structure.Order, evaluation *structure.OrderEvaluation) bool {
	err := p.o.Begin()
	if err != nil {
		seelog.Error("HccMysqlSvr::ServiceFinishOrder Begin error: ", err)
		p.o.Rollback()
		return false
	}
	if _, err = p.o.Insert(evaluation); err != nil {
		p.o.Rollback()
		seelog.Error("HccMysqlSvr::ServiceFinishOrder InsertOrUpdate error: ", err)
		return false
	}
	if _, err = p.o.Update(order, "order_status", "begin_time", "end_time"); err != nil {
		p.o.Rollback()
		seelog.Error("HccMysqlSvr::ServiceFinishOrder update order err: ", err)
		return false
	}
	p.o.Commit()
	return true
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) FinisOrder(order *structure.Order, updates ...string) bool {
	_, err := p.o.Update(order, updates...)
	if err != nil {
		seelog.Error("HccMysqlSvr::FinisOrder update order err: ", err)
		return false
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) GetOrderEvaluation(orderidx string) *structure.OrderEvaluation {
	evaluation := new(structure.OrderEvaluation)
	err := p.o.QueryTable("order_evaluation").Filter("order_idx", orderidx).One(evaluation)
	if err != nil {
		if err == orm.ErrNoRows {
			return evaluation
		}
		seelog.Error("HccMysqlSvr::GetOrderEvaluation err:", err)
		return nil
	}
	return evaluation
}

///////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) GetOrderFromTypeList() []*structure.OrderFromType {
	var list []*structure.OrderFromType

	_, err := p.o.QueryTable("order_from_type").All(&list)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetOrderFromTypeList  err:", err)
	}

	return list
}
