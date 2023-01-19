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

	p.NewService("HccHttpSvr", "")
	p.NewService("EesHttpSvr", "")
	//p.NewService("Httpsvr_svr", "")

}

func init() {
	gonet.RegService(&Mainsvr{})
}
