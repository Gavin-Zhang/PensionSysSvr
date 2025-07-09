package db

import (
	//"github.com/beego/beego/v2/client/orm"
	"yanglao/gonet/orm"

	"github.com/cihub/seelog"
)

func (p *HccMysqlSvr) GetPersonalSubsidy(clientidx string, date string) int {
	sql := "select sum(`duration`) as `subsidy_time` from `subsidy_record` where `client_idx`=? and `date` like ?;"

	var maps []orm.Params
	_, err := p.o.Raw(sql, clientidx, date).Values(&maps)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetPersonalSubsidy  err:", err)
	}

	value, ok := maps[0]["subsidy_time"]
	if !ok || value == nil {
		return 0
	}
	return 0
}
