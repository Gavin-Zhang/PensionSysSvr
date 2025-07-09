package service

import (
	"sync"
	"yanglao/gonet"
)

type Mainsvr struct {
	gonet.ActorModel
}

func (p *Mainsvr) Init(args string) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	p.NewService("HccMysqlSvr", &wg)
	wg.Add(1)
	p.NewService("EesMysqlSvr", &wg)
	wg.Add(1)
	p.NewService("JlyMysqlSvr", &wg)
	//	wg.Add(1)
	//	p.NewService("StoreMysqlSvr", &wg)

	//wg.Add(1)
	//p.NewService("MysqlSvr", &wg)

	p.NewService("HccHttpSvr", "")
	p.NewService("EesHttpSvr", "")
	p.NewService("JlyHttpSvr", "")
	//p.NewService("StoreHttpSvr", "")

}

func init() {
	gonet.RegService(&Mainsvr{})
}
