package db

import (
	"yanglao/controller"
	"yanglao/structure"

	"github.com/cihub/seelog"
)

func (p *Mysqlsvr) GetConsumptionTypes() []*structure.ConsumptionType {
	var list []*structure.ConsumptionType
	_, err := p.o.QueryTable("consumption_type").All(&list)
	if err != nil {
		seelog.Error("Mysqlsvr::GetConsumptionTypes  err:", err)
		return make([]*structure.ConsumptionType, 0)
	}
	return list
}

func (p *Mysqlsvr) GetPaymentTypes() []*structure.PaymentType {
	var list []*structure.PaymentType
	_, err := p.o.QueryTable("payment_type").All(&list)
	if err != nil {
		seelog.Error("Mysqlsvr::GetPaymentTypes  err:", err)
		return make([]*structure.PaymentType, 0)
	}
	return list
}

func (p *Mysqlsvr) GetOrders(page int, limit int, condition map[string]string) controller.Orders {
	var orders []*structure.Order
	qs := p.o.QueryTable("order")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&orders)
	if err != nil {
		seelog.Error("Mysqlsvr::GetOrders 1 err:", err)
		return controller.Orders{}
	}

	count, err := qs.Count()
	if err != nil {
		seelog.Error("Mysqlsvr::GetOrders 2 err:", err)
		return controller.Orders{}
	}

	back := controller.Orders{
		Count: int(count),
		Data:  orders,
	}
	return back
}

func (p *Mysqlsvr) AddOrder(order *structure.Order) bool {
	_, err := p.o.Insert(order)
	if err != nil {
		seelog.Error("Mysqlsvr::AddOrder save order err: ", err)
		return false
	}
	return true
}

func (p *Mysqlsvr) AssignOrder(order *structure.Order) bool {
	_, err := p.o.Update(order, "waiter", "waiter_phone", "assign_time", "order_status")
	if err != nil {
		seelog.Error("Mysqlsvr::AssignOrder update order err: ", err)
		return false
	}
	return true
}

func (p *Mysqlsvr) ServiceFinishOrder(order *structure.Order) bool {
	_, err := p.o.Update(order, "order_status", "begin_time", "end_time")
	if err != nil {
		seelog.Error("Mysqlsvr::ServiceFinishOrder update order err: ", err)
		return false
	}
	return true
}

func (p *Mysqlsvr) FinisOrder(order *structure.Order, updates ...string) bool {
	_, err := p.o.Update(order, updates...)
	if err != nil {
		seelog.Error("Mysqlsvr::FinisOrder update order err: ", err)
		return false
	}
	return true
}
