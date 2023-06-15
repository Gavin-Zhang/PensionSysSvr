package db

import (
	"fmt"
	"strings"

	"yanglao/constant"
	"yanglao/ees/controller"
	"yanglao/ees/structure"
	"yanglao/static"

	"github.com/orm"

	"github.com/cihub/seelog"
)

////////////////////////////////////////////////////////////////////////////////
// Handle Client
////////////////////////////////////////////////////////////////////////////////
//func (p *EesMysqlSvr) CheckChinaID(id string) uint16 {
//	exist := p.o.QueryTable("client").Filter("china_id", strings.ToUpper(id)).Exist()
//	if exist {
//		return constant.RegisterResult_RepeatID
//	}
//	return constant.RegisterResult_Success
//}

func (p *EesMysqlSvr) RegisterClient(client *structure.Client) string {
	if p.o.QueryTable("client").Filter("china_id", strings.ToUpper(client.ChinaId)).Exist() {
		return constant.Error_RepeatID.Error()
	}

	err := p.o.Begin()
	if err != nil {
		seelog.Error("EesMysqlSvr::RegisterClient Begin error: ", err)
		return constant.Error_Program.Error()
	}

	role_index := &structure.Indexs{Key: structure.Indexs_Role_Index}
	err = p.o.Read(role_index)
	if err != nil {
		seelog.Error("EesMysqlSvr::RegisterClient get role index err: ", err)
		return constant.Error_Program.Error()
	}

	index := role_index.Value
	role_index.Value += 1
	_, err = p.o.Update(role_index)
	if err != nil {
		p.o.Rollback()
		seelog.Error("EesMysqlSvr::RegisterClient update role index err: ", err)
		return constant.Error_Program.Error()
	}

	client.Idx = fmt.Sprintf("%s%010d", static.Db.EESIndexHead, index)
	_, err = p.o.Insert(client)
	if err != nil {
		p.o.Rollback()
		seelog.Error("EesMysqlSvr::RegisterClient save role err: ", err)
		return constant.Error_Program.Error()
	}

	p.o.Commit()
	return ""
}

func (p *EesMysqlSvr) GetClient(key string, value string) *structure.Client {
	client := new(structure.Client)
	if err := p.o.QueryTable("client").Filter(key, value).One(client); err != nil {
		seelog.Error("EesMysqlSvr::GetClient err: ", err)
		return nil
	}
	return client
}

func (p *EesMysqlSvr) GetClients(page int, limit int, condition map[string]string) controller.Clients {
	var clients []*structure.Client

	qs := p.o.QueryTable("client")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&clients)
	if err != nil {
		seelog.Error("EesMysqlSvr::GetClients 1 err:", err)
		return controller.Clients{}
	}

	//err = p.o.Raw("select count(*) from client").QueryRow(&count)
	count, err := qs.Count()
	if err != nil {
		seelog.Error("EesMysqlSvr::GetClients 2 err:", err)
		return controller.Clients{}
	}

	back := controller.Clients{
		Count: int(count),
		Data:  clients,
	}
	return back
}

func (p *EesMysqlSvr) SetAvatar(idx string, avatar string) bool {
	client := structure.Client{Idx: idx, Avatar: avatar}
	_, err := p.o.Update(&client, "avatar")
	if err != nil {
		seelog.Error("EesMysqlSvr::SetAvatar err: ", err)
		return false
	}
	return true
}

func (p *EesMysqlSvr) UpdateClient(client *structure.Client, change_chinaid bool) bool {
	_, err := p.o.Update(client, "china_id", "name", "phone", "addr",
		"healthy", "remarks", "contacts", "slow_ill")
	if err != nil {
		seelog.Error("EesMysqlSvr::UpdateClient err: ", err)
		return false
	}

	if change_chinaid {
		p.o.QueryTable("record").Filter("client_idx", client.Idx).Update(orm.Params{
			"china_id": client.ChinaId,
		})
	}

	return true
}
