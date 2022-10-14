package db

import (
	"yanglao/controller"
	"yanglao/structure"

	"github.com/cihub/seelog"
)

func (p *Mysqlsvr) GetWorkerClassList() []*structure.WorkerClass {
	var list []*structure.WorkerClass

	_, err := p.o.QueryTable("worker_class").All(&list)
	if err != nil {
		seelog.Error("Mysqlsvr::GetWorkerClassList  err:", err)
		return make([]*structure.WorkerClass, 0)
	}

	return list
}

func (p *Mysqlsvr) GetWorkers(page int, limit int, condition map[string]string) controller.Workers {
	var workers []*structure.Worker
	qs := p.o.QueryTable("worker")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&workers)
	if err != nil {
		seelog.Error("Mysqlsvr::GetWorker 1 err:", err)
		return controller.Workers{}
	}

	count, err := qs.Count()
	if err != nil {
		seelog.Error("Mysqlsvr::GetWorker 2 err:", err)
		return controller.Workers{}
	}

	back := controller.Workers{
		Count: int(count),
		Data:  workers,
	}
	return back
}
