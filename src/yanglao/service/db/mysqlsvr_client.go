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

func (p *Mysqlsvr) GetClients(page int, limit int) controller.Clients {
	var clients []*structure.Client
	_, err := p.o.QueryTable("client").Limit(limit, (page-1)*limit).All(&clients)
	if err != nil {
		seelog.Error("Mysqlsvr::GetClients 1 err:", err)
	}

	count := 0
	err = p.o.Raw("select count(*) from client").QueryRow(&count)
	if err != nil {
		seelog.Error("Mysqlsvr::GetClients 2 err:", err)
	}

	back := controller.Clients{
		Count: count,
		Data:  clients,
	}
	return back
}
