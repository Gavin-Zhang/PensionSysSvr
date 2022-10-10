package db

import (
	"strings"
	"time"

	"yanglao/controller"
	"yanglao/static"
	"yanglao/structure"

	"github.com/cihub/seelog"
)

func (p *Mysqlsvr) GetOrders(page int, limit int) controller.Orders {
	var orders []*structure.Order
	_, err := p.o.QueryTable("order").Limit(limit, (page-1)*limit).All(&orders)
	if err != nil {
		seelog.Error("Mysqlsvr::GetOrders 1 err:", err)
	}

	count := 0
	err = p.o.Raw("select count(*) from order").QueryRow(&count)
	if err != nil {
		seelog.Error("Mysqlsvr::GetOrders 2 err:", err)
	}

	back := controller.Orders{
		Count: count,
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
