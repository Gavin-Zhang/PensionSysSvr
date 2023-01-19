package db

import (
	"yanglao/ees/controller"
	"yanglao/ees/structure"

	"github.com/cihub/seelog"
)

func (p *EesMysqlSvr) GetWorkerClassList() []*structure.WorkerClass {
	var list []*structure.WorkerClass

	_, err := p.o.QueryTable("worker_class").All(&list)
	if err != nil {
		seelog.Error("EesMysqlSvr::GetWorkerClassList  err:", err)
		return make([]*structure.WorkerClass, 0)
	}

	return list
}

func (p *EesMysqlSvr) GetWorkers(page int, limit int, condition map[string]string) controller.Workers {
	var workers []*structure.Worker
	qs := p.o.QueryTable("worker")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&workers)
	if err != nil {
		seelog.Error("EesMysqlSvr::GetWorker 1 err:", err)
		return controller.Workers{}
	}

	count, err := qs.Count()
	if err != nil {
		seelog.Error("EesMysqlSvr::GetWorker 2 err:", err)
		return controller.Workers{}
	}

	back := controller.Workers{
		Count: int(count),
		Data:  workers,
	}
	return back
}
