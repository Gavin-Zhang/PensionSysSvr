package db

import (
	"yanglao/hcc/controller"
	"yanglao/hcc/structure"

	"github.com/cihub/seelog"
)

func (p *HccMysqlSvr) AddWorker(worker *structure.Worker) string {
	info := new(structure.Worker)
	err := p.o.QueryTable("worker").Filter("china_id", worker.ChinaId).Filter("class", worker.Class).Limit(1).One(info)
	if err == nil {
		return "重复注册"
	}

	_, err = p.o.Insert(worker)
	if err != nil {
		seelog.Error("HccMysqlSvr::AddWorker err: ", err)
		return err.Error()
	}
	return ""
}

func (p *HccMysqlSvr) GetWorkerClassList() []*structure.WorkerClass {
	var list []*structure.WorkerClass

	_, err := p.o.QueryTable("worker_class").All(&list)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetWorkerClassList  err:", err)
		return make([]*structure.WorkerClass, 0)
	}

	return list
}

func (p *HccMysqlSvr) GetWorkers(page int, limit int, condition map[string]string) controller.Workers {
	var workers []*structure.Worker
	qs := p.o.QueryTable("worker")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&workers)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetWorker 1 err:", err)
		return controller.Workers{}
	}

	count, err := qs.Count()
	if err != nil {
		seelog.Error("HccMysqlSvr::GetWorker 2 err:", err)
		return controller.Workers{}
	}

	back := controller.Workers{
		Count: int(count),
		Data:  workers,
	}
	return back
}
