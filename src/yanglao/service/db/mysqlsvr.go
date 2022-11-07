package db

import (
	"sync"
	"utils"

	"yanglao/static"
	"yanglao/structure"

	"github.com/cihub/seelog"

	"yanglao/gonet"
	"yanglao/gonet/orm"
)

type Mysqlsvr struct {
	gonet.ActorModel

	o orm.Ormer
}

func (p *Mysqlsvr) Init(wg *sync.WaitGroup) {
	p.RegName("mysqlsvr")
	utils.OutputInfo(p.GetName(), p.connect())
	wg.Done()
}

func (p *Mysqlsvr) connect() error {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DR_MySQL)

	orm.RegisterModel(new(structure.Indexs),
		new(structure.User),
		new(structure.Client),
		new(structure.Order),
		new(structure.Service),
		new(structure.ServiceClass),
		new(structure.Worker),
		new(structure.WorkerClass),
		new(structure.ConsumptionType),
		new(structure.PaymentType),
		new(structure.OrderAssign),
		new(structure.OrderEvaluation))

	orm.RegisterDataBase("default", "mysql", static.Db.Conn)

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		return err
	}

	p.o = orm.NewOrm()
	p.init_index()

	return nil
}

func (p *Mysqlsvr) init_index() {
	indexs := &structure.Indexs{}

	for _, key := range static.Db.Indexs {
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

func (p *Mysqlsvr) Login(account string, password string) *structure.User {
	user := new(structure.User)
	p.o.QueryTable(user).Filter("account", account).Filter("pass_word", password).One(user)
	seelog.Info("Login:", account, user)
	return user
}

func init() {
	gonet.RegService(&Mysqlsvr{})
}
