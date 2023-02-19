package db

import (
	"yanglao/hcc/structure"

	"github.com/cihub/seelog"
)

func (p *HccMysqlSvr) GetHouseKeepers(page int, limit int, condition map[string]string) []*structure.HouseKeeper {
	var keepers []*structure.HouseKeeper
	qs := p.o.QueryTable("house_keeper")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&keepers)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetHouseKeepers 1 err:", err)
	}

	return keepers
}
