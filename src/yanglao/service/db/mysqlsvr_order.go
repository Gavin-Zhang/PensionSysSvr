package db

import (
	"strings"
	"time"

	"yanglao/controller"
	"yanglao/static"
	"yanglao/structure"

	"github.com/cihub/seelog"
)

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
	order.Idx = static.Db.OrderHead + strings.Replace(time.Now().Format("20060102150405.000000000"), ".", "", -1)
	_, err := p.o.Insert(order)
	if err != nil {
		seelog.Error("Mysqlsvr::AddOrder save order err: ", err)
		return false
	}
	return true
}
