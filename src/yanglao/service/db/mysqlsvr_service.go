package db

import (
	"yanglao/controller"
	"yanglao/structure"

	"github.com/cihub/seelog"
)

func (p *Mysqlsvr) GetServiceClassList() []*structure.ServiceClass {
	var list []*structure.ServiceClass

	_, err := p.o.QueryTable("service_class").All(&list)
	if err != nil {
		seelog.Error("Mysqlsvr::GetServiceClassList  err:", err)
		return make([]*structure.ServiceClass, 0)
	}

	return list
}

func (p *Mysqlsvr) GetServices(page int, limit int, condition map[string]string) controller.Services {
	var services []*structure.Service
	qs := p.o.QueryTable("service")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&services)
	if err != nil {
		seelog.Error("Mysqlsvr::GetServices 1 err:", err)
		return controller.Services{}
	}

	count, err := qs.Count()
	if err != nil {
		seelog.Error("Mysqlsvr::GetServices 2 err:", err)
		return controller.Services{}
	}

	back := controller.Services{
		Count: int(count),
		Data:  services,
	}
	return back
}
