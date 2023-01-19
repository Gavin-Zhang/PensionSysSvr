package db

import (
	"yanglao/ees/controller"
	"yanglao/ees/structure"

	"github.com/cihub/seelog"
)

///////////////////////////////////////////////////////////////////////////////
func (p *EesMysqlSvr) GetRecord(idx string) *structure.Record {
	Record := new(structure.Record)
	err := p.o.QueryTable("record").Filter("idx", idx).One(Record)
	if err != nil {
		seelog.Error("EesMysqlSvr::GetRecord  err:", err)
		return nil
	}
	return Record
}

///////////////////////////////////////////////////////////////////////////////
func (p *EesMysqlSvr) GetRecords(page int, limit int, condition map[string]string) controller.Records {
	var Records []*structure.Record
	qs := p.o.QueryTable("Record")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).OrderBy("-created").All(&Records)
	if err != nil {
		seelog.Error("EesMysqlSvr::GetRecords 1 err:", err)
		return controller.Records{}
	}

	count, err := qs.Count()
	if err != nil {
		seelog.Error("EesMysqlSvr::GetRecords 2 err:", err)
		return controller.Records{}
	}

	back := controller.Records{
		Count: int(count),
		Data:  Records,
	}
	return back
}

///////////////////////////////////////////////////////////////////////////////
func (p *EesMysqlSvr) AddRecord(Record *structure.Record, workers []structure.RecordWorker) bool {
	err := p.o.Begin()
	if err != nil {
		seelog.Error("EesMysqlSvr::AddRecord Begin error: ", err)
		return false
	}

	if _, err := p.o.Insert(Record); err != nil {
		p.o.Rollback()
		seelog.Error("EesMysqlSvr::AddRecord save Record err: ", err)
		return false
	}

	for _, v := range workers {
		if _, err = p.o.Insert(&v); err != nil {
			p.o.Rollback()
			seelog.Error("EesMysqlSvr::AddRecord save RecordWorker err: ", err)
			return false
		}
	}

	p.o.Commit()
	return true
}

///////////////////////////////////////////////////////////////////////////////
func (p *EesMysqlSvr) GetRecordWorks(idx string) []*structure.RecordWorker {
	var workers []*structure.RecordWorker

	_, err := p.o.QueryTable("record_worker").Filter("record_idx", idx).All(&workers)
	if err != nil {
		seelog.Error("EesMysqlSvr::GetRecordWorks  err:", err)
	}

	return workers
}
