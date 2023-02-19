package subsidy

import (
	"utils"
)

type SubsidySvr struct {
	gonet.ActorModel

	o orm.Ormer
}

func (p *SubsidySvr) Init() {
	p.RegName("SubsidySvr")
	utils.OutputInfo(p.GetName(), p.connect())
}

func (p *SubsidySvr) GetSubsidyTime(clientidx string) {
	// TODO 获取个人月服务时间
	// TODO 获取全局月服务时间
	// TODO 获取预留服务时间
}

func init() {
	gonet.RegService(&SubsidySvr{})
}
