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
	p.NewService("Mysqlsvr", &wg)
	p.NewService("Httpsvr_cli", "")
	//p.NewService("Httpsvr_cli", "")

}

func init() {
	gonet.RegService(&Mainsvr{})
}
