package db

import (
	"fmt"
	"strings"

	"yanglao/constant"
	"yanglao/gonet/orm"
	"yanglao/hcc/controller"
	"yanglao/hcc/structure"
	"yanglao/static"

	"github.com/cihub/seelog"
)

////////////////////////////////////////////////////////////////////////////////
// Handle Client
////////////////////////////////////////////////////////////////////////////////
func (p *HccMysqlSvr) CheckChinaID(id string) uint16 {
	exist := p.o.QueryTable("client").Filter("china_id", strings.ToUpper(id)).Exist()
	if exist {
		return constant.RegisterResult_RepeatID
	}
	return constant.RegisterResult_Success
}

func (p *HccMysqlSvr) RegisterClient(client *structure.Client) bool {
	role_index := &structure.Indexs{Key: structure.Indexs_Role_Index}
	err := p.o.Read(role_index)
	if err != nil {
		seelog.Error("HccMysqlSvr::RegisterClient get role index err: ", err)
		return false
	}

	index := role_index.Value
	role_index.Value += 1
	_, err = p.o.Update(role_index)
	if err != nil {
		seelog.Error("HccMysqlSvr::RegisterClient update role index err: ", err)
		return false
	}

	client.Idx = fmt.Sprintf("%s%010d", static.Db.IndexHead, index)
	_, err = p.o.Insert(client)
	if err != nil {
		seelog.Error("HccMysqlSvr::RegisterClient save role err: ", err)
		return false
	}
	return true
}

func (p *HccMysqlSvr) GetClient(idx string) *structure.Client {
	client := new(structure.Client)
	if err := p.o.QueryTable("client").Filter("idx", idx).One(client); err != nil {
		seelog.Error("HccMysqlSvr::GetClient err: ", err)
		return nil
	}
	return client
}

func (p *HccMysqlSvr) GetClients(page int, limit int, condition map[string]string) controller.Clients {
	var clients []*structure.Client

	qs := p.o.QueryTable("client")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&clients)
	if err != nil {
		seelog.Error("HccMysqlSvr::GetClients 1 err:", err)
		return controller.Clients{}
	}

	//err = p.o.Raw("select count(*) from client").QueryRow(&count)
	count, err := qs.Count()
	if err != nil {
		seelog.Error("HccMysqlSvr::GetClients 2 err:", err)
		return controller.Clients{}
	}

	back := controller.Clients{
		Count: int(count),
		Data:  clients,
	}
	return back
}

func (p *HccMysqlSvr) SetAvatar(idx string, avatar string) bool {
	client := structure.Client{Idx: idx, Avatar: avatar}
	_, err := p.o.Update(&client, "avatar")
	if err != nil {
		seelog.Error("HccMysqlSvr::SetAvatar err: ", err)
		return false
	}
	return true
}

func (p *HccMysqlSvr) UpdateClient(client *structure.Client, change_chinaid bool) bool {
	_, err := p.o.Update(client, "china_id", "name", "phone", "addr",
		"type", "community", "healthy", "remarks", "contacts", "slow_ill")
	if err != nil {
		seelog.Error("HccMysqlSvr::UpdateClient err: ", err)
		return false
	}

	if change_chinaid {
		p.o.QueryTable("order").Filter("client_idx", client.Idx).Update(orm.Params{
			"china_id": client.ChinaId,
		})
	}

	return true
}
