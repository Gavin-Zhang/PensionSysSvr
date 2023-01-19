package db

import (
	"sync"
	"utils"

	"yanglao/ees/structure"
	"yanglao/static"

	"github.com/orm"

	"yanglao/gonet"
)

type EesMysqlSvr struct {
	gonet.ActorModel

	o orm.Ormer
}

func (p *EesMysqlSvr) Init(wg *sync.WaitGroup) {
	p.RegName("EesMysqlSvr")
	utils.OutputInfo(p.GetName(), p.connect())
	wg.Done()
}

func (p *EesMysqlSvr) connect() error {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DR_MySQL)

	orm.RegisterModel(new(structure.Indexs),
		new(structure.Client),
		new(structure.Record),
		new(structure.RecordWorker),
		new(structure.Worker))

	orm.RegisterDataBase("default", "mysql", static.Db.EESConn)

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		return err
	}

	p.o = orm.NewOrm()
	p.init_index()

	return nil
}

func (p *EesMysqlSvr) init_index() {
	indexs := &structure.Indexs{}

	for _, key := range static.Db.EESIndexs {
		indexs.Key = key

		err := p.o.Read(indexs, "key")
		if err != nil {
			if err == orm.ErrNoRows {
				indexs.Value = 1
				p.o.Insert(indexs)
				continue
			}
			utils.CheckError(err)
		}
	}
}

func init() {
	gonet.RegService(&EesMysqlSvr{})
}
