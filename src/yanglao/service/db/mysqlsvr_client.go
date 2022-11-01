package db

import (
	"fmt"
	"strings"

	"yanglao/constant"
	"yanglao/controller"
	"yanglao/static"
	"yanglao/structure"

	"github.com/cihub/seelog"
)

////////////////////////////////////////////////////////////////////////////////
// Handle Client
////////////////////////////////////////////////////////////////////////////////
func (p *Mysqlsvr) CheckChinaID(id string) uint16 {
	exist := p.o.QueryTable("client").Filter("china_id", strings.ToUpper(id)).Exist()
	if exist {
		return constant.RegisterResult_RepeatID
	}
	return constant.RegisterResult_Success
}

func (p *Mysqlsvr) RegisterClient(client *structure.Client) bool {
	role_index := &structure.Indexs{Key: structure.Indexs_Role_Index}
	err := p.o.Read(role_index)
	if err != nil {
		seelog.Error("Mysqlsvr::RegisterClient get role index err: ", err)
		return false
	}

	index := role_index.Value
	role_index.Value += 1
	_, err = p.o.Update(role_index)
	if err != nil {
		seelog.Error("Mysqlsvr::RegisterClient update role index err: ", err)
		return false
	}

	client.Idx = fmt.Sprintf("%s%010d", static.Db.IndexHead, index)
	_, err = p.o.Insert(client)
	if err != nil {
		seelog.Error("Mysqlsvr::RegisterClient save role err: ", err)
		return false
	}
	return true
}

func (p *Mysqlsvr) GetClients(page int, limit int, condition map[string]string) controller.Clients {
	var clients []*structure.Client

	seelog.Info(condition)
	qs := p.o.QueryTable("client")
	for k, v := range condition {
		qs = qs.Filter(k, v)
	}

	q := qs
	_, err := q.Limit(limit, (page-1)*limit).All(&clients)
	if err != nil {
		seelog.Error("Mysqlsvr::GetClients 1 err:", err)
		return controller.Clients{}
	}

	//err = p.o.Raw("select count(*) from client").QueryRow(&count)
	count, err := qs.Count()
	if err != nil {
		seelog.Error("Mysqlsvr::GetClients 2 err:", err)
		return controller.Clients{}
	}

	back := controller.Clients{
		Count: int(count),
		Data:  clients,
	}
	return back
}

func (p *Mysqlsvr) SetAvatar(idx string, avatar string) bool {
	client := structure.Client{Idx: idx, Avatar: avatar}
	_, err := p.o.Update(&client, "avatar")
	if err != nil {
		seelog.Error("Mysqlsvr::SetAvatar err: ", err)
		return false
	}
	return true
}
